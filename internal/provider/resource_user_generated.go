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

type UserResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID types.String `tfsdk:"id"`
	Uid types.Int64 `tfsdk:"uid"`
	Username types.String `tfsdk:"username"`
	Home types.String `tfsdk:"home"`
	Shell types.String `tfsdk:"shell"`
	FullName types.String `tfsdk:"full_name"`
	Smb types.Bool `tfsdk:"smb"`
	UsernsIdmap types.Int64 `tfsdk:"userns_idmap"`
	Group types.Int64 `tfsdk:"group"`
	Groups types.List `tfsdk:"groups"`
	PasswordDisabled types.Bool `tfsdk:"password_disabled"`
	SshPasswordEnabled types.Bool `tfsdk:"ssh_password_enabled"`
	Sshpubkey types.String `tfsdk:"sshpubkey"`
	Locked types.Bool `tfsdk:"locked"`
	SudoCommands types.List `tfsdk:"sudo_commands"`
	SudoCommandsNopasswd types.List `tfsdk:"sudo_commands_nopasswd"`
	Email types.String `tfsdk:"email"`
	GroupCreate types.Bool `tfsdk:"group_create"`
	HomeCreate types.Bool `tfsdk:"home_create"`
	HomeMode types.String `tfsdk:"home_mode"`
	Password types.String `tfsdk:"password"`
	RandomPassword types.Bool `tfsdk:"random_password"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"uid": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "UNIX UID. If not provided, it is automatically filled with the next one available.",
			},
			"username": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "String used to uniquely identify the user on the server. In order to be portable across     systems,",
			},
			"home": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The local file system path for the user account's home directory. Typically, this is required only i",
			},
			"shell": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Available choices can be retrieved with `user.shell_choices`.",
			},
			"full_name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Comment field to provide additional information about the user account. Typically, this is     the f",
			},
			"smb": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash o",
			},
			"userns_idmap": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to al",
			},
			"group": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "The group entry `id` for the user's primary group. This is not the same as the Unix group `gid` valu",
			},
			"groups": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of additional groups to which the user belongs. NOTE: Groups are identified by their group ent",
			},
			"password_disabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set to `true` password authentication for the user account is disabled.  NOTE: Users with passwor",
			},
			"ssh_password_enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Allow the user to authenticate to the TrueNAS SSH server using a password.  WARNING: The established",
			},
			"sshpubkey": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server.",
			},
			"locked": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS se",
			},
			"sudo_commands": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "An array of commands the user may execute with elevated privileges. User is prompted for password   ",
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "An array of commands the user may execute with elevated privileges. User is *not* prompted for passw",
			},
			"email": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and",
			},
			"group_create": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set to `true`, the TrueNAS server automatically creates a new local group as the user's primary g",
			},
			"home_create": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Create a new home directory for the user in the specified `home` path. ",
			},
			"home_mode": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Filesystem permission to set on the user's home directory. ",
			},
			"password": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The password for the user account. This is required if `random_password` is not set. ",
			},
			"random_password": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Generate a random 20 character password for the user.",
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Uid.IsNull() {
		params["uid"] = data.Uid.ValueInt64()
	}
	if !data.Username.IsNull() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Home.IsNull() {
		params["home"] = data.Home.ValueString()
	}
	if !data.Shell.IsNull() {
		params["shell"] = data.Shell.ValueString()
	}
	if !data.FullName.IsNull() {
		params["full_name"] = data.FullName.ValueString()
	}
	if !data.Smb.IsNull() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Group.IsNull() {
		params["group"] = data.Group.ValueInt64()
	}
	if !data.Groups.IsNull() {
		var groupsList []string
		data.Groups.ElementsAs(ctx, &groupsList, false)
		params["groups"] = groupsList
	}
	if !data.PasswordDisabled.IsNull() {
		params["password_disabled"] = data.PasswordDisabled.ValueBool()
	}
	if !data.SshPasswordEnabled.IsNull() {
		params["ssh_password_enabled"] = data.SshPasswordEnabled.ValueBool()
	}
	if !data.Sshpubkey.IsNull() {
		params["sshpubkey"] = data.Sshpubkey.ValueString()
	}
	if !data.Locked.IsNull() {
		params["locked"] = data.Locked.ValueBool()
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
	if !data.Email.IsNull() {
		params["email"] = data.Email.ValueString()
	}
	if !data.GroupCreate.IsNull() {
		params["group_create"] = data.GroupCreate.ValueBool()
	}
	if !data.HomeCreate.IsNull() {
		params["home_create"] = data.HomeCreate.ValueBool()
	}
	if !data.HomeMode.IsNull() {
		params["home_mode"] = data.HomeMode.ValueString()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}
	if !data.RandomPassword.IsNull() {
		params["random_password"] = data.RandomPassword.ValueBool()
	}

	result, err := r.client.Call("user.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create user: %s", err))
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

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("user.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read user: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["uid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Uid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["home"]; ok && v != nil {
			data.Home = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["shell"]; ok && v != nil {
			data.Shell = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["full_name"]; ok && v != nil {
			data.FullName = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["smb"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Smb = types.BoolValue(bv) }
		}
		if v, ok := resultMap["userns_idmap"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.UsernsIdmap = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["group"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Group = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["groups"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Groups, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["password_disabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.PasswordDisabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["ssh_password_enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.SshPasswordEnabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["sshpubkey"]; ok && v != nil {
			data.Sshpubkey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Locked = types.BoolValue(bv) }
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
		if v, ok := resultMap["email"]; ok && v != nil {
			data.Email = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["group_create"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.GroupCreate = types.BoolValue(bv) }
		}
		if v, ok := resultMap["home_create"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.HomeCreate = types.BoolValue(bv) }
		}
		if v, ok := resultMap["home_mode"]; ok && v != nil {
			data.HomeMode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["random_password"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.RandomPassword = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserResourceModel
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
	if !data.Uid.IsNull() {
		params["uid"] = data.Uid.ValueInt64()
	}
	if !data.Username.IsNull() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Home.IsNull() {
		params["home"] = data.Home.ValueString()
	}
	if !data.Shell.IsNull() {
		params["shell"] = data.Shell.ValueString()
	}
	if !data.FullName.IsNull() {
		params["full_name"] = data.FullName.ValueString()
	}
	if !data.Smb.IsNull() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Group.IsNull() {
		params["group"] = data.Group.ValueInt64()
	}
	if !data.Groups.IsNull() {
		var groupsList []string
		data.Groups.ElementsAs(ctx, &groupsList, false)
		params["groups"] = groupsList
	}
	if !data.PasswordDisabled.IsNull() {
		params["password_disabled"] = data.PasswordDisabled.ValueBool()
	}
	if !data.SshPasswordEnabled.IsNull() {
		params["ssh_password_enabled"] = data.SshPasswordEnabled.ValueBool()
	}
	if !data.Sshpubkey.IsNull() {
		params["sshpubkey"] = data.Sshpubkey.ValueString()
	}
	if !data.Locked.IsNull() {
		params["locked"] = data.Locked.ValueBool()
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
	if !data.Email.IsNull() {
		params["email"] = data.Email.ValueString()
	}
	if !data.GroupCreate.IsNull() {
		params["group_create"] = data.GroupCreate.ValueBool()
	}
	if !data.HomeCreate.IsNull() {
		params["home_create"] = data.HomeCreate.ValueBool()
	}
	if !data.HomeMode.IsNull() {
		params["home_mode"] = data.HomeMode.ValueString()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}
	if !data.RandomPassword.IsNull() {
		params["random_password"] = data.RandomPassword.ValueBool()
	}

	_, err = r.client.Call("user.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update user: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("user.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete user: %s", err))
		return
	}
}
