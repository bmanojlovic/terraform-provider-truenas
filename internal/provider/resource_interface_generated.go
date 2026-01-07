package provider

import (
	"context"
	"fmt"

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
	FailoverGroup types.String `tfsdk:"failover_group"`
	FailoverVhid types.String `tfsdk:"failover_vhid"`
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
	VlanPcp types.String `tfsdk:"vlan_pcp"`
	Mtu types.String `tfsdk:"mtu"`
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
				Required: false,
				Optional: true,
			},
			"failover_critical": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"failover_group": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"failover_vhid": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"failover_aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"failover_virtual_aliases": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"bridge_members": schema.ListAttribute{
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
			"vlan_pcp": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"mtu": schema.StringAttribute{
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

	params := map[string]interface{}{
		"name": data.Name.ValueString(),
		"description": data.Description.ValueString(),
		"type": data.Type.ValueString(),
		"ipv4_dhcp": data.Ipv4Dhcp.ValueBool(),
		"ipv6_auto": data.Ipv6Auto.ValueBool(),
		"failover_critical": data.FailoverCritical.ValueBool(),
		"failover_group": data.FailoverGroup.ValueString(),
		"failover_vhid": data.FailoverVhid.ValueString(),
		"enable_learning": data.EnableLearning.ValueBool(),
		"stp": data.Stp.ValueBool(),
		"lag_protocol": data.LagProtocol.ValueString(),
		"xmit_hash_policy": data.XmitHashPolicy.ValueString(),
		"lacpdu_rate": data.LacpduRate.ValueString(),
		"vlan_parent_interface": data.VlanParentInterface.ValueString(),
		"vlan_tag": data.VlanTag.ValueInt64(),
		"vlan_pcp": data.VlanPcp.ValueString(),
		"mtu": data.Mtu.ValueString(),
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

	_, err := r.client.Call("interface.get_instance", data.ID.ValueString())
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

	params := map[string]interface{}{
		"name": data.Name.ValueString(),
		"description": data.Description.ValueString(),
		"type": data.Type.ValueString(),
		"ipv4_dhcp": data.Ipv4Dhcp.ValueBool(),
		"ipv6_auto": data.Ipv6Auto.ValueBool(),
		"failover_critical": data.FailoverCritical.ValueBool(),
		"failover_group": data.FailoverGroup.ValueString(),
		"failover_vhid": data.FailoverVhid.ValueString(),
		"enable_learning": data.EnableLearning.ValueBool(),
		"stp": data.Stp.ValueBool(),
		"lag_protocol": data.LagProtocol.ValueString(),
		"xmit_hash_policy": data.XmitHashPolicy.ValueString(),
		"lacpdu_rate": data.LacpduRate.ValueString(),
		"vlan_parent_interface": data.VlanParentInterface.ValueString(),
		"vlan_tag": data.VlanTag.ValueInt64(),
		"vlan_pcp": data.VlanPcp.ValueString(),
		"mtu": data.Mtu.ValueString(),
	}

	_, err := r.client.Call("interface.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("interface.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
