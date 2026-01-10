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

var _ datasource.DataSource = &GroupsDataSource{}

func NewGroupsDataSource() datasource.DataSource {
	return &GroupsDataSource{}
}

type GroupsDataSource struct {
	client *client.Client
}

type GroupsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type GroupsItemModel struct {
	ID types.String `tfsdk:"id"`
	Gid types.Int64 `tfsdk:"gid"`
	Name types.String `tfsdk:"name"`
	Builtin types.Bool `tfsdk:"builtin"`
	Smb types.Bool `tfsdk:"smb"`
	UsernsIdmap types.Int64 `tfsdk:"userns_idmap"`
	Group types.String `tfsdk:"group"`
	Local types.Bool `tfsdk:"local"`
	Sid types.String `tfsdk:"sid"`
	Immutable types.Bool `tfsdk:"immutable"`
}

func (d *GroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *GroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query groups with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of groups resources",
				NestedObject: schema.NestedAttributeObject{
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
			"immutable": schema.BoolAttribute{
				Computed: true,
				Description: "This is a read-only field showing if the group entry can be changed. If `True`, the group is immutab",
			},
					},
				},
			},
		},
	}
}

func (d *GroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("group.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query groups: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]GroupsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := GroupsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["gid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Gid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
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
		if v, ok := resultMap["local"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Local = types.BoolValue(bv) }
		}
		if v, ok := resultMap["sid"]; ok && v != nil {
			itemModel.Sid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["immutable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Immutable = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"builtin": types.BoolType,
			"gid": types.Int64Type,
			"group": types.StringType,
			"id": types.StringType,
			"immutable": types.BoolType,
			"local": types.BoolType,
			"name": types.StringType,
			"sid": types.StringType,
			"smb": types.BoolType,
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
