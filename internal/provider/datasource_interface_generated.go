package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &InterfaceDataSource{}

func NewInterfaceDataSource() datasource.DataSource {
	return &InterfaceDataSource{}
}

type InterfaceDataSource struct {
	client *client.Client
}

type InterfaceDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Fake                types.Bool   `tfsdk:"fake"`
	Type                types.String `tfsdk:"type"`
	State               types.String `tfsdk:"state"`
	Aliases             types.List   `tfsdk:"aliases"`
	Ipv4Dhcp            types.Bool   `tfsdk:"ipv4_dhcp"`
	Ipv6Auto            types.Bool   `tfsdk:"ipv6_auto"`
	Description         types.String `tfsdk:"description"`
	Mtu                 types.Int64  `tfsdk:"mtu"`
	VlanParentInterface types.String `tfsdk:"vlan_parent_interface"`
	VlanTag             types.Int64  `tfsdk:"vlan_tag"`
	VlanPcp             types.Int64  `tfsdk:"vlan_pcp"`
	LagProtocol         types.String `tfsdk:"lag_protocol"`
	LagPorts            types.List   `tfsdk:"lag_ports"`
	BridgeMembers       types.List   `tfsdk:"bridge_members"`
	EnableLearning      types.Bool   `tfsdk:"enable_learning"`
}

func (d *InterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface"
}

func (d *InterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the network interface.",
			},
			"fake": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is a fake/simulated interface for testing purposes.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of interface (PHYSICAL, BRIDGE, LINK_AGGREGATION, VLAN, etc.).",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Current runtime state information for the interface.",
			},
			"aliases": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of IP address aliases configured on the interface.",
			},
			"ipv4_dhcp": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether IPv4 DHCP is enabled for automatic IP address assignment.",
			},
			"ipv6_auto": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether IPv6 autoconfiguration is enabled.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description of the interface.",
			},
			"mtu": schema.Int64Attribute{
				Computed:    true,
				Description: "Maximum transmission unit size for the interface.",
			},
			"vlan_parent_interface": schema.StringAttribute{
				Computed:    true,
				Description: "Parent interface for VLAN configuration.",
			},
			"vlan_tag": schema.Int64Attribute{
				Computed:    true,
				Description: "VLAN tag number for VLAN interfaces.",
			},
			"vlan_pcp": schema.Int64Attribute{
				Computed:    true,
				Description: "Priority Code Point for VLAN traffic prioritization.",
			},
			"lag_protocol": schema.StringAttribute{
				Computed:    true,
				Description: "Link aggregation protocol (LACP, FAILOVER, LOADBALANCE, etc.).",
			},
			"lag_ports": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of interface names that are members of this link aggregation group.",
			},
			"bridge_members": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of interface names that are members of this bridge.",
			},
			"enable_learning": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether MAC address learning is enabled for bridge interfaces.",
			},
		},
	}
}

func (d *InterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("interface.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read interface: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
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
	if v, ok := resultMap["fake"]; ok {
		if bv, ok := v.(bool); ok {
			data.Fake = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["state"]; ok {
		switch val := v.(type) {
		case string:
			data.State = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.State = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.State = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["aliases"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.Aliases, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["ipv4_dhcp"]; ok {
		if bv, ok := v.(bool); ok {
			data.Ipv4Dhcp = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["ipv6_auto"]; ok {
		if bv, ok := v.(bool); ok {
			data.Ipv6Auto = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["description"]; ok {
		switch val := v.(type) {
		case string:
			data.Description = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Description = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["mtu"]; ok {
		if v == nil {
			data.Mtu = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.Mtu = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.Mtu = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["vlan_parent_interface"]; ok {
		if v == nil {
			data.VlanParentInterface = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.VlanParentInterface = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.VlanParentInterface = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.VlanParentInterface = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["vlan_tag"]; ok {
		if v == nil {
			data.VlanTag = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.VlanTag = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.VlanTag = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["vlan_pcp"]; ok {
		if v == nil {
			data.VlanPcp = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.VlanPcp = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.VlanPcp = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["lag_protocol"]; ok {
		switch val := v.(type) {
		case string:
			data.LagProtocol = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.LagProtocol = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.LagProtocol = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["lag_ports"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.LagPorts, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["bridge_members"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.BridgeMembers, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["enable_learning"]; ok {
		if bv, ok := v.(bool); ok {
			data.EnableLearning = types.BoolValue(bv)
		}
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.Fake.IsUnknown() {
		data.Fake = types.BoolNull()
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.State.IsUnknown() {
		data.State = types.StringNull()
	}
	if data.Aliases.IsUnknown() {
		data.Aliases, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Ipv4Dhcp.IsUnknown() {
		data.Ipv4Dhcp = types.BoolNull()
	}
	if data.Ipv6Auto.IsUnknown() {
		data.Ipv6Auto = types.BoolNull()
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Mtu.IsUnknown() {
		data.Mtu = types.Int64Null()
	}
	if data.VlanParentInterface.IsUnknown() {
		data.VlanParentInterface = types.StringNull()
	}
	if data.VlanTag.IsUnknown() {
		data.VlanTag = types.Int64Null()
	}
	if data.VlanPcp.IsUnknown() {
		data.VlanPcp = types.Int64Null()
	}
	if data.LagProtocol.IsUnknown() {
		data.LagProtocol = types.StringNull()
	}
	if data.LagPorts.IsUnknown() {
		data.LagPorts, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.BridgeMembers.IsUnknown() {
		data.BridgeMembers, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.EnableLearning.IsUnknown() {
		data.EnableLearning = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
