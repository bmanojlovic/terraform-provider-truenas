package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type FilesystemPutResource struct {
	client *client.Client
}

type FilesystemPutResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
	Mode    types.Int64  `tfsdk:"mode"`
}

func NewFilesystemPutResource() resource.Resource {
	return &FilesystemPutResource{}
}

func (r *FilesystemPutResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_put"
}

func (r *FilesystemPutResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Upload file to TrueNAS filesystem",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier (same as path)",
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Destination path on TrueNAS",
			},
			"content": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded file content",
			},
			"mode": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "File permissions (octal, e.g., 0644)",
			},
		},
	}
}

func (r *FilesystemPutResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FilesystemPutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FilesystemPutResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Decode base64 content
	fileContent, err := base64.StdEncoding.DecodeString(data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Content", fmt.Sprintf("Failed to decode base64 content: %s", err.Error()))
		return
	}

	// Prepare JSON data
	jsonData := map[string]interface{}{
		"path": data.Path.ValueString(),
	}
	if !data.Mode.IsNull() {
		jsonData["mode"] = data.Mode.ValueInt64()
	}

	// Upload file
	_, err = r.client.UploadFile("/api/v2.0/filesystem/put", jsonData, fileContent, data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to upload file: %s", err.Error()))
		return
	}

	data.ID = data.Path
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemPutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FilesystemPutResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *FilesystemPutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FilesystemPutResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Decode base64 content
	fileContent, err := base64.StdEncoding.DecodeString(data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Content", fmt.Sprintf("Failed to decode base64 content: %s", err.Error()))
		return
	}

	// Prepare JSON data
	jsonData := map[string]interface{}{
		"path": data.Path.ValueString(),
	}
	if !data.Mode.IsNull() {
		jsonData["mode"] = data.Mode.ValueInt64()
	}

	// Upload file
	_, err = r.client.UploadFile("/api/v2.0/filesystem/put", jsonData, fileContent, data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to upload file: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemPutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Optionally delete the file - for now just remove from state
}
