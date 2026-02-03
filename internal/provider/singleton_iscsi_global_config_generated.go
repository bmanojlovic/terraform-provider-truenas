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

type Iscsi_GlobalConfigResource struct {
	client *client.Client
}

type Iscsi_GlobalConfigResourceModel struct {
	Basename types.String `tfsdk:"basename"`
	ListenPort types.Int64 `tfsdk:"listen_port"`
	PoolAvailThreshold types.Int64 `tfsdk:"pool_avail_threshold"`
	Alua types.Bool `tfsdk:"alua"`
	Iser types.Bool `tfsdk:"iser"`
}

func NewIscsi_GlobalConfigResource() resource.Resource {
	return &Iscsi_GlobalConfigResource{}
}

func (r *Iscsi_GlobalConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_global_config"
}

func (r *Iscsi_GlobalConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "iSCSI global configuration",
		Attributes: map[string]schema.Attribute{
			"basename": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Base name prefix for iSCSI target IQNs."},
			"listen_port": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "TCP port number for iSCSI connections."},
			"pool_avail_threshold": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Pool available space threshold percentage or `null` to disable."},
			"alua": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether Asymmetric Logical Unit Access (ALUA) is enabled. Enabling is limited to TrueNAS Enterprise-licensed     high availability systems. ALUA only works when configured on both the client and serve"},
			"iser": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether iSCSI Extensions for RDMA (iSER) are enabled. Enabling is limited to TrueNAS Enterprise-licensed     systems and requires the system and network environment have Remote Direct Memory Access (R"},
		},
	}
}

func (r *Iscsi_GlobalConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Iscsi_GlobalConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data Iscsi_GlobalConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Basename.IsNull() && !data.Basename.IsUnknown() { params["basename"] = data.Basename.ValueString() }
	if !data.ListenPort.IsNull() && !data.ListenPort.IsUnknown() { params["listen_port"] = data.ListenPort.ValueInt64() }
	if !data.PoolAvailThreshold.IsNull() && !data.PoolAvailThreshold.IsUnknown() { params["pool_avail_threshold"] = data.PoolAvailThreshold.ValueInt64() }
	if !data.Alua.IsNull() && !data.Alua.IsUnknown() { params["alua"] = data.Alua.ValueBool() }
	if !data.Iser.IsNull() && !data.Iser.IsUnknown() { params["iser"] = data.Iser.ValueBool() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("iscsi.global.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply iscsi_global config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("iscsi.global.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read iscsi_global config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["basename"]; ok { if v == nil { data.Basename = types.StringNull() } else if s, ok := v.(string); ok { data.Basename = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Basename = types.StringValue(string(j)) } } }
		if v, ok := m["listen_port"]; ok { if v == nil { data.ListenPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.ListenPort = types.Int64Value(int64(f)) } }
		if v, ok := m["pool_avail_threshold"]; ok { if v == nil { data.PoolAvailThreshold = types.Int64Null() } else if f, ok := v.(float64); ok { data.PoolAvailThreshold = types.Int64Value(int64(f)) } }
		if v, ok := m["alua"]; ok { if v == nil { data.Alua = types.BoolNull() } else if b, ok := v.(bool); ok { data.Alua = types.BoolValue(b) } }
		if v, ok := m["iser"]; ok { if v == nil { data.Iser = types.BoolNull() } else if b, ok := v.(bool); ok { data.Iser = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Iscsi_GlobalConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data Iscsi_GlobalConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("iscsi.global.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read iscsi_global config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["basename"]; ok { if v == nil { data.Basename = types.StringNull() } else if s, ok := v.(string); ok { data.Basename = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Basename = types.StringValue(string(j)) } } }
		if v, ok := m["listen_port"]; ok { if v == nil { data.ListenPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.ListenPort = types.Int64Value(int64(f)) } }
		if v, ok := m["pool_avail_threshold"]; ok { if v == nil { data.PoolAvailThreshold = types.Int64Null() } else if f, ok := v.(float64); ok { data.PoolAvailThreshold = types.Int64Value(int64(f)) } }
		if v, ok := m["alua"]; ok { if v == nil { data.Alua = types.BoolNull() } else if b, ok := v.(bool); ok { data.Alua = types.BoolValue(b) } }
		if v, ok := m["iser"]; ok { if v == nil { data.Iser = types.BoolNull() } else if b, ok := v.(bool); ok { data.Iser = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Iscsi_GlobalConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data Iscsi_GlobalConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Basename.IsNull() && !data.Basename.IsUnknown() { params["basename"] = data.Basename.ValueString() }
	if !data.ListenPort.IsNull() && !data.ListenPort.IsUnknown() { params["listen_port"] = data.ListenPort.ValueInt64() }
	if !data.PoolAvailThreshold.IsNull() && !data.PoolAvailThreshold.IsUnknown() { params["pool_avail_threshold"] = data.PoolAvailThreshold.ValueInt64() }
	if !data.Alua.IsNull() && !data.Alua.IsUnknown() { params["alua"] = data.Alua.ValueBool() }
	if !data.Iser.IsNull() && !data.Iser.IsUnknown() { params["iser"] = data.Iser.ValueBool() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("iscsi.global.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update iscsi_global config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("iscsi.global.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read iscsi_global config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["basename"]; ok { if v == nil { data.Basename = types.StringNull() } else if s, ok := v.(string); ok { data.Basename = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Basename = types.StringValue(string(j)) } } }
		if v, ok := m["listen_port"]; ok { if v == nil { data.ListenPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.ListenPort = types.Int64Value(int64(f)) } }
		if v, ok := m["pool_avail_threshold"]; ok { if v == nil { data.PoolAvailThreshold = types.Int64Null() } else if f, ok := v.(float64); ok { data.PoolAvailThreshold = types.Int64Value(int64(f)) } }
		if v, ok := m["alua"]; ok { if v == nil { data.Alua = types.BoolNull() } else if b, ok := v.(bool); ok { data.Alua = types.BoolValue(b) } }
		if v, ok := m["iser"]; ok { if v == nil { data.Iser = types.BoolNull() } else if b, ok := v.(bool); ok { data.Iser = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Iscsi_GlobalConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
