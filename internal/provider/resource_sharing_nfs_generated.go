package provider

import (
	"context"
	"fmt"
	"strconv"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

func (r *SharingNfsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SharingNfsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a NFS Share.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Local path to be exported. ",
			},
			"aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "IGNORED for now. ",
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "User comment associated with share. ",
			},
			"networks": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of authorized networks that are allowed to access the share having format     \"network/mask\" CI",
			},
			"hosts": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of IP's/hostnames which are allowed to access the share. No quotes or spaces are allowed. Each ",
			},
			"ro": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Export the share as read only. ",
			},
			"maproot_user": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Map root user client to a specified user. ",
			},
			"maproot_group": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Map root group client to a specified group. ",
			},
			"mapall_user": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Map all client users to a specified user. ",
			},
			"mapall_group": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Map all client groups to a specified group. ",
			},
			"security": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Specify the security schema. ",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable or disable the share. ",
			},
			"expose_snapshots": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enterprise feature to enable access to the ZFS snapshot directory for the export. Export path must b",
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
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Aliases.IsNull() {
		var aliasesList []string
		data.Aliases.ElementsAs(ctx, &aliasesList, false)
		params["aliases"] = aliasesList
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Networks.IsNull() {
		var networksList []string
		data.Networks.ElementsAs(ctx, &networksList, false)
		params["networks"] = networksList
	}
	if !data.Hosts.IsNull() {
		var hostsList []string
		data.Hosts.ElementsAs(ctx, &hostsList, false)
		params["hosts"] = hostsList
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
	if !data.Security.IsNull() {
		var securityList []string
		data.Security.ElementsAs(ctx, &securityList, false)
		params["security"] = securityList
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ExposeSnapshots.IsNull() {
		params["expose_snapshots"] = data.ExposeSnapshots.ValueBool()
	}

	result, err := r.client.Call("sharing.nfs.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create sharing_nfs: %s", err))
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

func (r *SharingNfsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	result, err := r.client.Call("sharing.nfs.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read sharing_nfs: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["path"]; ok && v != nil {
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["aliases"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Aliases, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			data.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["networks"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Networks, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["hosts"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Hosts, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["ro"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Ro = types.BoolValue(bv) }
		}
		if v, ok := resultMap["maproot_user"]; ok && v != nil {
			data.MaprootUser = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["maproot_group"]; ok && v != nil {
			data.MaprootGroup = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mapall_user"]; ok && v != nil {
			data.MapallUser = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mapall_group"]; ok && v != nil {
			data.MapallGroup = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["security"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Security, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["expose_snapshots"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.ExposeSnapshots = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingNfsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SharingNfsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	params := map[string]interface{}{}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Aliases.IsNull() {
		var aliasesList []string
		data.Aliases.ElementsAs(ctx, &aliasesList, false)
		params["aliases"] = aliasesList
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Networks.IsNull() {
		var networksList []string
		data.Networks.ElementsAs(ctx, &networksList, false)
		params["networks"] = networksList
	}
	if !data.Hosts.IsNull() {
		var hostsList []string
		data.Hosts.ElementsAs(ctx, &hostsList, false)
		params["hosts"] = hostsList
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
	if !data.Security.IsNull() {
		var securityList []string
		data.Security.ElementsAs(ctx, &securityList, false)
		params["security"] = securityList
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ExposeSnapshots.IsNull() {
		params["expose_snapshots"] = data.ExposeSnapshots.ValueBool()
	}

	_, err = r.client.Call("sharing.nfs.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update sharing_nfs: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingNfsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SharingNfsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	_, err = r.client.Call("sharing.nfs.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete sharing_nfs: %s", err))
		return
	}
}
