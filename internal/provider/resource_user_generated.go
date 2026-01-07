package provider

import (
	"context"
	"fmt"

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
	Uid types.String `tfsdk:"uid"`
	Username types.String `tfsdk:"username"`
	Home types.String `tfsdk:"home"`
	Shell types.String `tfsdk:"shell"`
	FullName types.String `tfsdk:"full_name"`
	Smb types.Bool `tfsdk:"smb"`
	UsernsIdmap types.String `tfsdk:"userns_idmap"`
	Group types.String `tfsdk:"group"`
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
			"uid": schema.StringAttribute{
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
			"userns_idmap": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"group": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"groups": schema.ListAttribute{
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
				Required: false,
				Optional: true,
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
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

	params := map[string]interface{}{
		"uid": data.Uid.ValueString(),
		"username": data.Username.ValueString(),
		"home": data.Home.ValueString(),
		"shell": data.Shell.ValueString(),
		"full_name": data.FullName.ValueString(),
		"smb": data.Smb.ValueBool(),
		"userns_idmap": data.UsernsIdmap.ValueString(),
		"group": data.Group.ValueString(),
		"password_disabled": data.PasswordDisabled.ValueBool(),
		"ssh_password_enabled": data.SshPasswordEnabled.ValueBool(),
		"sshpubkey": data.Sshpubkey.ValueString(),
		"locked": data.Locked.ValueBool(),
		"email": data.Email.ValueString(),
		"group_create": data.GroupCreate.ValueBool(),
		"home_create": data.HomeCreate.ValueBool(),
		"home_mode": data.HomeMode.ValueString(),
		"password": data.Password.ValueString(),
		"random_password": data.RandomPassword.ValueBool(),
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

	_, err := r.client.Call("user.get_instance", data.ID.ValueString())
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

	params := map[string]interface{}{
		"uid": data.Uid.ValueString(),
		"username": data.Username.ValueString(),
		"home": data.Home.ValueString(),
		"shell": data.Shell.ValueString(),
		"full_name": data.FullName.ValueString(),
		"smb": data.Smb.ValueBool(),
		"userns_idmap": data.UsernsIdmap.ValueString(),
		"group": data.Group.ValueString(),
		"password_disabled": data.PasswordDisabled.ValueBool(),
		"ssh_password_enabled": data.SshPasswordEnabled.ValueBool(),
		"sshpubkey": data.Sshpubkey.ValueString(),
		"locked": data.Locked.ValueBool(),
		"email": data.Email.ValueString(),
		"group_create": data.GroupCreate.ValueBool(),
		"home_create": data.HomeCreate.ValueBool(),
		"home_mode": data.HomeMode.ValueString(),
		"password": data.Password.ValueString(),
		"random_password": data.RandomPassword.ValueBool(),
	}

	_, err := r.client.Call("user.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("user.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
