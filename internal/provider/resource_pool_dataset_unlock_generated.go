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

type PoolDatasetUnlockResource struct {
	client *client.Client
}

type PoolDatasetUnlockResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DatasetId   types.String `tfsdk:"dataset_id"`
	Options     types.String `tfsdk:"options"`
	FileContent types.String `tfsdk:"file_content"`
}

func NewPoolDatasetUnlockResource() resource.Resource {
	return &PoolDatasetUnlockResource{}
}

func (r *PoolDatasetUnlockResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_dataset_unlock"
}

func (r *PoolDatasetUnlockResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Unlock dataset `id` (and its children if `unlock_options.recursive` is `true`).  If `id` dataset is not encrypted an exception will be raised. There is one exception: when `id` is a root dataset and `",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"dataset_id": schema.StringAttribute{Required: true, MarkdownDescription: "The dataset ID (full path) to unlock."},
			"options":    schema.StringAttribute{Optional: true, MarkdownDescription: "Options for unlocking including force settings, recursion, and dataset-specific keys."},
			"file_content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded file content for upload (optional, only needed for key file uploads)",
			},
		},
	}
}

func (r *PoolDatasetUnlockResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolDatasetUnlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolDatasetUnlockResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := make(map[string]interface{})
	params["id"] = data.DatasetId.ValueString()
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
	endpoint := "/api/v2.0/pool/dataset/unlock"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute pool.dataset.unlock: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	// Set ID
	data.ID = data.DatasetId
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolDatasetUnlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolDatasetUnlockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// File upload resources are write-only, state is maintained
}

func (r *PoolDatasetUnlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolDatasetUnlockResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from state
	var state PoolDatasetUnlockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID

	// Build parameters
	params := make(map[string]interface{})
	params["id"] = data.DatasetId.ValueString()
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
	endpoint := "/api/v2.0/pool/dataset/unlock"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute pool.dataset.unlock: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolDatasetUnlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Note: TrueNAS does not provide an API method to delete files
	// The uploaded file will remain on the filesystem after destroy
	// Manual cleanup required if needed
}
