package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &UsersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

type UsersDataSource struct {
	client *client.Client
}

type UsersDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type UsersItemModel struct {
	ID types.String `tfsdk:"id"`
	Uid types.Int64 `tfsdk:"uid"`
	Username types.String `tfsdk:"username"`
	Unixhash types.String `tfsdk:"unixhash"`
	Smbhash types.String `tfsdk:"smbhash"`
	Home types.String `tfsdk:"home"`
	Shell types.String `tfsdk:"shell"`
	FullName types.String `tfsdk:"full_name"`
	Builtin types.Bool `tfsdk:"builtin"`
	Smb types.Bool `tfsdk:"smb"`
	UsernsIdmap types.Int64 `tfsdk:"userns_idmap"`
	Group types.String `tfsdk:"group"`
	PasswordDisabled types.Bool `tfsdk:"password_disabled"`
	SshPasswordEnabled types.Bool `tfsdk:"ssh_password_enabled"`
	Sshpubkey types.String `tfsdk:"sshpubkey"`
	Locked types.Bool `tfsdk:"locked"`
	Email types.String `tfsdk:"email"`
	Local types.Bool `tfsdk:"local"`
	Immutable types.Bool `tfsdk:"immutable"`
	TwofactorAuthConfigured types.Bool `tfsdk:"twofactor_auth_configured"`
	Sid types.String `tfsdk:"sid"`
	LastPasswordChange types.String `tfsdk:"last_password_change"`
	PasswordAge types.Int64 `tfsdk:"password_age"`
	PasswordHistory types.String `tfsdk:"password_history"`
	PasswordChangeRequired types.Bool `tfsdk:"password_change_required"`
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query users with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of users resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"uid": schema.Int64Attribute{
				Computed: true,
				Description: "A non-negative integer used to identify a system user. TrueNAS uses this value for permission     ch",
			},
			"username": schema.StringAttribute{
				Computed: true,
				Description: "A string used to identify a user. Local accounts must use characters from the POSIX portable filenam",
			},
			"unixhash": schema.StringAttribute{
				Computed: true,
				Description: "Hashed password for local accounts. This value is `null` for accounts provided by directory services",
			},
			"smbhash": schema.StringAttribute{
				Computed: true,
				Description: "NT hash of the local account password for `smb` users. This value is `null` for accounts provided by",
			},
			"home": schema.StringAttribute{
				Computed: true,
				Description: "The local file system path for the user account's home directory. Typically, this is required only i",
			},
			"shell": schema.StringAttribute{
				Computed: true,
				Description: "Available choices can be retrieved with `user.shell_choices`.",
			},
			"full_name": schema.StringAttribute{
				Computed: true,
				Description: "Comment field to provide additional information about the user account. Typically, this is     the f",
			},
			"builtin": schema.BoolAttribute{
				Computed: true,
				Description: "If `true`, the user account is an internal system account for the TrueNAS server. Typically, one sho",
			},
			"smb": schema.BoolAttribute{
				Computed: true,
				Description: "The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash o",
			},
			"userns_idmap": schema.Int64Attribute{
				Computed: true,
				Description: "Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to al",
			},
			"group": schema.StringAttribute{
				Computed: true,
				Description: "The primary group of the user account. ",
			},
			"password_disabled": schema.BoolAttribute{
				Computed: true,
				Description: "If set to `true` password authentication for the user account is disabled.  NOTE: Users with passwor",
			},
			"ssh_password_enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Allow the user to authenticate to the TrueNAS SSH server using a password.  WARNING: The established",
			},
			"sshpubkey": schema.StringAttribute{
				Computed: true,
				Description: "SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server.",
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				Description: "If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS se",
			},
			"email": schema.StringAttribute{
				Computed: true,
				Description: "Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and",
			},
			"local": schema.BoolAttribute{
				Computed: true,
				Description: "If `true`, the account is local to the TrueNAS server. If `false`, the account is provided by a dire",
			},
			"immutable": schema.BoolAttribute{
				Computed: true,
				Description: "If `true`, the account is system-provided and most fields related to it may not be changed. ",
			},
			"twofactor_auth_configured": schema.BoolAttribute{
				Computed: true,
				Description: "If `true`, the account has been configured for two-factor authentication. Users are prompted for a  ",
			},
			"sid": schema.StringAttribute{
				Computed: true,
				Description: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses     t",
			},
			"last_password_change": schema.StringAttribute{
				Computed: true,
				Description: "The date of the last password change for local user accounts.",
			},
			"password_age": schema.Int64Attribute{
				Computed: true,
				Description: "The age in days of the password for local user accounts.",
			},
			"password_history": schema.StringAttribute{
				Computed: true,
				Description: "This contains hashes of the ten most recent passwords used by local user accounts, and is     for en",
			},
			"password_change_required": schema.BoolAttribute{
				Computed: true,
				Description: "Password change for local user account is required on next login.",
			},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("user.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query users: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]UsersItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := UsersItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["uid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Uid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			itemModel.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["unixhash"]; ok && v != nil {
			itemModel.Unixhash = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["smbhash"]; ok && v != nil {
			itemModel.Smbhash = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["home"]; ok && v != nil {
			itemModel.Home = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["shell"]; ok && v != nil {
			itemModel.Shell = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["full_name"]; ok && v != nil {
			itemModel.FullName = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["builtin"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Builtin = types.BoolValue(bv) }
		}
		if v, ok := resultMap["smb"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Smb = types.BoolValue(bv) }
		}
		if v, ok := resultMap["userns_idmap"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.UsernsIdmap = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["group"]; ok && v != nil {
			itemModel.Group = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password_disabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.PasswordDisabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["ssh_password_enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.SshPasswordEnabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["sshpubkey"]; ok && v != nil {
			itemModel.Sshpubkey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Locked = types.BoolValue(bv) }
		}
		if v, ok := resultMap["email"]; ok && v != nil {
			itemModel.Email = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["local"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Local = types.BoolValue(bv) }
		}
		if v, ok := resultMap["immutable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Immutable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["twofactor_auth_configured"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.TwofactorAuthConfigured = types.BoolValue(bv) }
		}
		if v, ok := resultMap["sid"]; ok && v != nil {
			itemModel.Sid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["last_password_change"]; ok && v != nil {
			itemModel.LastPasswordChange = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password_age"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.PasswordAge = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["password_history"]; ok && v != nil {
			itemModel.PasswordHistory = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password_change_required"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.PasswordChangeRequired = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"builtin": types.BoolType,
			"email": types.StringType,
			"full_name": types.StringType,
			"group": types.StringType,
			"home": types.StringType,
			"id": types.StringType,
			"immutable": types.BoolType,
			"last_password_change": types.StringType,
			"local": types.BoolType,
			"locked": types.BoolType,
			"password_age": types.Int64Type,
			"password_change_required": types.BoolType,
			"password_disabled": types.BoolType,
			"password_history": types.StringType,
			"shell": types.StringType,
			"sid": types.StringType,
			"smb": types.BoolType,
			"smbhash": types.StringType,
			"ssh_password_enabled": types.BoolType,
			"sshpubkey": types.StringType,
			"twofactor_auth_configured": types.BoolType,
			"uid": types.Int64Type,
			"unixhash": types.StringType,
			"username": types.StringType,
			"userns_idmap": types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
