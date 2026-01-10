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

var _ datasource.DataSource = &InterfacesDataSource{}

func NewInterfacesDataSource() datasource.DataSource {
	return &InterfacesDataSource{}
}

type InterfacesDataSource struct {
	client *client.Client
}

type InterfacesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type InterfacesItemModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Fake types.Bool `tfsdk:"fake"`
	Type types.String `tfsdk:"type"`
	State types.String `tfsdk:"state"`
	Ipv4Dhcp types.Bool `tfsdk:"ipv4_dhcp"`
	Ipv6Auto types.Bool `tfsdk:"ipv6_auto"`
	Description types.String `tfsdk:"description"`
	Mtu types.Int64 `tfsdk:"mtu"`
	VlanParentInterface types.String `tfsdk:"vlan_parent_interface"`
	VlanTag types.Int64 `tfsdk:"vlan_tag"`
	VlanPcp types.Int64 `tfsdk:"vlan_pcp"`
	LagProtocol types.String `tfsdk:"lag_protocol"`
	EnableLearning types.Bool `tfsdk:"enable_learning"`
}

func (d *InterfacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interfaces"
}

func (d *InterfacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query Interfaces with `query-filters` and `query-options`",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of interfaces resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the network interface.",
			},
			"fake": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this is a fake/simulated interface for testing purposes.",
			},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "Type of interface (PHYSICAL, BRIDGE, LINK_AGGREGATION, VLAN, etc.).",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Current runtime state information for the interface.",
			},
			"ipv4_dhcp": schema.BoolAttribute{
				Computed: true,
				Description: "Whether IPv4 DHCP is enabled for automatic IP address assignment.",
			},
			"ipv6_auto": schema.BoolAttribute{
				Computed: true,
				Description: "Whether IPv6 autoconfiguration is enabled.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable description of the interface.",
			},
			"mtu": schema.Int64Attribute{
				Computed: true,
				Description: "Maximum transmission unit size for the interface.",
			},
			"vlan_parent_interface": schema.StringAttribute{
				Computed: true,
				Description: "Parent interface for VLAN configuration.",
			},
			"vlan_tag": schema.Int64Attribute{
				Computed: true,
				Description: "VLAN tag number for VLAN interfaces.",
			},
			"vlan_pcp": schema.Int64Attribute{
				Computed: true,
				Description: "Priority Code Point for VLAN traffic prioritization.",
			},
			"lag_protocol": schema.StringAttribute{
				Computed: true,
				Description: "Link aggregation protocol (LACP, FAILOVER, LOADBALANCE, etc.).",
			},
			"enable_learning": schema.BoolAttribute{
				Computed: true,
				Description: "Whether MAC address learning is enabled for bridge interfaces.",
			},
					},
				},
			},
		},
	}
}

func (d *InterfacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InterfacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfacesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("interface.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query interfaces: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]InterfacesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := InterfacesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["fake"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Fake = types.BoolValue(bv) }
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ipv4_dhcp"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Ipv4Dhcp = types.BoolValue(bv) }
		}
		if v, ok := resultMap["ipv6_auto"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Ipv6Auto = types.BoolValue(bv) }
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mtu"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Mtu = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["vlan_parent_interface"]; ok && v != nil {
			itemModel.VlanParentInterface = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vlan_tag"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.VlanTag = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["vlan_pcp"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.VlanPcp = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["lag_protocol"]; ok && v != nil {
			itemModel.LagProtocol = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enable_learning"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.EnableLearning = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"enable_learning": types.BoolType,
			"fake": types.BoolType,
			"id": types.StringType,
			"ipv4_dhcp": types.BoolType,
			"ipv6_auto": types.BoolType,
			"lag_protocol": types.StringType,
			"mtu": types.Int64Type,
			"name": types.StringType,
			"state": types.StringType,
			"type": types.StringType,
			"vlan_parent_interface": types.StringType,
			"vlan_pcp": types.Int64Type,
			"vlan_tag": types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
