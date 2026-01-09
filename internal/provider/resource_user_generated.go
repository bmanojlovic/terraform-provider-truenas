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

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS user resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"uid": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"username": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"home": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"shell": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"full_name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"smb": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"userns_idmap": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"group": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"password_disabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"ssh_password_enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"sshpubkey": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"locked": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"sudo_commands": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"email": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"group_create": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"home_create": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"home_mode": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"password": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"random_password": schema.BoolAttribute{
				Required: false,
				Optional: true,
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
	params["username"] = data.Username.ValueString()
	if !data.Home.IsNull() {
		params["home"] = data.Home.ValueString()
	}
	if !data.Shell.IsNull() {
		params["shell"] = data.Shell.ValueString()
	}
	params["full_name"] = data.FullName.ValueString()
	if !data.Smb.IsNull() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Group.IsNull() {
		params["group"] = data.Group.ValueInt64()
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

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("user.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Uid.IsNull() {
		params["uid"] = data.Uid.ValueInt64()
	}
	params["username"] = data.Username.ValueString()
	if !data.Home.IsNull() {
		params["home"] = data.Home.ValueString()
	}
	if !data.Shell.IsNull() {
		params["shell"] = data.Shell.ValueString()
	}
	params["full_name"] = data.FullName.ValueString()
	if !data.Smb.IsNull() {
		params["smb"] = data.Smb.ValueBool()
	}
	if !data.UsernsIdmap.IsNull() {
		params["userns_idmap"] = data.UsernsIdmap.ValueInt64()
	}
	if !data.Group.IsNull() {
		params["group"] = data.Group.ValueInt64()
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

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("user.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("user.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
