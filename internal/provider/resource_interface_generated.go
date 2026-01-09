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

func (r *InterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS interface resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"ipv4_dhcp": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"ipv6_auto": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"failover_critical": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"failover_group": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"failover_vhid": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"failover_aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"failover_virtual_aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"bridge_members": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"enable_learning": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"stp": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"lag_protocol": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"xmit_hash_policy": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"lacpdu_rate": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"lag_ports": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"vlan_parent_interface": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"vlan_tag": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"vlan_pcp": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"mtu": schema.Int64Attribute{
				Required: false,
				Optional: true,
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
	params["type"] = data.Type.ValueString()
	if !data.Ipv4Dhcp.IsNull() {
		params["ipv4_dhcp"] = data.Ipv4Dhcp.ValueBool()
	}
	if !data.Ipv6Auto.IsNull() {
		params["ipv6_auto"] = data.Ipv6Auto.ValueBool()
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

func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InterfaceResourceModel
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

	_, err = r.client.Call("interface.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
	params["type"] = data.Type.ValueString()
	if !data.Ipv4Dhcp.IsNull() {
		params["ipv4_dhcp"] = data.Ipv4Dhcp.ValueBool()
	}
	if !data.Ipv6Auto.IsNull() {
		params["ipv6_auto"] = data.Ipv6Auto.ValueBool()
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

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("interface.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceResourceModel
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

	_, err = r.client.Call("interface.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
