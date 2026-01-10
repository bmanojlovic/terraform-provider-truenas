package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type NvmetNamespaceResource struct {
	client *client.Client
}

type NvmetNamespaceResourceModel struct {
	ID types.String `tfsdk:"id"`
	Nsid types.Int64 `tfsdk:"nsid"`
	DeviceType types.String `tfsdk:"device_type"`
	DevicePath types.String `tfsdk:"device_path"`
	Filesize types.Int64 `tfsdk:"filesize"`
	Enabled types.Bool `tfsdk:"enabled"`
	SubsysId types.Int64 `tfsdk:"subsys_id"`
}

func NewNvmetNamespaceResource() resource.Resource {
	return &NvmetNamespaceResource{}
}

func (r *NvmetNamespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_namespace"
}

func (r *NvmetNamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NvmetNamespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a NVMe target namespace in a subsystem (`subsys`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"nsid": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Namespace ID (NSID).  Each namespace within a subsystem has an associated NSID, unique within that s",
			},
			"device_type": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Type of device (or file) used to implement the namespace. ",
			},
			"device_path": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Normalized path to the device or file for the namespace.",
			},
			"filesize": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "When `device_type` is \"FILE\" then this will be the size of the file in bytes.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If `enabled` is `False` then the namespace will not be accessible.  Some namespace configuration cha",
			},
			"subsys_id": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the NVMe-oF subsystem to contain this namespace.",
			},
		},
	}
}

func (r *NvmetNamespaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NvmetNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NvmetNamespaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Nsid.IsNull() {
		params["nsid"] = data.Nsid.ValueInt64()
	}
	if !data.DeviceType.IsNull() {
		params["device_type"] = data.DeviceType.ValueString()
	}
	if !data.DevicePath.IsNull() {
		params["device_path"] = data.DevicePath.ValueString()
	}
	if !data.Filesize.IsNull() {
		params["filesize"] = data.Filesize.ValueInt64()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.SubsysId.IsNull() {
		params["subsys_id"] = data.SubsysId.ValueInt64()
	}

	result, err := r.client.Call("nvmet.namespace.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create nvmet_namespace: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetNamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("nvmet.namespace.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read nvmet_namespace: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["nsid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Nsid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["device_type"]; ok && v != nil {
			data.DeviceType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["device_path"]; ok && v != nil {
			data.DevicePath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["filesize"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Filesize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["subsys_id"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.SubsysId = types.Int64Value(int64(fv)) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetNamespaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NvmetNamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Nsid.IsNull() {
		params["nsid"] = data.Nsid.ValueInt64()
	}
	if !data.DeviceType.IsNull() {
		params["device_type"] = data.DeviceType.ValueString()
	}
	if !data.DevicePath.IsNull() {
		params["device_path"] = data.DevicePath.ValueString()
	}
	if !data.Filesize.IsNull() {
		params["filesize"] = data.Filesize.ValueInt64()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.SubsysId.IsNull() {
		params["subsys_id"] = data.SubsysId.ValueInt64()
	}

	_, err = r.client.Call("nvmet.namespace.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update nvmet_namespace: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetNamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetNamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("nvmet.namespace.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete nvmet_namespace: %s", err))
		return
	}
}
