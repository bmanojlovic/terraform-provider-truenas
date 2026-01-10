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

type PrivilegeResource struct {
	client *client.Client
}

type PrivilegeResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	LocalGroups types.List `tfsdk:"local_groups"`
	DsGroups types.List `tfsdk:"ds_groups"`
	Roles types.List `tfsdk:"roles"`
	WebShell types.Bool `tfsdk:"web_shell"`
}

func NewPrivilegeResource() resource.Resource {
	return &PrivilegeResource{}
}

func (r *PrivilegeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privilege"
}

func (r *PrivilegeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PrivilegeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a privilege.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Display name of the privilege.",
			},
			"local_groups": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of local group IDs to assign to this privilege.",
			},
			"ds_groups": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of directory service group IDs or SIDs to assign to this privilege.",
			},
			"roles": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of role names included in this privilege.",
			},
			"web_shell": schema.BoolAttribute{
				Required: true,
				Optional: false,
				Description: "Whether this privilege grants access to the web shell.",
			},
		},
	}
}

func (r *PrivilegeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PrivilegeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PrivilegeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.LocalGroups.IsNull() {
		var local_groupsList []string
		data.LocalGroups.ElementsAs(ctx, &local_groupsList, false)
		params["local_groups"] = local_groupsList
	}
	if !data.DsGroups.IsNull() {
		var ds_groupsList []string
		data.DsGroups.ElementsAs(ctx, &ds_groupsList, false)
		params["ds_groups"] = ds_groupsList
	}
	if !data.Roles.IsNull() {
		var rolesList []string
		data.Roles.ElementsAs(ctx, &rolesList, false)
		params["roles"] = rolesList
	}
	if !data.WebShell.IsNull() {
		params["web_shell"] = data.WebShell.ValueBool()
	}

	result, err := r.client.Call("privilege.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create privilege: %s", err))
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

func (r *PrivilegeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PrivilegeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("privilege.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read privilege: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["local_groups"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.LocalGroups, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["ds_groups"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.DsGroups, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["roles"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Roles, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["web_shell"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.WebShell = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PrivilegeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PrivilegeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PrivilegeResourceModel
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
	if !data.LocalGroups.IsNull() {
		var local_groupsList []string
		data.LocalGroups.ElementsAs(ctx, &local_groupsList, false)
		params["local_groups"] = local_groupsList
	}
	if !data.DsGroups.IsNull() {
		var ds_groupsList []string
		data.DsGroups.ElementsAs(ctx, &ds_groupsList, false)
		params["ds_groups"] = ds_groupsList
	}
	if !data.Roles.IsNull() {
		var rolesList []string
		data.Roles.ElementsAs(ctx, &rolesList, false)
		params["roles"] = rolesList
	}
	if !data.WebShell.IsNull() {
		params["web_shell"] = data.WebShell.ValueBool()
	}

	_, err = r.client.Call("privilege.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update privilege: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PrivilegeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PrivilegeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("privilege.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete privilege: %s", err))
		return
	}
}
