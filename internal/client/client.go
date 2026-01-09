package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	host       string
	token      string
	conn       *websocket.Conn
	httpClient *http.Client
	mu         sync.Mutex
	requests   map[string]chan DDPResponse
	nextID     int
}

type DDPMessage struct {
	Msg     string      `json:"msg"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	ID      string      `json:"id,omitempty"`
	Version string      `json:"version,omitempty"`
	Support []string    `json:"support,omitempty"`
}

type DDPResponse struct {
	Msg     string      `json:"msg"`
	ID      string      `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Session string      `json:"session,omitempty"`
}

func NewClient(host, token string) (*Client, error) {
	return &Client{
		host:     host,
		token:    token,
		requests: make(map[string]chan DDPResponse),
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) Connect() error {
	// Force HTTP/1.1 for WebSocket upgrade
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
	
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+c.token)
	
	url := fmt.Sprintf("wss://%s/websocket", c.host)
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %v", err)
	}
	
	c.conn = conn
	
	// Start message handler
	go c.handleMessages()
	
	// Send DDP connect
	connectMsg := DDPMessage{
		Msg:     "connect",
		Version: "1",
		Support: []string{"1"},
	}
	
	if err := c.conn.WriteJSON(connectMsg); err != nil {
		return fmt.Errorf("failed to send connect: %v", err)
	}
	
	// Wait for connected response
	time.Sleep(200 * time.Millisecond)
	
	// Authenticate
	authResp, err := c.call("auth.login_with_api_key", []interface{}{c.token})
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	
	if result, ok := authResp.Result.(bool); !ok || !result {
		return fmt.Errorf("authentication failed: %v", authResp.Error)
	}
	
	log.Println("WebSocket DDP connection established and authenticated")
	return nil
}

func (c *Client) handleMessages() {
	for {
		var response DDPResponse
		if err := c.conn.ReadJSON(&response); err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}
		
		if response.ID != "" {
			c.mu.Lock()
			if ch, exists := c.requests[response.ID]; exists {
				ch <- response
				delete(c.requests, response.ID)
			}
			c.mu.Unlock()
		}
	}
}

func (c *Client) call(method string, params interface{}) (*DDPResponse, error) {
	c.mu.Lock()
	c.nextID++
	id := fmt.Sprintf("req%d", c.nextID)
	respChan := make(chan DDPResponse, 1)
	c.requests[id] = respChan
	
	msg := DDPMessage{
		Msg:    "method",
		Method: method,
		Params: params,
		ID:     id,
	}
	
	if err := c.conn.WriteJSON(msg); err != nil {
		delete(c.requests, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("failed to send message: %v", err)
	}
	c.mu.Unlock()
	
	select {
	case response := <-respChan:
		return &response, nil
	case <-time.After(30 * time.Second):
		c.mu.Lock()
		delete(c.requests, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("request timeout")
	}
}

func (c *Client) Call(method string, params interface{}) (interface{}, error) {
	// For DDP protocol, params should be wrapped in array unless already an array
	var ddpParams interface{}
	switch method {
	case "vm.update", "vm.stop":
		// These expect [id, data] format - params should already be correct
		ddpParams = params
	case "vm.delete", "vm.get_instance":
		// These expect integer ID parameter
		if idStr, ok := params.(string); ok {
			// Convert string ID to integer
			if id, err := strconv.Atoi(idStr); err == nil {
				ddpParams = []interface{}{id}
			} else {
				return nil, fmt.Errorf("invalid ID format: %s", idStr)
			}
		} else {
			ddpParams = []interface{}{params}
		}
	default:
		// vm.create and others expect [data] format
		ddpParams = []interface{}{params}
	}
	
	response, err := c.call(method, ddpParams)
	if err != nil {
		return nil, err
	}
	
	if response.Error != nil {
		return nil, fmt.Errorf("%s failed: %v", method, response.Error)
	}
	
	return response.Result, nil
}

// UploadFile performs a multipart file upload to the specified endpoint
func (c *Client) UploadFile(endpoint string, jsonData map[string]interface{}, fileContent []byte, filename string) (interface{}, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Add JSON data part
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %v", err)
	}
	if err := writer.WriteField("data", string(jsonBytes)); err != nil {
		return nil, fmt.Errorf("failed to write data field: %v", err)
	}
	
	// Add file part
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		return nil, fmt.Errorf("failed to write file content: %v", err)
	}
	
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}
	
	// Create HTTP request
	url := fmt.Sprintf("https://%s%s", c.host, endpoint)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.token)
	
	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	return result, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
