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

type VirtVolumeImport_IsoResource struct {
	client *client.Client
}

type VirtVolumeImport_IsoResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VirtVolumeImportIso types.String `tfsdk:"virt_volume_import_iso"`
	FileContent         types.String `tfsdk:"file_content"`
}

func NewVirtVolumeImport_IsoResource() resource.Resource {
	return &VirtVolumeImport_IsoResource{}
}

func (r *VirtVolumeImport_IsoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virt_volume_import_iso"
}

func (r *VirtVolumeImport_IsoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Upload via virt.volume.import_iso",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"virt_volume_import_iso": schema.StringAttribute{Required: true, MarkdownDescription: "VirtVolumeImportIsoArgs parameters."},
			"file_content": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Base64-encoded file content for upload (optional, only needed for key file uploads)",
			},
		},
	}
}

func (r *VirtVolumeImport_IsoResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VirtVolumeImport_IsoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtVolumeImport_IsoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := make(map[string]interface{})
	params["virt_volume_import_iso"] = data.VirtVolumeImportIso.ValueString()

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
	endpoint := "/api/v2.0/virt/volume/import_iso"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute virt.volume.import_iso: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	// Set ID
	data.ID = data.VirtVolumeImportIso
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtVolumeImport_IsoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtVolumeImport_IsoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// File upload resources are write-only, state is maintained
}

func (r *VirtVolumeImport_IsoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtVolumeImport_IsoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from state
	var state VirtVolumeImport_IsoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID

	// Build parameters
	params := make(map[string]interface{})
	params["virt_volume_import_iso"] = data.VirtVolumeImportIso.ValueString()

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
	endpoint := "/api/v2.0/virt/volume/import_iso"
	_, err := r.client.UploadFile(endpoint, params, fileContent, "upload")
	if err != nil {
		resp.Diagnostics.AddError("Upload Failed", fmt.Sprintf("Failed to execute virt.volume.import_iso: %s", err.Error()))
		return
	}

	// Note: Upload returns job ID but we don't wait - job completes in background

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtVolumeImport_IsoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Note: TrueNAS does not provide an API method to delete files
	// The uploaded file will remain on the filesystem after destroy
	// Manual cleanup required if needed
}
