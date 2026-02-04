package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type VirtVolumeResource struct {
	client *client.Client
}

type VirtVolumeResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ContentType types.String `tfsdk:"content_type"`
	Size        types.Int64  `tfsdk:"size"`
	StoragePool types.String `tfsdk:"storage_pool"`
}

func NewVirtVolumeResource() resource.Resource {
	return &VirtVolumeResource{}
}

func (r *VirtVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virt_volume"
}

func (r *VirtVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *VirtVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS virt_volume resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Name for the new virtualization volume (alphanumeric, dashes, dots, underscores).",
			},
			"content_type": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"size": schema.Int64Attribute{
				Required:    false,
				Optional:    true,
				Description: "New size for the volume in MB (minimum 512MB).",
			},
			"storage_pool": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Storage pool in which to create the volume. This must be one of pools listed     in virt.global.conf",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *VirtVolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VirtVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtVolumeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.ContentType.IsNull() && !data.ContentType.IsUnknown() {
		params["content_type"] = data.ContentType.ValueString()
	}
	if !data.Size.IsNull() && !data.Size.IsUnknown() {
		params["size"] = data.Size.ValueInt64()
	}
	if !data.StoragePool.IsNull() && !data.StoragePool.IsUnknown() {
		params["storage_pool"] = data.StoragePool.ValueString()
	}

	result, err := r.client.Call("virt.volume.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create virt_volume: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	// Read back to populate computed fields
	id := data.ID.ValueString()
	result, err = r.client.Call("virt.volume.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back virt_volume: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	result, err := r.client.Call("virt.volume.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read virt_volume: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtVolumeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VirtVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = state.ID.ValueString()

	params := map[string]interface{}{}
	if !data.Size.IsNull() && !data.Size.IsUnknown() {
		params["size"] = data.Size.ValueInt64()
	}

	_, err = r.client.Call("virt.volume.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update virt_volume: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	_, err := r.client.Call("virt.volume.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete virt_volume: %s", err))
		return
	}
}
