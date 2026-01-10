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
				Required: false,
				Optional: true,
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Ipv4Dhcp.IsNull() {
		params["ipv4_dhcp"] = data.Ipv4Dhcp.ValueBool()
	}
	if !data.Ipv6Auto.IsNull() {
		params["ipv6_auto"] = data.Ipv6Auto.ValueBool()
	}
	if !data.Aliases.IsNull() {
		var aliasesList []string
		data.Aliases.ElementsAs(ctx, &aliasesList, false)
		params["aliases"] = aliasesList
	}
	if !data.FailoverCritical.IsNull() {
		params["failover_critical"] = data.FailoverCritical.ValueBool()
	}
	if !data.FailoverGroup.IsNull() {
		params["failover_group"] = data.FailoverGroup.ValueInt64()
	}
	if !data.FailoverVhid.IsNull() {
		params["failover_vhid"] = data.FailoverVhid.ValueInt64()
	}
	if !data.FailoverAliases.IsNull() {
		var failover_aliasesList []string
		data.FailoverAliases.ElementsAs(ctx, &failover_aliasesList, false)
		params["failover_aliases"] = failover_aliasesList
	}
	if !data.FailoverVirtualAliases.IsNull() {
		var failover_virtual_aliasesList []string
		data.FailoverVirtualAliases.ElementsAs(ctx, &failover_virtual_aliasesList, false)
		params["failover_virtual_aliases"] = failover_virtual_aliasesList
	}
	if !data.BridgeMembers.IsNull() {
		var bridge_membersList []string
		data.BridgeMembers.ElementsAs(ctx, &bridge_membersList, false)
		params["bridge_members"] = bridge_membersList
	}
	if !data.EnableLearning.IsNull() {
		params["enable_learning"] = data.EnableLearning.ValueBool()
	}
	if !data.Stp.IsNull() {
		params["stp"] = data.Stp.ValueBool()
	}
	if !data.LagProtocol.IsNull() {
		params["lag_protocol"] = data.LagProtocol.ValueString()
	}
	if !data.XmitHashPolicy.IsNull() {
		params["xmit_hash_policy"] = data.XmitHashPolicy.ValueString()
	}
	if !data.LacpduRate.IsNull() {
		params["lacpdu_rate"] = data.LacpduRate.ValueString()
	}
	if !data.LagPorts.IsNull() {
		var lag_portsList []string
		data.LagPorts.ElementsAs(ctx, &lag_portsList, false)
		params["lag_ports"] = lag_portsList
	}
	if !data.VlanParentInterface.IsNull() {
		params["vlan_parent_interface"] = data.VlanParentInterface.ValueString()
	}
	if !data.VlanTag.IsNull() {
		params["vlan_tag"] = data.VlanTag.ValueInt64()
	}
	if !data.VlanPcp.IsNull() {
		params["vlan_pcp"] = data.VlanPcp.ValueInt64()
	}
	if !data.Mtu.IsNull() {
		params["mtu"] = data.Mtu.ValueInt64()
	}

	result, err := r.client.Call("interface.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create interface: %s", err))
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

func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("interface.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read interface: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ipv4_dhcp"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Ipv4Dhcp = types.BoolValue(bv) }
		}
		if v, ok := resultMap["ipv6_auto"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Ipv6Auto = types.BoolValue(bv) }
		}
		if v, ok := resultMap["aliases"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Aliases, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["failover_critical"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.FailoverCritical = types.BoolValue(bv) }
		}
		if v, ok := resultMap["failover_group"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.FailoverGroup = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["failover_vhid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.FailoverVhid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["failover_aliases"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.FailoverAliases, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["failover_virtual_aliases"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.FailoverVirtualAliases, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["bridge_members"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.BridgeMembers, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["enable_learning"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.EnableLearning = types.BoolValue(bv) }
		}
		if v, ok := resultMap["stp"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Stp = types.BoolValue(bv) }
		}
		if v, ok := resultMap["lag_protocol"]; ok && v != nil {
			data.LagProtocol = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["xmit_hash_policy"]; ok && v != nil {
			data.XmitHashPolicy = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lacpdu_rate"]; ok && v != nil {
			data.LacpduRate = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lag_ports"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.LagPorts, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["vlan_parent_interface"]; ok && v != nil {
			data.VlanParentInterface = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vlan_tag"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.VlanTag = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["vlan_pcp"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.VlanPcp = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["mtu"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Mtu = types.Int64Value(int64(fv)) }
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

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Ipv4Dhcp.IsNull() {
		params["ipv4_dhcp"] = data.Ipv4Dhcp.ValueBool()
	}
	if !data.Ipv6Auto.IsNull() {
		params["ipv6_auto"] = data.Ipv6Auto.ValueBool()
	}
	if !data.Aliases.IsNull() {
		var aliasesList []string
		data.Aliases.ElementsAs(ctx, &aliasesList, false)
		params["aliases"] = aliasesList
	}
	if !data.FailoverCritical.IsNull() {
		params["failover_critical"] = data.FailoverCritical.ValueBool()
	}
	if !data.FailoverGroup.IsNull() {
		params["failover_group"] = data.FailoverGroup.ValueInt64()
	}
	if !data.FailoverVhid.IsNull() {
		params["failover_vhid"] = data.FailoverVhid.ValueInt64()
	}
	if !data.FailoverAliases.IsNull() {
		var failover_aliasesList []string
		data.FailoverAliases.ElementsAs(ctx, &failover_aliasesList, false)
		params["failover_aliases"] = failover_aliasesList
	}
	if !data.FailoverVirtualAliases.IsNull() {
		var failover_virtual_aliasesList []string
		data.FailoverVirtualAliases.ElementsAs(ctx, &failover_virtual_aliasesList, false)
		params["failover_virtual_aliases"] = failover_virtual_aliasesList
	}
	if !data.BridgeMembers.IsNull() {
		var bridge_membersList []string
		data.BridgeMembers.ElementsAs(ctx, &bridge_membersList, false)
		params["bridge_members"] = bridge_membersList
	}
	if !data.EnableLearning.IsNull() {
		params["enable_learning"] = data.EnableLearning.ValueBool()
	}
	if !data.Stp.IsNull() {
		params["stp"] = data.Stp.ValueBool()
	}
	if !data.LagProtocol.IsNull() {
		params["lag_protocol"] = data.LagProtocol.ValueString()
	}
	if !data.XmitHashPolicy.IsNull() {
		params["xmit_hash_policy"] = data.XmitHashPolicy.ValueString()
	}
	if !data.LacpduRate.IsNull() {
		params["lacpdu_rate"] = data.LacpduRate.ValueString()
	}
	if !data.LagPorts.IsNull() {
		var lag_portsList []string
		data.LagPorts.ElementsAs(ctx, &lag_portsList, false)
		params["lag_ports"] = lag_portsList
	}
	if !data.VlanParentInterface.IsNull() {
		params["vlan_parent_interface"] = data.VlanParentInterface.ValueString()
	}
	if !data.VlanTag.IsNull() {
		params["vlan_tag"] = data.VlanTag.ValueInt64()
	}
	if !data.VlanPcp.IsNull() {
		params["vlan_pcp"] = data.VlanPcp.ValueInt64()
	}
	if !data.Mtu.IsNull() {
		params["mtu"] = data.Mtu.ValueInt64()
	}

	_, err = r.client.Call("interface.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update interface: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("interface.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete interface: %s", err))
		return
	}
}
