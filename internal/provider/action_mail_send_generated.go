package provider

import (
	"context"
	"encoding/base64"

	"fmt"
	"time"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ActionMailSendResource struct {
	client *client.Client
}

type ActionMailSendResourceModel struct {
	Message types.String `tfsdk:"message"`
	Config  types.String `tfsdk:"config"`
	// File upload (optional)
	FileContent types.String `tfsdk:"file_content"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewActionMailSendResource() resource.Resource {
	return &ActionMailSendResource{}
}

func (r *ActionMailSendResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action_mail_send"
}

func (r *ActionMailSendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sends mail using configured mail settings.",
		Attributes: map[string]schema.Attribute{
			"message": schema.StringAttribute{Required: true, MarkdownDescription: "Email message content and configuration."},
			"config":  schema.StringAttribute{Optional: true, MarkdownDescription: "Optional mail configuration overrides for this message."},
			"file_content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded file content for upload (optional)",
			},
			"action_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Action execution identifier",
			},
			"job_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Background job ID (if applicable)",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Job state: SUCCESS, FAILED, or RUNNING",
			},
			"progress": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "Job progress percentage (0-100)",
			},
			"result": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Action result data",
			},
			"error": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Error message if action failed",
			},
		},
	}
}

func (r *ActionMailSendResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
		return
	}
	r.client = client
}

func (r *ActionMailSendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ActionMailSendResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := make(map[string]interface{})
	params["message"] = data.Message.ValueString()
	if !data.Config.IsNull() {
		params["config"] = data.Config.ValueString()
	}

	// Prepare file content if provided
	var fileContent []byte
	var err error
	if !data.FileContent.IsNull() {
		fileContent, err = base64.StdEncoding.DecodeString(data.FileContent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid File Content", fmt.Sprintf("Failed to decode base64: %s", err.Error()))
			return
		}
	}

	// Execute via HTTP multipart upload
	endpoint := "/api/v2.0/mail/send"
	result, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute mail.send: %s", err.Error()))
		return
	}

	// Store job ID if returned
	if jobID, ok := result.(float64); ok {
		data.JobID = types.Int64Value(int64(jobID))
	}

	// Actions are fire-and-forget - mark as success immediately
	data.State = types.StringValue("SUBMITTED")
	data.Progress = types.Float64Value(0.0)
	data.Result = types.StringValue(fmt.Sprintf("%v", result))
	data.Error = types.StringValue("")

	// Generate ID from timestamp
	data.ActionID = types.StringValue(fmt.Sprintf("mail.send-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ActionMailSendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data ActionMailSendResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *ActionMailSendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Actions cannot be updated - force recreation
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated. Please destroy and recreate the resource.")
}

func (r *ActionMailSendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Actions cannot be undone - just remove from state
}
