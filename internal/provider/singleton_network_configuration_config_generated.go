package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Network_ConfigurationConfigResource struct {
	client *client.Client
}

type Network_ConfigurationConfigResourceModel struct {
	Hostname types.String `tfsdk:"hostname"`
	Domain types.String `tfsdk:"domain"`
	Ipv4Gateway types.String `tfsdk:"ipv4gateway"`
	Ipv6Gateway types.String `tfsdk:"ipv6gateway"`
	Nameserver1 types.String `tfsdk:"nameserver1"`
	Nameserver2 types.String `tfsdk:"nameserver2"`
	Nameserver3 types.String `tfsdk:"nameserver3"`
	Httpproxy types.String `tfsdk:"httpproxy"`
	ServiceAnnouncement types.String `tfsdk:"service_announcement"`
	Activity types.String `tfsdk:"activity"`
	HostnameB types.String `tfsdk:"hostname_b"`
	HostnameVirtual types.String `tfsdk:"hostname_virtual"`
}

func NewNetwork_ConfigurationConfigResource() resource.Resource {
	return &Network_ConfigurationConfigResource{}
}

func (r *Network_ConfigurationConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_configuration_config"
}

func (r *Network_ConfigurationConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Network configuration",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "System hostname."},
			"domain": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "System domain name."},
			"ipv4gateway": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Used instead of the default gateway provided by DHCP."},
			"ipv6gateway": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "IPv6 default gateway address."},
			"nameserver1": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Primary DNS server."},
			"nameserver2": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Secondary DNS server."},
			"nameserver3": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Tertiary DNS server."},
			"httpproxy": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Must be provided if a proxy is to be used for network operations."},
			"service_announcement": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Determines the broadcast protocols that will be used to advertise the server."},
			"activity": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Network activity filtering configuration."},
			"hostname_b": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Hostname for the second controller in HA configurations or `null`."},
			"hostname_virtual": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Virtual hostname for HA configurations or `null`."},
		},
	}
}

func (r *Network_ConfigurationConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Network_ConfigurationConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data Network_ConfigurationConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Hostname.IsNull() && !data.Hostname.IsUnknown() { params["hostname"] = data.Hostname.ValueString() }
	if !data.Domain.IsNull() && !data.Domain.IsUnknown() { params["domain"] = data.Domain.ValueString() }
	if !data.Ipv4Gateway.IsNull() && !data.Ipv4Gateway.IsUnknown() { params["ipv4gateway"] = data.Ipv4Gateway.ValueString() }
	if !data.Ipv6Gateway.IsNull() && !data.Ipv6Gateway.IsUnknown() { params["ipv6gateway"] = data.Ipv6Gateway.ValueString() }
	if !data.Nameserver1.IsNull() && !data.Nameserver1.IsUnknown() { params["nameserver1"] = data.Nameserver1.ValueString() }
	if !data.Nameserver2.IsNull() && !data.Nameserver2.IsUnknown() { params["nameserver2"] = data.Nameserver2.ValueString() }
	if !data.Nameserver3.IsNull() && !data.Nameserver3.IsUnknown() { params["nameserver3"] = data.Nameserver3.ValueString() }
	if !data.Httpproxy.IsNull() && !data.Httpproxy.IsUnknown() { params["httpproxy"] = data.Httpproxy.ValueString() }
	if !data.ServiceAnnouncement.IsNull() && !data.ServiceAnnouncement.IsUnknown() { params["service_announcement"] = data.ServiceAnnouncement.ValueString() }
	if !data.Activity.IsNull() && !data.Activity.IsUnknown() { params["activity"] = data.Activity.ValueString() }
	if !data.HostnameB.IsNull() && !data.HostnameB.IsUnknown() { params["hostname_b"] = data.HostnameB.ValueString() }
	if !data.HostnameVirtual.IsNull() && !data.HostnameVirtual.IsUnknown() { params["hostname_virtual"] = data.HostnameVirtual.ValueString() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("network.configuration.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply network_configuration config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("network.configuration.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read network_configuration config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["hostname"]; ok { if v == nil { data.Hostname = types.StringNull() } else if s, ok := v.(string); ok { data.Hostname = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Hostname = types.StringValue(string(j)) } } }
		if v, ok := m["domain"]; ok { if v == nil { data.Domain = types.StringNull() } else if s, ok := v.(string); ok { data.Domain = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Domain = types.StringValue(string(j)) } } }
		if v, ok := m["ipv4gateway"]; ok { if v == nil { data.Ipv4Gateway = types.StringNull() } else if s, ok := v.(string); ok { data.Ipv4Gateway = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Ipv4Gateway = types.StringValue(string(j)) } } }
		if v, ok := m["ipv6gateway"]; ok { if v == nil { data.Ipv6Gateway = types.StringNull() } else if s, ok := v.(string); ok { data.Ipv6Gateway = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Ipv6Gateway = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver1"]; ok { if v == nil { data.Nameserver1 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver1 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver1 = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver2"]; ok { if v == nil { data.Nameserver2 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver2 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver2 = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver3"]; ok { if v == nil { data.Nameserver3 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver3 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver3 = types.StringValue(string(j)) } } }
		if v, ok := m["httpproxy"]; ok { if v == nil { data.Httpproxy = types.StringNull() } else if s, ok := v.(string); ok { data.Httpproxy = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Httpproxy = types.StringValue(string(j)) } } }
		if v, ok := m["service_announcement"]; ok { if v == nil { data.ServiceAnnouncement = types.StringNull() } else if s, ok := v.(string); ok { data.ServiceAnnouncement = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.ServiceAnnouncement = types.StringValue(string(j)) } } }
		if v, ok := m["activity"]; ok { if v == nil { data.Activity = types.StringNull() } else if s, ok := v.(string); ok { data.Activity = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Activity = types.StringValue(string(j)) } } }
		if v, ok := m["hostname_b"]; ok { if v == nil { data.HostnameB = types.StringNull() } else if s, ok := v.(string); ok { data.HostnameB = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.HostnameB = types.StringValue(string(j)) } } }
		if v, ok := m["hostname_virtual"]; ok { if v == nil { data.HostnameVirtual = types.StringNull() } else if s, ok := v.(string); ok { data.HostnameVirtual = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.HostnameVirtual = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Network_ConfigurationConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data Network_ConfigurationConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("network.configuration.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read network_configuration config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["hostname"]; ok { if v == nil { data.Hostname = types.StringNull() } else if s, ok := v.(string); ok { data.Hostname = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Hostname = types.StringValue(string(j)) } } }
		if v, ok := m["domain"]; ok { if v == nil { data.Domain = types.StringNull() } else if s, ok := v.(string); ok { data.Domain = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Domain = types.StringValue(string(j)) } } }
		if v, ok := m["ipv4gateway"]; ok { if v == nil { data.Ipv4Gateway = types.StringNull() } else if s, ok := v.(string); ok { data.Ipv4Gateway = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Ipv4Gateway = types.StringValue(string(j)) } } }
		if v, ok := m["ipv6gateway"]; ok { if v == nil { data.Ipv6Gateway = types.StringNull() } else if s, ok := v.(string); ok { data.Ipv6Gateway = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Ipv6Gateway = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver1"]; ok { if v == nil { data.Nameserver1 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver1 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver1 = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver2"]; ok { if v == nil { data.Nameserver2 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver2 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver2 = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver3"]; ok { if v == nil { data.Nameserver3 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver3 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver3 = types.StringValue(string(j)) } } }
		if v, ok := m["httpproxy"]; ok { if v == nil { data.Httpproxy = types.StringNull() } else if s, ok := v.(string); ok { data.Httpproxy = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Httpproxy = types.StringValue(string(j)) } } }
		if v, ok := m["service_announcement"]; ok { if v == nil { data.ServiceAnnouncement = types.StringNull() } else if s, ok := v.(string); ok { data.ServiceAnnouncement = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.ServiceAnnouncement = types.StringValue(string(j)) } } }
		if v, ok := m["activity"]; ok { if v == nil { data.Activity = types.StringNull() } else if s, ok := v.(string); ok { data.Activity = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Activity = types.StringValue(string(j)) } } }
		if v, ok := m["hostname_b"]; ok { if v == nil { data.HostnameB = types.StringNull() } else if s, ok := v.(string); ok { data.HostnameB = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.HostnameB = types.StringValue(string(j)) } } }
		if v, ok := m["hostname_virtual"]; ok { if v == nil { data.HostnameVirtual = types.StringNull() } else if s, ok := v.(string); ok { data.HostnameVirtual = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.HostnameVirtual = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Network_ConfigurationConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data Network_ConfigurationConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Hostname.IsNull() && !data.Hostname.IsUnknown() { params["hostname"] = data.Hostname.ValueString() }
	if !data.Domain.IsNull() && !data.Domain.IsUnknown() { params["domain"] = data.Domain.ValueString() }
	if !data.Ipv4Gateway.IsNull() && !data.Ipv4Gateway.IsUnknown() { params["ipv4gateway"] = data.Ipv4Gateway.ValueString() }
	if !data.Ipv6Gateway.IsNull() && !data.Ipv6Gateway.IsUnknown() { params["ipv6gateway"] = data.Ipv6Gateway.ValueString() }
	if !data.Nameserver1.IsNull() && !data.Nameserver1.IsUnknown() { params["nameserver1"] = data.Nameserver1.ValueString() }
	if !data.Nameserver2.IsNull() && !data.Nameserver2.IsUnknown() { params["nameserver2"] = data.Nameserver2.ValueString() }
	if !data.Nameserver3.IsNull() && !data.Nameserver3.IsUnknown() { params["nameserver3"] = data.Nameserver3.ValueString() }
	if !data.Httpproxy.IsNull() && !data.Httpproxy.IsUnknown() { params["httpproxy"] = data.Httpproxy.ValueString() }
	if !data.ServiceAnnouncement.IsNull() && !data.ServiceAnnouncement.IsUnknown() { params["service_announcement"] = data.ServiceAnnouncement.ValueString() }
	if !data.Activity.IsNull() && !data.Activity.IsUnknown() { params["activity"] = data.Activity.ValueString() }
	if !data.HostnameB.IsNull() && !data.HostnameB.IsUnknown() { params["hostname_b"] = data.HostnameB.ValueString() }
	if !data.HostnameVirtual.IsNull() && !data.HostnameVirtual.IsUnknown() { params["hostname_virtual"] = data.HostnameVirtual.ValueString() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("network.configuration.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update network_configuration config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("network.configuration.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read network_configuration config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["hostname"]; ok { if v == nil { data.Hostname = types.StringNull() } else if s, ok := v.(string); ok { data.Hostname = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Hostname = types.StringValue(string(j)) } } }
		if v, ok := m["domain"]; ok { if v == nil { data.Domain = types.StringNull() } else if s, ok := v.(string); ok { data.Domain = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Domain = types.StringValue(string(j)) } } }
		if v, ok := m["ipv4gateway"]; ok { if v == nil { data.Ipv4Gateway = types.StringNull() } else if s, ok := v.(string); ok { data.Ipv4Gateway = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Ipv4Gateway = types.StringValue(string(j)) } } }
		if v, ok := m["ipv6gateway"]; ok { if v == nil { data.Ipv6Gateway = types.StringNull() } else if s, ok := v.(string); ok { data.Ipv6Gateway = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Ipv6Gateway = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver1"]; ok { if v == nil { data.Nameserver1 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver1 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver1 = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver2"]; ok { if v == nil { data.Nameserver2 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver2 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver2 = types.StringValue(string(j)) } } }
		if v, ok := m["nameserver3"]; ok { if v == nil { data.Nameserver3 = types.StringNull() } else if s, ok := v.(string); ok { data.Nameserver3 = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Nameserver3 = types.StringValue(string(j)) } } }
		if v, ok := m["httpproxy"]; ok { if v == nil { data.Httpproxy = types.StringNull() } else if s, ok := v.(string); ok { data.Httpproxy = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Httpproxy = types.StringValue(string(j)) } } }
		if v, ok := m["service_announcement"]; ok { if v == nil { data.ServiceAnnouncement = types.StringNull() } else if s, ok := v.(string); ok { data.ServiceAnnouncement = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.ServiceAnnouncement = types.StringValue(string(j)) } } }
		if v, ok := m["activity"]; ok { if v == nil { data.Activity = types.StringNull() } else if s, ok := v.(string); ok { data.Activity = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Activity = types.StringValue(string(j)) } } }
		if v, ok := m["hostname_b"]; ok { if v == nil { data.HostnameB = types.StringNull() } else if s, ok := v.(string); ok { data.HostnameB = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.HostnameB = types.StringValue(string(j)) } } }
		if v, ok := m["hostname_virtual"]; ok { if v == nil { data.HostnameVirtual = types.StringNull() } else if s, ok := v.(string); ok { data.HostnameVirtual = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.HostnameVirtual = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Network_ConfigurationConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
