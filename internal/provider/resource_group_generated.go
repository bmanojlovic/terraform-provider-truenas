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

type GroupResource struct {
	client *client.Client
}

type GroupResourceModel struct {
	ID types.String `tfsdk:"id"`
	Gid types.Int64 `tfsdk:"gid"`
	Name types.String `tfsdk:"name"`
	SudoCommands types.List `tfsdk:"sudo_commands"`
	SudoCommandsNopasswd types.List `tfsdk:"sudo_commands_nopasswd"`
	Smb types.Bool `tfsdk:"smb"`
	UsernsIdmap types.Int64 `tfsdk:"userns_idmap"`
	Users types.List `tfsdk:"users"`
}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"gid": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "If `null`, it is automatically filled with the next one available.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "",
			},
			"sudo_commands": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "A list of commands that group members may execute with elevated privileges. User is prompted for pas",
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "A list of commands that group members may execute with elevated privileges. User is not prompted for",
			},
			"smb": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT grou",
			},
			"userns_idmap": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to a",
			},
			"users": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "A list a API user identifiers for local users who are members of this group. These IDs match the `id",
			},
		},
	}
}

func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Gid.IsNull() {
		params["gid"] = data.Gid.ValueInt64()
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.SudoCommands.IsNull() {
		var sudo_commandsList []string
		data.SudoCommands.ElementsAs(ctx, &sudo_commandsList, false)
		params["sudo_commands"] = sudo_commandsList
	}
	if !data.SudoCommandsNopasswd.IsNull() {
		var sudo_commands_nopasswdList []string
		data.SudoCommandsNopasswd.ElementsAs(ctx, &sudo_commands_nopasswdList, false)
		params["sudo_commands_nopasswd"] = sudo_commands_nopasswdList
	}
	if !data.Smb.IsNull() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Users.IsNull() {
		var usersList []string
		data.Users.ElementsAs(ctx, &usersList, false)
		params["users"] = usersList
	}

	result, err := r.client.Call("group.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create group: %s", err))
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

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("group.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read group: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["gid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Gid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sudo_commands"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.SudoCommands, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["sudo_commands_nopasswd"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.SudoCommandsNopasswd, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["smb"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Smb = types.BoolValue(bv) }
		}
		if v, ok := resultMap["userns_idmap"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.UsernsIdmap = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["users"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Users, _ = types.ListValue(types.StringType, strVals)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state GroupResourceModel
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
	if !data.Gid.IsNull() {
		params["gid"] = data.Gid.ValueInt64()
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.SudoCommands.IsNull() {
		var sudo_commandsList []string
		data.SudoCommands.ElementsAs(ctx, &sudo_commandsList, false)
		params["sudo_commands"] = sudo_commandsList
	}
	if !data.SudoCommandsNopasswd.IsNull() {
		var sudo_commands_nopasswdList []string
		data.SudoCommandsNopasswd.ElementsAs(ctx, &sudo_commands_nopasswdList, false)
		params["sudo_commands_nopasswd"] = sudo_commands_nopasswdList
	}
	if !data.Smb.IsNull() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Users.IsNull() {
		var usersList []string
		data.Users.ElementsAs(ctx, &usersList, false)
		params["users"] = usersList
	}

	_, err = r.client.Call("group.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update group: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("group.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete group: %s", err))
		return
	}
}
