package provider

import (
	"context"
	"encoding/base64"

	"fmt"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilesystemPutResource struct {
	client *client.Client
}

type FilesystemPutResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Path        types.String `tfsdk:"path"`
	Options     types.String `tfsdk:"options"`
	FileContent types.String `tfsdk:"file_content"`
}

func NewFilesystemPutResource() resource.Resource {
	return &FilesystemPutResource{}
}

func (r *FilesystemPutResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_put"
}

func (r *FilesystemPutResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Job to put contents to `path`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"path":    schema.StringAttribute{Required: true, MarkdownDescription: "Path where the file should be written."},
			"options": schema.StringAttribute{Optional: true, MarkdownDescription: "Options controlling file writing behavior."},
			"file_content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded file content for upload (optional, only needed for key file uploads)",
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

	// Build parameters
	params := make(map[string]interface{})
	params["path"] = data.Path.ValueString()
	if !data.Options.IsNull() {
		params["options"] = data.Options.ValueString()
	}

	// Decode file content if provided
	var fileContent []byte
	if !data.FileContent.IsNull() && data.FileContent.ValueString() != "" {
		var err error
		fileContent, err = base64.StdEncoding.DecodeString(data.FileContent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid File Content", fmt.Sprintf("Failed to decode base64: %s", err.Error()))
			return
		}
	}

	// Execute via HTTP multipart upload
	endpoint := "/api/v2.0/filesystem/put"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute filesystem.put: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	// Set ID
	data.ID = data.Path
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemPutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FilesystemPutResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// File upload resources are write-only, state is maintained
}

func (r *FilesystemPutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FilesystemPutResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from state
	var state FilesystemPutResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID

	// Build parameters
	params := make(map[string]interface{})
	params["path"] = data.Path.ValueString()
	if !data.Options.IsNull() {
		params["options"] = data.Options.ValueString()
	}

	// Decode file content if provided
	var fileContent []byte
	if !data.FileContent.IsNull() && data.FileContent.ValueString() != "" {
		var err error
		fileContent, err = base64.StdEncoding.DecodeString(data.FileContent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid File Content", fmt.Sprintf("Failed to decode base64: %s", err.Error()))
			return
		}
	}

	// Execute via HTTP multipart upload
	endpoint := "/api/v2.0/filesystem/put"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute filesystem.put: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemPutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Note: TrueNAS does not provide an API method to delete files
	// The uploaded file will remain on the filesystem after destroy
	// Manual cleanup required if needed
}
