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

type IscsiExtentResource struct {
	client *client.Client
}

type IscsiExtentResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Disk types.String `tfsdk:"disk"`
	Serial types.String `tfsdk:"serial"`
	Path types.String `tfsdk:"path"`
	Filesize types.Int64 `tfsdk:"filesize"`
	Blocksize types.Int64 `tfsdk:"blocksize"`
	Pblocksize types.Bool `tfsdk:"pblocksize"`
	AvailThreshold types.Int64 `tfsdk:"avail_threshold"`
	Comment types.String `tfsdk:"comment"`
	InsecureTpc types.Bool `tfsdk:"insecure_tpc"`
	Xen types.Bool `tfsdk:"xen"`
	Rpm types.String `tfsdk:"rpm"`
	Ro types.Bool `tfsdk:"ro"`
	Enabled types.Bool `tfsdk:"enabled"`
	ProductId types.String `tfsdk:"product_id"`
}

func NewIscsiExtentResource() resource.Resource {
	return &IscsiExtentResource{}
}

func (r *IscsiExtentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_extent"
}

func (r *IscsiExtentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiExtentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Extent.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name of the iSCSI extent.",
			},
			"type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Type of the extent storage backend.",
			},
			"disk": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Disk device to use for the extent or `null` if using a file.",
			},
			"serial": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Serial number for the extent or `null` to auto-generate.",
			},
			"path": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "File path for file-based extents or `null` if using a disk.",
			},
			"filesize": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Size of the file-based extent in bytes.",
			},
			"blocksize": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Block size for the extent in bytes.",
			},
			"pblocksize": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to use physical block size reporting.",
			},
			"avail_threshold": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Available space threshold percentage or `null` to disable.",
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional comment describing the extent.",
			},
			"insecure_tpc": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable insecure Third Party Copy (TPC) operations.",
			},
			"xen": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable Xen compatibility mode.",
			},
			"rpm": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Reported RPM type for the extent.",
			},
			"ro": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the extent is read-only.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the extent is enabled and available for use.",
			},
			"product_id": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Product ID string for the extent or `null` for default.",
			},
		},
	}
}

func (r *IscsiExtentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiExtentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Disk.IsNull() {
		params["disk"] = data.Disk.ValueString()
	}
	if !data.Serial.IsNull() {
		params["serial"] = data.Serial.ValueString()
	}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Filesize.IsNull() {
		params["filesize"] = data.Filesize.ValueInt64()
	}
	if !data.Blocksize.IsNull() {
		params["blocksize"] = data.Blocksize.ValueInt64()
	}
	if !data.Pblocksize.IsNull() {
		params["pblocksize"] = data.Pblocksize.ValueBool()
	}
	if !data.AvailThreshold.IsNull() {
		params["avail_threshold"] = data.AvailThreshold.ValueInt64()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.InsecureTpc.IsNull() {
		params["insecure_tpc"] = data.InsecureTpc.ValueBool()
	}
	if !data.Xen.IsNull() {
		params["xen"] = data.Xen.ValueBool()
	}
	if !data.Rpm.IsNull() {
		params["rpm"] = data.Rpm.ValueString()
	}
	if !data.Ro.IsNull() {
		params["ro"] = data.Ro.ValueBool()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ProductId.IsNull() {
		params["product_id"] = data.ProductId.ValueString()
	}

	result, err := r.client.Call("iscsi.extent.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_extent: %s", err))
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

func (r *IscsiExtentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("iscsi.extent.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_extent: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["disk"]; ok && v != nil {
			data.Disk = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["serial"]; ok && v != nil {
			data.Serial = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["filesize"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Filesize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["blocksize"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Blocksize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["pblocksize"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Pblocksize = types.BoolValue(bv) }
		}
		if v, ok := resultMap["avail_threshold"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.AvailThreshold = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			data.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["insecure_tpc"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.InsecureTpc = types.BoolValue(bv) }
		}
		if v, ok := resultMap["xen"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Xen = types.BoolValue(bv) }
		}
		if v, ok := resultMap["rpm"]; ok && v != nil {
			data.Rpm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ro"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Ro = types.BoolValue(bv) }
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["product_id"]; ok && v != nil {
			data.ProductId = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiExtentResourceModel
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Disk.IsNull() {
		params["disk"] = data.Disk.ValueString()
	}
	if !data.Serial.IsNull() {
		params["serial"] = data.Serial.ValueString()
	}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Filesize.IsNull() {
		params["filesize"] = data.Filesize.ValueInt64()
	}
	if !data.Blocksize.IsNull() {
		params["blocksize"] = data.Blocksize.ValueInt64()
	}
	if !data.Pblocksize.IsNull() {
		params["pblocksize"] = data.Pblocksize.ValueBool()
	}
	if !data.AvailThreshold.IsNull() {
		params["avail_threshold"] = data.AvailThreshold.ValueInt64()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.InsecureTpc.IsNull() {
		params["insecure_tpc"] = data.InsecureTpc.ValueBool()
	}
	if !data.Xen.IsNull() {
		params["xen"] = data.Xen.ValueBool()
	}
	if !data.Rpm.IsNull() {
		params["rpm"] = data.Rpm.ValueString()
	}
	if !data.Ro.IsNull() {
		params["ro"] = data.Ro.ValueBool()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ProductId.IsNull() {
		params["product_id"] = data.ProductId.ValueString()
	}

	_, err = r.client.Call("iscsi.extent.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_extent: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("iscsi.extent.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_extent: %s", err))
		return
	}
}
