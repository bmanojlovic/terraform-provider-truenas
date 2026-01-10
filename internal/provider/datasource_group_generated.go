package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"strconv"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &GroupDataSource{}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

type GroupDataSource struct {
	client *client.Client
}

type GroupDataSourceModel struct {
	ID types.String `tfsdk:"id"`
	Gid types.Int64 `tfsdk:"gid"`
	Name types.String `tfsdk:"name"`
	Builtin types.Bool `tfsdk:"builtin"`
	SudoCommands types.List `tfsdk:"sudo_commands"`
	SudoCommandsNopasswd types.List `tfsdk:"sudo_commands_nopasswd"`
	Smb types.Bool `tfsdk:"smb"`
	UsernsIdmap types.Int64 `tfsdk:"userns_idmap"`
	Group types.String `tfsdk:"group"`
	Local types.Bool `tfsdk:"local"`
	Sid types.String `tfsdk:"sid"`
	Roles types.List `tfsdk:"roles"`
	Users types.List `tfsdk:"users"`
	Immutable types.Bool `tfsdk:"immutable"`
}

func (d *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"gid": schema.Int64Attribute{
				Computed: true,
				Description: "A non-negative integer used to identify a group. TrueNAS uses this value for permission checks and m",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "A string used to identify a group.",
			},
			"builtin": schema.BoolAttribute{
				Computed: true,
				Description: "If `True`, the group is an internal system account for the TrueNAS server. Typically, one should    ",
			},
			"sudo_commands": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "A list of commands that group members may execute with elevated privileges. User is prompted for pas",
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "A list of commands that group members may execute with elevated privileges. User is not prompted for",
			},
			"smb": schema.BoolAttribute{
				Computed: true,
				Description: "If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT grou",
			},
			"userns_idmap": schema.Int64Attribute{
				Computed: true,
				Description: "Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to a",
			},
			"group": schema.StringAttribute{
				Computed: true,
				Description: "A string used to identify a group. Identical to the `name` key. ",
			},
			"local": schema.BoolAttribute{
				Computed: true,
				Description: "If `True`, the group is local to the TrueNAS server. If `False`, the group is provided by a director",
			},
			"sid": schema.StringAttribute{
				Computed: true,
				Description: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses this ",
			},
			"roles": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "List of roles assigned to this groups. Roles control administrative access to TrueNAS through the we",
			},
			"users": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "A list a API user identifiers for local users who are members of this group. These IDs match the `id",
			},
			"immutable": schema.BoolAttribute{
				Computed: true,
				Description: "This is a read-only field showing if the group entry can be changed. If `True`, the group is immutab",
			},
		},
	}
}

func (d *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("group.get_instance", func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read group: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["gid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Gid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["builtin"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Builtin = types.BoolValue(bv) }
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
		if v, ok := resultMap["group"]; ok && v != nil {
			data.Group = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["local"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Local = types.BoolValue(bv) }
		}
		if v, ok := resultMap["sid"]; ok && v != nil {
			data.Sid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["roles"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Roles, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["users"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Users, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["immutable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Immutable = types.BoolValue(bv) }
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
