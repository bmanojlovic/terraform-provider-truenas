package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type VmDeviceResource struct {
	client *client.Client
}

type VmDeviceResourceModel struct {
	ID types.String `tfsdk:"id"`
	Attributes types.String `tfsdk:"attributes"`
	Vm types.Int64 `tfsdk:"vm"`
	Order types.Int64 `tfsdk:"order"`
}

func NewVmDeviceResource() resource.Resource {
	return &VmDeviceResource{}
}

func (r *VmDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_device"
}

func (r *VmDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *VmDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS vm_device resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"attributes": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Device-specific configuration attributes.",
			},
			"vm": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the virtual machine this device belongs to.",
			},
			"order": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Boot order priority for this device. `null` for automatic assignment.",
			},
		},
	}
}

func (r *VmDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VmDeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Vm.IsNull() && !data.Vm.IsUnknown() {
		params["vm"] = data.Vm.ValueInt64()
	}
	if !data.Order.IsNull() && !data.Order.IsUnknown() {
		params["order"] = data.Order.ValueInt64()
	}

	result, err := r.client.Call("vm.device.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("vm.device.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state VmDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Vm.IsNull() && !data.Vm.IsUnknown() {
		params["vm"] = data.Vm.ValueInt64()
	}
	if !data.Order.IsNull() && !data.Order.IsUnknown() {
		params["order"] = data.Order.ValueInt64()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("vm.device.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	// Try to delete device
	// If VM is running, this will fail - that's OK, VM deletion will cascade
	_, err = r.client.Call("vm.device.delete", resourceID)
	if err != nil {
		// Ignore ENOENT - device already deleted (e.g., VM was deleted with force=true)
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		// Ignore "VM must be stopped" errors - VM will cascade delete devices
		if strings.Contains(err.Error(), "stop") || strings.Contains(err.Error(), "running") {
			tflog.Info(ctx, "Device deletion skipped - VM is running, will be cleaned up by VM deletion")
			return
		}
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Delete zvol if it was created by this device
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err == nil {
			if createZvol, ok := attributesObj["create_zvol"].(bool); ok && createZvol {
				if zvolName, ok := attributesObj["zvol_name"].(string); ok && zvolName != "" {
					// Delete the zvol dataset - API returns null if dataset doesn't exist
					deleteParams := []interface{}{zvolName, map[string]interface{}{"force": true}}
					if _, err := r.client.Call("pool.dataset.delete", deleteParams); err != nil {
						// Log warning but don't fail - zvol might already be deleted
						tflog.Warn(ctx, "Failed to delete zvol", map[string]interface{}{"zvol": zvolName, "error": err.Error()})
					}
				}
			}
		}
	}
}
