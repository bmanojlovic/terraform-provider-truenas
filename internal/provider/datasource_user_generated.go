package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *client.Client
}

type UserDataSourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Uid                     types.Int64  `tfsdk:"uid"`
	Username                types.String `tfsdk:"username"`
	Unixhash                types.String `tfsdk:"unixhash"`
	Smbhash                 types.String `tfsdk:"smbhash"`
	Home                    types.String `tfsdk:"home"`
	Shell                   types.String `tfsdk:"shell"`
	FullName                types.String `tfsdk:"full_name"`
	Builtin                 types.Bool   `tfsdk:"builtin"`
	Smb                     types.Bool   `tfsdk:"smb"`
	UsernsIdmap             types.Int64  `tfsdk:"userns_idmap"`
	Group                   types.String `tfsdk:"group"`
	Groups                  types.List   `tfsdk:"groups"`
	PasswordDisabled        types.Bool   `tfsdk:"password_disabled"`
	SshPasswordEnabled      types.Bool   `tfsdk:"ssh_password_enabled"`
	Sshpubkey               types.String `tfsdk:"sshpubkey"`
	Locked                  types.Bool   `tfsdk:"locked"`
	SudoCommands            types.List   `tfsdk:"sudo_commands"`
	SudoCommandsNopasswd    types.List   `tfsdk:"sudo_commands_nopasswd"`
	Email                   types.String `tfsdk:"email"`
	Local                   types.Bool   `tfsdk:"local"`
	Immutable               types.Bool   `tfsdk:"immutable"`
	TwofactorAuthConfigured types.Bool   `tfsdk:"twofactor_auth_configured"`
	Sid                     types.String `tfsdk:"sid"`
	LastPasswordChange      types.String `tfsdk:"last_password_change"`
	PasswordAge             types.Int64  `tfsdk:"password_age"`
	PasswordHistory         types.List   `tfsdk:"password_history"`
	PasswordChangeRequired  types.Bool   `tfsdk:"password_change_required"`
	Roles                   types.List   `tfsdk:"roles"`
	ApiKeys                 types.List   `tfsdk:"api_keys"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"uid": schema.Int64Attribute{
				Computed:    true,
				Description: "A non-negative integer used to identify a system user. TrueNAS uses this value for permission     ch",
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Description: "A string used to identify a user. Local accounts must use characters from the POSIX portable filenam",
			},
			"unixhash": schema.StringAttribute{
				Computed:    true,
				Description: "Hashed password for local accounts. This value is `null` for accounts provided by directory services",
			},
			"smbhash": schema.StringAttribute{
				Computed:    true,
				Description: "NT hash of the local account password for `smb` users. This value is `null` for accounts provided by",
			},
			"home": schema.StringAttribute{
				Computed:    true,
				Description: "The local file system path for the user account's home directory. Typically, this is required only i",
			},
			"shell": schema.StringAttribute{
				Computed:    true,
				Description: "Available choices can be retrieved with `user.shell_choices`.",
			},
			"full_name": schema.StringAttribute{
				Computed:    true,
				Description: "Comment field to provide additional information about the user account. Typically, this is     the f",
			},
			"builtin": schema.BoolAttribute{
				Computed:    true,
				Description: "If `true`, the user account is an internal system account for the TrueNAS server. Typically, one sho",
			},
			"smb": schema.BoolAttribute{
				Computed:    true,
				Description: "The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash o",
			},
			"userns_idmap": schema.Int64Attribute{
				Computed:    true,
				Description: "Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to al",
			},
			"group": schema.StringAttribute{
				Computed:    true,
				Description: "The primary group of the user account. ",
			},
			"groups": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of additional groups to which the user belongs. NOTE: Groups are identified by their group ent",
			},
			"password_disabled": schema.BoolAttribute{
				Computed:    true,
				Description: "If set to `true` password authentication for the user account is disabled.  NOTE: Users with passwor",
			},
			"ssh_password_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Allow the user to authenticate to the TrueNAS SSH server using a password.  WARNING: The established",
			},
			"sshpubkey": schema.StringAttribute{
				Computed:    true,
				Description: "SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server.",
			},
			"locked": schema.BoolAttribute{
				Computed:    true,
				Description: "If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS se",
			},
			"sudo_commands": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "An array of commands the user may execute with elevated privileges. User is prompted for password   ",
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "An array of commands the user may execute with elevated privileges. User is *not* prompted for passw",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and",
			},
			"local": schema.BoolAttribute{
				Computed:    true,
				Description: "If `true`, the account is local to the TrueNAS server. If `false`, the account is provided by a dire",
			},
			"immutable": schema.BoolAttribute{
				Computed:    true,
				Description: "If `true`, the account is system-provided and most fields related to it may not be changed. ",
			},
			"twofactor_auth_configured": schema.BoolAttribute{
				Computed:    true,
				Description: "If `true`, the account has been configured for two-factor authentication. Users are prompted for a  ",
			},
			"sid": schema.StringAttribute{
				Computed:    true,
				Description: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses     t",
			},
			"last_password_change": schema.StringAttribute{
				Computed:    true,
				Description: "The date of the last password change for local user accounts.",
			},
			"password_age": schema.Int64Attribute{
				Computed:    true,
				Description: "The age in days of the password for local user accounts.",
			},
			"password_history": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "This contains hashes of the ten most recent passwords used by local user accounts, and is     for en",
			},
			"password_change_required": schema.BoolAttribute{
				Computed:    true,
				Description: "Password change for local user account is required on next login.",
			},
			"roles": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of roles assigned to this user's groups. Roles control administrative access to TrueNAS throug",
			},
			"api_keys": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of API key IDs associated with this user account for programmatic access.",
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("user.get_instance", func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read user: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["uid"]; ok {
		switch val := v.(type) {
		case float64:
			data.Uid = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Uid = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["username"]; ok {
		switch val := v.(type) {
		case string:
			data.Username = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Username = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["unixhash"]; ok {
		if v == nil {
			data.Unixhash = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Unixhash = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Unixhash = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Unixhash = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["smbhash"]; ok {
		if v == nil {
			data.Smbhash = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Smbhash = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Smbhash = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Smbhash = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["home"]; ok {
		switch val := v.(type) {
		case string:
			data.Home = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Home = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Home = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["shell"]; ok {
		switch val := v.(type) {
		case string:
			data.Shell = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Shell = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Shell = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["full_name"]; ok {
		switch val := v.(type) {
		case string:
			data.FullName = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.FullName = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.FullName = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["builtin"]; ok {
		if bv, ok := v.(bool); ok {
			data.Builtin = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["smb"]; ok {
		if bv, ok := v.(bool); ok {
			data.Smb = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["userns_idmap"]; ok {
		switch val := v.(type) {
		case float64:
			data.UsernsIdmap = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.UsernsIdmap = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["group"]; ok {
		switch val := v.(type) {
		case string:
			data.Group = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Group = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Group = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["groups"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.Groups, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["password_disabled"]; ok {
		if bv, ok := v.(bool); ok {
			data.PasswordDisabled = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["ssh_password_enabled"]; ok {
		if bv, ok := v.(bool); ok {
			data.SshPasswordEnabled = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["sshpubkey"]; ok {
		if v == nil {
			data.Sshpubkey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Sshpubkey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Sshpubkey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Sshpubkey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["locked"]; ok {
		if bv, ok := v.(bool); ok {
			data.Locked = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["sudo_commands"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.SudoCommands, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["sudo_commands_nopasswd"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.SudoCommandsNopasswd, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["email"]; ok {
		if v == nil {
			data.Email = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Email = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Email = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Email = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["local"]; ok {
		if bv, ok := v.(bool); ok {
			data.Local = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["immutable"]; ok {
		if bv, ok := v.(bool); ok {
			data.Immutable = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["twofactor_auth_configured"]; ok {
		if bv, ok := v.(bool); ok {
			data.TwofactorAuthConfigured = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["sid"]; ok {
		if v == nil {
			data.Sid = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Sid = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Sid = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Sid = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["last_password_change"]; ok {
		if v == nil {
			data.LastPasswordChange = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LastPasswordChange = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LastPasswordChange = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LastPasswordChange = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["password_age"]; ok {
		if v == nil {
			data.PasswordAge = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.PasswordAge = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.PasswordAge = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["password_history"]; ok {
		if v == nil {
		} else {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr {
					strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
				}
				data.PasswordHistory, _ = types.ListValue(types.StringType, strVals)
			}
		}
	}
	if v, ok := resultMap["password_change_required"]; ok {
		if bv, ok := v.(bool); ok {
			data.PasswordChangeRequired = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["roles"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.Roles, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["api_keys"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.ApiKeys, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if data.Uid.IsUnknown() {
		data.Uid = types.Int64Null()
	}
	if data.Username.IsUnknown() {
		data.Username = types.StringNull()
	}
	if data.Unixhash.IsUnknown() {
		data.Unixhash = types.StringNull()
	}
	if data.Smbhash.IsUnknown() {
		data.Smbhash = types.StringNull()
	}
	if data.Home.IsUnknown() {
		data.Home = types.StringNull()
	}
	if data.Shell.IsUnknown() {
		data.Shell = types.StringNull()
	}
	if data.FullName.IsUnknown() {
		data.FullName = types.StringNull()
	}
	if data.Builtin.IsUnknown() {
		data.Builtin = types.BoolNull()
	}
	if data.Smb.IsUnknown() {
		data.Smb = types.BoolNull()
	}
	if data.UsernsIdmap.IsUnknown() {
		data.UsernsIdmap = types.Int64Null()
	}
	if data.Group.IsUnknown() {
		data.Group = types.StringNull()
	}
	if data.Groups.IsUnknown() {
		data.Groups, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.PasswordDisabled.IsUnknown() {
		data.PasswordDisabled = types.BoolNull()
	}
	if data.SshPasswordEnabled.IsUnknown() {
		data.SshPasswordEnabled = types.BoolNull()
	}
	if data.Sshpubkey.IsUnknown() {
		data.Sshpubkey = types.StringNull()
	}
	if data.Locked.IsUnknown() {
		data.Locked = types.BoolNull()
	}
	if data.SudoCommands.IsUnknown() {
		data.SudoCommands, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.SudoCommandsNopasswd.IsUnknown() {
		data.SudoCommandsNopasswd, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Email.IsUnknown() {
		data.Email = types.StringNull()
	}
	if data.Local.IsUnknown() {
		data.Local = types.BoolNull()
	}
	if data.Immutable.IsUnknown() {
		data.Immutable = types.BoolNull()
	}
	if data.TwofactorAuthConfigured.IsUnknown() {
		data.TwofactorAuthConfigured = types.BoolNull()
	}
	if data.Sid.IsUnknown() {
		data.Sid = types.StringNull()
	}
	if data.LastPasswordChange.IsUnknown() {
		data.LastPasswordChange = types.StringNull()
	}
	if data.PasswordAge.IsUnknown() {
		data.PasswordAge = types.Int64Null()
	}
	if data.PasswordHistory.IsUnknown() {
		data.PasswordHistory, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.PasswordChangeRequired.IsUnknown() {
		data.PasswordChangeRequired = types.BoolNull()
	}
	if data.Roles.IsUnknown() {
		data.Roles, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.ApiKeys.IsUnknown() {
		data.ApiKeys, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
