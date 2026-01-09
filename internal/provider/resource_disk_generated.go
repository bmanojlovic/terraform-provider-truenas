package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type DiskResource struct {
	client *client.Client
}

type DiskResourceModel struct {
	ID types.String `tfsdk:"id"`
	Number types.Int64 `tfsdk:"number"`
	Lunid types.String `tfsdk:"lunid"`
	Description types.String `tfsdk:"description"`
	Hddstandby types.String `tfsdk:"hddstandby"`
	Advpowermgmt types.String `tfsdk:"advpowermgmt"`
	Bus types.String `tfsdk:"bus"`
	Enclosure types.String `tfsdk:"enclosure"`
	Pool types.String `tfsdk:"pool"`
	Passwd types.String `tfsdk:"passwd"`
}

func NewDiskResource() resource.Resource {
	return &DiskResource{}
}

func (r *DiskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disk"
}

func (r *DiskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS disk resource (update-only)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"number": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"lunid": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"hddstandby": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"advpowermgmt": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"bus": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enclosure": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"pool": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"passwd": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *DiskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DiskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DiskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update-only resource: ID must be provided
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "ID is required for update-only resources. Use terraform import to manage existing resources.")
		return
	}

	// Convert string ID to integer for API call
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID must be a valid integer")
		return
	}

	params := map[string]interface{}{}
	if !data.Number.IsNull() {
		params["number"] = data.Number.ValueInt64()
	}
	if !data.Lunid.IsNull() {
		params["lunid"] = data.Lunid.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Hddstandby.IsNull() {
		params["hddstandby"] = data.Hddstandby.ValueString()
	}
	if !data.Advpowermgmt.IsNull() {
		params["advpowermgmt"] = data.Advpowermgmt.ValueString()
	}
	if !data.Bus.IsNull() {
		params["bus"] = data.Bus.ValueString()
	}
	if !data.Enclosure.IsNull() {
		params["enclosure"] = data.Enclosure.ValueString()
	}
	if !data.Pool.IsNull() {
		params["pool"] = data.Pool.ValueString()
	}
	if !data.Passwd.IsNull() {
		params["passwd"] = data.Passwd.ValueString()
	}

	_, err = r.client.Call("disk.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DiskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for API call
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID must be a valid integer")
		return
	}

	result, err := r.client.Call("disk.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	_ = result // TODO: Map result to data fields

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DiskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for API call
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID must be a valid integer")
		return
	}

	params := map[string]interface{}{}
	if !data.Number.IsNull() {
		params["number"] = data.Number.ValueInt64()
	}
	if !data.Lunid.IsNull() {
		params["lunid"] = data.Lunid.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Hddstandby.IsNull() {
		params["hddstandby"] = data.Hddstandby.ValueString()
	}
	if !data.Advpowermgmt.IsNull() {
		params["advpowermgmt"] = data.Advpowermgmt.ValueString()
	}
	if !data.Bus.IsNull() {
		params["bus"] = data.Bus.ValueString()
	}
	if !data.Enclosure.IsNull() {
		params["enclosure"] = data.Enclosure.ValueString()
	}
	if !data.Pool.IsNull() {
		params["pool"] = data.Pool.ValueString()
	}
	if !data.Passwd.IsNull() {
		params["passwd"] = data.Passwd.ValueString()
	}

	_, err = r.client.Call("disk.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Update-only resource: just remove from state, don't delete on server
}

func (r *DiskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID (e.g., "ssh" for SSH service)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
