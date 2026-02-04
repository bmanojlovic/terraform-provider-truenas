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

type ConfigUploadResource struct {
	client *client.Client
}

type ConfigUploadResourceModel struct {
	ID types.String `tfsdk:"id"`

	FileContent types.String `tfsdk:"file_content"`
}

func NewConfigUploadResource() resource.Resource {
	return &ConfigUploadResource{}
}

func (r *ConfigUploadResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_upload"
}

func (r *ConfigUploadResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Accepts a configuration file via job pipe.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},

			"file_content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded file content for upload (optional, only needed for key file uploads)",
			},
		},
	}
}

func (r *ConfigUploadResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConfigUploadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConfigUploadResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := make(map[string]interface{})

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
	endpoint := "/api/v2.0/config/upload"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute config.upload: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	// Set ID
	data.ID = types.StringValue("config.upload")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigUploadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConfigUploadResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// File upload resources are write-only, state is maintained
}

func (r *ConfigUploadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConfigUploadResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from state
	var state ConfigUploadResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID

	// Build parameters
	params := make(map[string]interface{})

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
	endpoint := "/api/v2.0/config/upload"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute config.upload: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigUploadResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Note: TrueNAS does not provide an API method to delete files
	// The uploaded file will remain on the filesystem after destroy
	// Manual cleanup required if needed
}
