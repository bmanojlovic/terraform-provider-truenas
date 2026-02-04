package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "A string used to identify a group.",
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
	if !data.Gid.IsNull() && !data.Gid.IsUnknown() {
		params["gid"] = data.Gid.ValueInt64()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.SudoCommands.IsNull() && !data.SudoCommands.IsUnknown() {
		var sudo_commandsList []string
		data.SudoCommands.ElementsAs(ctx, &sudo_commandsList, false)
		params["sudo_commands"] = sudo_commandsList
	}
	if !data.SudoCommandsNopasswd.IsNull() && !data.SudoCommandsNopasswd.IsUnknown() {
		var sudo_commands_nopasswdList []string
		data.SudoCommandsNopasswd.ElementsAs(ctx, &sudo_commands_nopasswdList, false)
		params["sudo_commands_nopasswd"] = sudo_commands_nopasswdList
	}
	if !data.Smb.IsNull() && !data.Smb.IsUnknown() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() && !data.UsernsIdmap.IsUnknown() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Users.IsNull() && !data.Users.IsUnknown() {
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
	var id interface{}
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("group.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back group: %s", err))
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
		if v, ok := resultMap["userns_idmap"]; ok {
			if v == nil {
				data.UsernsIdmap = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.UsernsIdmap = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.UsernsIdmap = types.Int64Value(int64(fv)) }
					}
				}
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("group.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read group: %s", err))
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
		if v, ok := resultMap["userns_idmap"]; ok {
			if v == nil {
				data.UsernsIdmap = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.UsernsIdmap = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.UsernsIdmap = types.Int64Value(int64(fv)) }
					}
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.SudoCommands.IsNull() && !data.SudoCommands.IsUnknown() {
		var sudo_commandsList []string
		data.SudoCommands.ElementsAs(ctx, &sudo_commandsList, false)
		params["sudo_commands"] = sudo_commandsList
	}
	if !data.SudoCommandsNopasswd.IsNull() && !data.SudoCommandsNopasswd.IsUnknown() {
		var sudo_commands_nopasswdList []string
		data.SudoCommandsNopasswd.ElementsAs(ctx, &sudo_commands_nopasswdList, false)
		params["sudo_commands_nopasswd"] = sudo_commands_nopasswdList
	}
	if !data.Smb.IsNull() && !data.Smb.IsUnknown() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() && !data.UsernsIdmap.IsUnknown() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Users.IsNull() && !data.Users.IsUnknown() {
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	id = []interface{}{id, map[string]interface{}{}}

	_, err = r.client.Call("group.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete group: %s", err))
		return
	}
}
