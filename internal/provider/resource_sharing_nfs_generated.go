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

type SharingNfsResource struct {
	client *client.Client
}

type SharingNfsResourceModel struct {
	ID types.String `tfsdk:"id"`
	Path types.String `tfsdk:"path"`
	Aliases types.List `tfsdk:"aliases"`
	Comment types.String `tfsdk:"comment"`
	Networks types.List `tfsdk:"networks"`
	Hosts types.List `tfsdk:"hosts"`
	Ro types.Bool `tfsdk:"ro"`
	MaprootUser types.String `tfsdk:"maproot_user"`
	MaprootGroup types.String `tfsdk:"maproot_group"`
	MapallUser types.String `tfsdk:"mapall_user"`
	MapallGroup types.String `tfsdk:"mapall_group"`
	Security types.List `tfsdk:"security"`
	Enabled types.Bool `tfsdk:"enabled"`
	ExposeSnapshots types.Bool `tfsdk:"expose_snapshots"`
}

func NewSharingNfsResource() resource.Resource {
	return &SharingNfsResource{}
}

func (r *SharingNfsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sharing_nfs"
}

func (r *SharingNfsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS sharing_nfs resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"networks": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"hosts": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"ro": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"maproot_user": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"maproot_group": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"mapall_user": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"mapall_group": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"security": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"expose_snapshots": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *SharingNfsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SharingNfsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["path"] = data.Path.ValueString()
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Ro.IsNull() {
		params["ro"] = data.Ro.ValueBool()
	}
	if !data.MaprootUser.IsNull() {
		params["maproot_user"] = data.MaprootUser.ValueString()
	}
	if !data.MaprootGroup.IsNull() {
		params["maproot_group"] = data.MaprootGroup.ValueString()
	}
	if !data.MapallUser.IsNull() {
		params["mapall_user"] = data.MapallUser.ValueString()
	}
	if !data.MapallGroup.IsNull() {
		params["mapall_group"] = data.MapallGroup.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ExposeSnapshots.IsNull() {
		params["expose_snapshots"] = data.ExposeSnapshots.ValueBool()
	}

	result, err := r.client.Call("sharing/nfs.create", params)
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

func (r *SharingNfsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("sharing/nfs.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingNfsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state SharingNfsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["path"] = data.Path.ValueString()
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Ro.IsNull() {
		params["ro"] = data.Ro.ValueBool()
	}
	if !data.MaprootUser.IsNull() {
		params["maproot_user"] = data.MaprootUser.ValueString()
	}
	if !data.MaprootGroup.IsNull() {
		params["maproot_group"] = data.MaprootGroup.ValueString()
	}
	if !data.MapallUser.IsNull() {
		params["mapall_user"] = data.MapallUser.ValueString()
	}
	if !data.MapallGroup.IsNull() {
		params["mapall_group"] = data.MapallGroup.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ExposeSnapshots.IsNull() {
		params["expose_snapshots"] = data.ExposeSnapshots.ValueBool()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("sharing/nfs.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingNfsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("sharing/nfs.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
