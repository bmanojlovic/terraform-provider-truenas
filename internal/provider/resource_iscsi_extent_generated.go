package provider

import (
	"context"
	"fmt"
	"strconv"

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
	Filesize types.String `tfsdk:"filesize"`
	Blocksize types.Int64 `tfsdk:"blocksize"`
	Pblocksize types.Bool `tfsdk:"pblocksize"`
	AvailThreshold types.String `tfsdk:"avail_threshold"`
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

func (r *IscsiExtentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS iscsi_extent resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"disk": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"serial": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"path": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"filesize": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"blocksize": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"pblocksize": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"avail_threshold": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"insecure_tpc": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"xen": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"rpm": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"ro": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"product_id": schema.StringAttribute{
				Required: false,
				Optional: true,
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
	params["name"] = data.Name.ValueString()
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
		params["filesize"] = data.Filesize.ValueString()
	}
	if !data.Blocksize.IsNull() {
		params["blocksize"] = data.Blocksize.ValueInt64()
	}
	if !data.Pblocksize.IsNull() {
		params["pblocksize"] = data.Pblocksize.ValueBool()
	}
	if !data.AvailThreshold.IsNull() {
		params["avail_threshold"] = data.AvailThreshold.ValueString()
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

	result, err := r.client.Call("iscsi/extent.create", params)
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

func (r *IscsiExtentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("iscsi/extent.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
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
		params["filesize"] = data.Filesize.ValueString()
	}
	if !data.Blocksize.IsNull() {
		params["blocksize"] = data.Blocksize.ValueInt64()
	}
	if !data.Pblocksize.IsNull() {
		params["pblocksize"] = data.Pblocksize.ValueBool()
	}
	if !data.AvailThreshold.IsNull() {
		params["avail_threshold"] = data.AvailThreshold.ValueString()
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

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("iscsi/extent.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("iscsi/extent.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
