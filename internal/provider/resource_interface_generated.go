package provider

import (
	"context"
	"fmt"
	"strings"
"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InterfaceResource struct {
	client *client.Client
}

type InterfaceResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type types.String `tfsdk:"type"`
	Ipv4Dhcp types.Bool `tfsdk:"ipv4_dhcp"`
	Ipv6Auto types.Bool `tfsdk:"ipv6_auto"`
	Aliases types.List `tfsdk:"aliases"`
	FailoverCritical types.Bool `tfsdk:"failover_critical"`
	FailoverGroup types.Int64 `tfsdk:"failover_group"`
	FailoverVhid types.Int64 `tfsdk:"failover_vhid"`
	FailoverAliases types.List `tfsdk:"failover_aliases"`
	FailoverVirtualAliases types.List `tfsdk:"failover_virtual_aliases"`
	BridgeMembers types.List `tfsdk:"bridge_members"`
	EnableLearning types.Bool `tfsdk:"enable_learning"`
	Stp types.Bool `tfsdk:"stp"`
	LagProtocol types.String `tfsdk:"lag_protocol"`
	XmitHashPolicy types.String `tfsdk:"xmit_hash_policy"`
	LacpduRate types.String `tfsdk:"lacpdu_rate"`
	LagPorts types.List `tfsdk:"lag_ports"`
	VlanParentInterface types.String `tfsdk:"vlan_parent_interface"`
	VlanTag types.Int64 `tfsdk:"vlan_tag"`
	VlanPcp types.Int64 `tfsdk:"vlan_pcp"`
	Mtu types.Int64 `tfsdk:"mtu"`
}

func NewInterfaceResource() resource.Resource {
	return &InterfaceResource{}
}

func (r *InterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface"
}

func (r *InterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *InterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create virtual interfaces (Link Aggregation, VLAN)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Generate a name if not provided based on `type`, e.g. \"br0\", \"bond1\", \"vlan0\".",
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Human-readable description of the interface.",
			},
			"type": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Type of interface to create.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ipv4_dhcp": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable IPv4 DHCP for automatic IP address assignment.",
			},
			"ipv6_auto": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable IPv6 autoconfiguration.",
			},
			"aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of IP address aliases to configure on the interface.",
			},
			"failover_critical": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether this interface is critical for failover functionality. Critical interfaces are monitored for",
			},
			"failover_group": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Failover group identifier for clustering. Interfaces in the same group fail over together during    ",
			},
			"failover_vhid": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Virtual Host ID for VRRP failover configuration. Must be unique within the VRRP group and match     ",
			},
			"failover_aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of IP aliases for failover configuration. These IPs are assigned to the interface during normal",
			},
			"failover_virtual_aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of virtual IP aliases for failover configuration. These are shared IPs that float between nodes",
			},
			"bridge_members": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of interfaces to add as members of this bridge.",
			},
			"enable_learning": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable MAC address learning for bridge interfaces. When enabled, the bridge learns MAC addresses    ",
			},
			"stp": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable Spanning Tree Protocol for bridge interfaces. STP prevents network loops by blocking redundan",
			},
			"lag_protocol": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Link aggregation protocol to use for bonding interfaces. LACP uses 802.3ad dynamic negotiation,     ",
			},
			"xmit_hash_policy": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Transmit hash policy for load balancing in link aggregation. LAYER2 uses MAC addresses, LAYER2+3 add",
			},
			"lacpdu_rate": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "LACP data unit transmission rate. SLOW sends LACPDUs every 30 seconds, FAST sends every 1 second for",
			},
			"lag_ports": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of interface names to include in the link aggregation group.",
			},
			"vlan_parent_interface": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Parent interface for VLAN configuration.",
			},
			"vlan_tag": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "VLAN tag number (1-4094).",
			},
			"vlan_pcp": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Priority Code Point for VLAN traffic prioritization (0-7). Values 0-7 map to different QoS priority ",
			},
			"mtu": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum transmission unit size for the interface (68-9216 bytes).",
			},
		},
	}
}

func (r *InterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Ipv4Dhcp.IsNull() && !data.Ipv4Dhcp.IsUnknown() {
		params["ipv4_dhcp"] = data.Ipv4Dhcp.ValueBool()
	}
	if !data.Ipv6Auto.IsNull() && !data.Ipv6Auto.IsUnknown() {
		params["ipv6_auto"] = data.Ipv6Auto.ValueBool()
	}
	if !data.Aliases.IsNull() && !data.Aliases.IsUnknown() {
		var aliasesList []string
		data.Aliases.ElementsAs(ctx, &aliasesList, false)
		var aliasesObjs []map[string]interface{}
		for _, jsonStr := range aliasesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse aliases item: %s", err))
				return
			}
			aliasesObjs = append(aliasesObjs, obj)
		}
		params["aliases"] = aliasesObjs
	}
	if !data.FailoverCritical.IsNull() && !data.FailoverCritical.IsUnknown() {
		params["failover_critical"] = data.FailoverCritical.ValueBool()
	}
	if !data.FailoverGroup.IsNull() && !data.FailoverGroup.IsUnknown() {
		params["failover_group"] = data.FailoverGroup.ValueInt64()
	}
	if !data.FailoverVhid.IsNull() && !data.FailoverVhid.IsUnknown() {
		params["failover_vhid"] = data.FailoverVhid.ValueInt64()
	}
	if !data.FailoverAliases.IsNull() && !data.FailoverAliases.IsUnknown() {
		var failover_aliasesList []string
		data.FailoverAliases.ElementsAs(ctx, &failover_aliasesList, false)
		var failover_aliasesObjs []map[string]interface{}
		for _, jsonStr := range failover_aliasesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse failover_aliases item: %s", err))
				return
			}
			failover_aliasesObjs = append(failover_aliasesObjs, obj)
		}
		params["failover_aliases"] = failover_aliasesObjs
	}
	if !data.FailoverVirtualAliases.IsNull() && !data.FailoverVirtualAliases.IsUnknown() {
		var failover_virtual_aliasesList []string
		data.FailoverVirtualAliases.ElementsAs(ctx, &failover_virtual_aliasesList, false)
		var failover_virtual_aliasesObjs []map[string]interface{}
		for _, jsonStr := range failover_virtual_aliasesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse failover_virtual_aliases item: %s", err))
				return
			}
			failover_virtual_aliasesObjs = append(failover_virtual_aliasesObjs, obj)
		}
		params["failover_virtual_aliases"] = failover_virtual_aliasesObjs
	}
	if !data.BridgeMembers.IsNull() && !data.BridgeMembers.IsUnknown() {
		var bridge_membersList []string
		data.BridgeMembers.ElementsAs(ctx, &bridge_membersList, false)
		params["bridge_members"] = bridge_membersList
	}
	if !data.EnableLearning.IsNull() && !data.EnableLearning.IsUnknown() {
		params["enable_learning"] = data.EnableLearning.ValueBool()
	}
	if !data.Stp.IsNull() && !data.Stp.IsUnknown() {
		params["stp"] = data.Stp.ValueBool()
	}
	if !data.LagProtocol.IsNull() && !data.LagProtocol.IsUnknown() {
		params["lag_protocol"] = data.LagProtocol.ValueString()
	}
	if !data.XmitHashPolicy.IsNull() && !data.XmitHashPolicy.IsUnknown() {
		params["xmit_hash_policy"] = data.XmitHashPolicy.ValueString()
	}
	if !data.LacpduRate.IsNull() && !data.LacpduRate.IsUnknown() {
		params["lacpdu_rate"] = data.LacpduRate.ValueString()
	}
	if !data.LagPorts.IsNull() && !data.LagPorts.IsUnknown() {
		var lag_portsList []string
		data.LagPorts.ElementsAs(ctx, &lag_portsList, false)
		params["lag_ports"] = lag_portsList
	}
	if !data.VlanParentInterface.IsNull() && !data.VlanParentInterface.IsUnknown() {
		params["vlan_parent_interface"] = data.VlanParentInterface.ValueString()
	}
	if !data.VlanTag.IsNull() && !data.VlanTag.IsUnknown() {
		params["vlan_tag"] = data.VlanTag.ValueInt64()
	}
	if !data.VlanPcp.IsNull() && !data.VlanPcp.IsUnknown() {
		params["vlan_pcp"] = data.VlanPcp.ValueInt64()
	}
	if !data.Mtu.IsNull() && !data.Mtu.IsUnknown() {
		params["mtu"] = data.Mtu.ValueInt64()
	}

	result, err := r.client.Call("interface.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create interface: %s", err))
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
	id := data.ID.ValueString()
	result, err = r.client.Call("interface.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back interface: %s", err))
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
		if v, ok := resultMap["xmit_hash_policy"]; ok {
			switch val := v.(type) {
			case string:
				data.XmitHashPolicy = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["lacpdu_rate"]; ok {
			switch val := v.(type) {
			case string:
				data.LacpduRate = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LacpduRate = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LacpduRate = types.StringValue(fmt.Sprintf("%v", v))
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
						if fv, ok := parsed.(float64); ok { data.Mtu = types.Int64Value(int64(fv)) }
					}
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	result, err := r.client.Call("interface.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read interface: %s", err))
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
		if v, ok := resultMap["xmit_hash_policy"]; ok {
			switch val := v.(type) {
			case string:
				data.XmitHashPolicy = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["lacpdu_rate"]; ok {
			switch val := v.(type) {
			case string:
				data.LacpduRate = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LacpduRate = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LacpduRate = types.StringValue(fmt.Sprintf("%v", v))
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
						if fv, ok := parsed.(float64); ok { data.Mtu = types.Int64Value(int64(fv)) }
					}
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = state.ID.ValueString()

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Ipv4Dhcp.IsNull() && !data.Ipv4Dhcp.IsUnknown() {
		params["ipv4_dhcp"] = data.Ipv4Dhcp.ValueBool()
	}
	if !data.Ipv6Auto.IsNull() && !data.Ipv6Auto.IsUnknown() {
		params["ipv6_auto"] = data.Ipv6Auto.ValueBool()
	}
	if !data.Aliases.IsNull() && !data.Aliases.IsUnknown() {
		var aliasesList []string
		data.Aliases.ElementsAs(ctx, &aliasesList, false)
		var aliasesObjs []map[string]interface{}
		for _, jsonStr := range aliasesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse aliases item: %s", err))
				return
			}
			aliasesObjs = append(aliasesObjs, obj)
		}
		params["aliases"] = aliasesObjs
	}
	if !data.FailoverCritical.IsNull() && !data.FailoverCritical.IsUnknown() {
		params["failover_critical"] = data.FailoverCritical.ValueBool()
	}
	if !data.FailoverGroup.IsNull() && !data.FailoverGroup.IsUnknown() {
		params["failover_group"] = data.FailoverGroup.ValueInt64()
	}
	if !data.FailoverVhid.IsNull() && !data.FailoverVhid.IsUnknown() {
		params["failover_vhid"] = data.FailoverVhid.ValueInt64()
	}
	if !data.FailoverAliases.IsNull() && !data.FailoverAliases.IsUnknown() {
		var failover_aliasesList []string
		data.FailoverAliases.ElementsAs(ctx, &failover_aliasesList, false)
		var failover_aliasesObjs []map[string]interface{}
		for _, jsonStr := range failover_aliasesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse failover_aliases item: %s", err))
				return
			}
			failover_aliasesObjs = append(failover_aliasesObjs, obj)
		}
		params["failover_aliases"] = failover_aliasesObjs
	}
	if !data.FailoverVirtualAliases.IsNull() && !data.FailoverVirtualAliases.IsUnknown() {
		var failover_virtual_aliasesList []string
		data.FailoverVirtualAliases.ElementsAs(ctx, &failover_virtual_aliasesList, false)
		var failover_virtual_aliasesObjs []map[string]interface{}
		for _, jsonStr := range failover_virtual_aliasesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse failover_virtual_aliases item: %s", err))
				return
			}
			failover_virtual_aliasesObjs = append(failover_virtual_aliasesObjs, obj)
		}
		params["failover_virtual_aliases"] = failover_virtual_aliasesObjs
	}
	if !data.BridgeMembers.IsNull() && !data.BridgeMembers.IsUnknown() {
		var bridge_membersList []string
		data.BridgeMembers.ElementsAs(ctx, &bridge_membersList, false)
		params["bridge_members"] = bridge_membersList
	}
	if !data.EnableLearning.IsNull() && !data.EnableLearning.IsUnknown() {
		params["enable_learning"] = data.EnableLearning.ValueBool()
	}
	if !data.Stp.IsNull() && !data.Stp.IsUnknown() {
		params["stp"] = data.Stp.ValueBool()
	}
	if !data.LagProtocol.IsNull() && !data.LagProtocol.IsUnknown() {
		params["lag_protocol"] = data.LagProtocol.ValueString()
	}
	if !data.XmitHashPolicy.IsNull() && !data.XmitHashPolicy.IsUnknown() {
		params["xmit_hash_policy"] = data.XmitHashPolicy.ValueString()
	}
	if !data.LacpduRate.IsNull() && !data.LacpduRate.IsUnknown() {
		params["lacpdu_rate"] = data.LacpduRate.ValueString()
	}
	if !data.LagPorts.IsNull() && !data.LagPorts.IsUnknown() {
		var lag_portsList []string
		data.LagPorts.ElementsAs(ctx, &lag_portsList, false)
		params["lag_ports"] = lag_portsList
	}
	if !data.VlanParentInterface.IsNull() && !data.VlanParentInterface.IsUnknown() {
		params["vlan_parent_interface"] = data.VlanParentInterface.ValueString()
	}
	if !data.VlanTag.IsNull() && !data.VlanTag.IsUnknown() {
		params["vlan_tag"] = data.VlanTag.ValueInt64()
	}
	if !data.VlanPcp.IsNull() && !data.VlanPcp.IsUnknown() {
		params["vlan_pcp"] = data.VlanPcp.ValueInt64()
	}
	if !data.Mtu.IsNull() && !data.Mtu.IsUnknown() {
		params["mtu"] = data.Mtu.ValueInt64()
	}

	_, err = r.client.Call("interface.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update interface: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("interface.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back interface: %s", readErr))
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
		if v, ok := resultMap["xmit_hash_policy"]; ok {
			switch val := v.(type) {
			case string:
				data.XmitHashPolicy = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["lacpdu_rate"]; ok {
			switch val := v.(type) {
			case string:
				data.LacpduRate = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LacpduRate = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LacpduRate = types.StringValue(fmt.Sprintf("%v", v))
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
						if fv, ok := parsed.(float64); ok { data.Mtu = types.Int64Value(int64(fv)) }
					}
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	_, err := r.client.Call("interface.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete interface: %s", err))
		return
	}
}
