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

type NfsConfigResource struct {
	client *client.Client
}

type NfsConfigResourceModel struct {
	Servers types.Int64 `tfsdk:"servers"`
	AllowNonroot types.Bool `tfsdk:"allow_nonroot"`
	V4Krb types.Bool `tfsdk:"v4_krb"`
	V4Domain types.String `tfsdk:"v4_domain"`
	MountdPort types.Int64 `tfsdk:"mountd_port"`
	RpcstatdPort types.Int64 `tfsdk:"rpcstatd_port"`
	RpclockdPort types.Int64 `tfsdk:"rpclockd_port"`
	MountdLog types.Bool `tfsdk:"mountd_log"`
	StatdLockdLog types.Bool `tfsdk:"statd_lockd_log"`
	UserdManageGids types.Bool `tfsdk:"userd_manage_gids"`
	Rdma types.Bool `tfsdk:"rdma"`
}

func NewNfsConfigResource() resource.Resource {
	return &NfsConfigResource{}
}

func (r *NfsConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_config"
}

func (r *NfsConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NFS service configuration",
		Attributes: map[string]schema.Attribute{
			"servers": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Specify the number of nfsd. Default: Number of nfsd is equal number of CPU. "},
			"allow_nonroot": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Allow non-root mount requests.  This equates to 'insecure' share option. "},
			"v4_krb": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Force Kerberos authentication on NFS shares. "},
			"v4_domain": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Specify a DNS domain (NFSv4 only). "},
			"mountd_port": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Specify the mountd port binding. "},
			"rpcstatd_port": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Specify the rpc.statd port binding. "},
			"rpclockd_port": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Specify the rpc.lockd port binding. "},
			"mountd_log": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable or disable mountd logging. "},
			"statd_lockd_log": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable or disable statd and lockd logging. "},
			"userd_manage_gids": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable to allow server to manage gids. "},
			"rdma": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable or disable NFS over RDMA.  Requires RDMA capable NIC. "},
		},
	}
}

func (r *NfsConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NfsConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NfsConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Servers.IsNull() && !data.Servers.IsUnknown() { params["servers"] = data.Servers.ValueInt64() }
	if !data.AllowNonroot.IsNull() && !data.AllowNonroot.IsUnknown() { params["allow_nonroot"] = data.AllowNonroot.ValueBool() }
	if !data.V4Krb.IsNull() && !data.V4Krb.IsUnknown() { params["v4_krb"] = data.V4Krb.ValueBool() }
	if !data.V4Domain.IsNull() && !data.V4Domain.IsUnknown() { params["v4_domain"] = data.V4Domain.ValueString() }
	if !data.MountdPort.IsNull() && !data.MountdPort.IsUnknown() { params["mountd_port"] = data.MountdPort.ValueInt64() }
	if !data.RpcstatdPort.IsNull() && !data.RpcstatdPort.IsUnknown() { params["rpcstatd_port"] = data.RpcstatdPort.ValueInt64() }
	if !data.RpclockdPort.IsNull() && !data.RpclockdPort.IsUnknown() { params["rpclockd_port"] = data.RpclockdPort.ValueInt64() }
	if !data.MountdLog.IsNull() && !data.MountdLog.IsUnknown() { params["mountd_log"] = data.MountdLog.ValueBool() }
	if !data.StatdLockdLog.IsNull() && !data.StatdLockdLog.IsUnknown() { params["statd_lockd_log"] = data.StatdLockdLog.ValueBool() }
	if !data.UserdManageGids.IsNull() && !data.UserdManageGids.IsUnknown() { params["userd_manage_gids"] = data.UserdManageGids.ValueBool() }
	if !data.Rdma.IsNull() && !data.Rdma.IsUnknown() { params["rdma"] = data.Rdma.ValueBool() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("nfs.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply nfs config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("nfs.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read nfs config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["servers"]; ok { if v == nil { data.Servers = types.Int64Null() } else if f, ok := v.(float64); ok { data.Servers = types.Int64Value(int64(f)) } }
		if v, ok := m["allow_nonroot"]; ok { if v == nil { data.AllowNonroot = types.BoolNull() } else if b, ok := v.(bool); ok { data.AllowNonroot = types.BoolValue(b) } }
		if v, ok := m["v4_krb"]; ok { if v == nil { data.V4Krb = types.BoolNull() } else if b, ok := v.(bool); ok { data.V4Krb = types.BoolValue(b) } }
		if v, ok := m["v4_domain"]; ok { if v == nil { data.V4Domain = types.StringNull() } else if s, ok := v.(string); ok { data.V4Domain = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.V4Domain = types.StringValue(string(j)) } } }
		if v, ok := m["mountd_port"]; ok { if v == nil { data.MountdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.MountdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["rpcstatd_port"]; ok { if v == nil { data.RpcstatdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.RpcstatdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["rpclockd_port"]; ok { if v == nil { data.RpclockdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.RpclockdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["mountd_log"]; ok { if v == nil { data.MountdLog = types.BoolNull() } else if b, ok := v.(bool); ok { data.MountdLog = types.BoolValue(b) } }
		if v, ok := m["statd_lockd_log"]; ok { if v == nil { data.StatdLockdLog = types.BoolNull() } else if b, ok := v.(bool); ok { data.StatdLockdLog = types.BoolValue(b) } }
		if v, ok := m["userd_manage_gids"]; ok { if v == nil { data.UserdManageGids = types.BoolNull() } else if b, ok := v.(bool); ok { data.UserdManageGids = types.BoolValue(b) } }
		if v, ok := m["rdma"]; ok { if v == nil { data.Rdma = types.BoolNull() } else if b, ok := v.(bool); ok { data.Rdma = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NfsConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NfsConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("nfs.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read nfs config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["servers"]; ok { if v == nil { data.Servers = types.Int64Null() } else if f, ok := v.(float64); ok { data.Servers = types.Int64Value(int64(f)) } }
		if v, ok := m["allow_nonroot"]; ok { if v == nil { data.AllowNonroot = types.BoolNull() } else if b, ok := v.(bool); ok { data.AllowNonroot = types.BoolValue(b) } }
		if v, ok := m["v4_krb"]; ok { if v == nil { data.V4Krb = types.BoolNull() } else if b, ok := v.(bool); ok { data.V4Krb = types.BoolValue(b) } }
		if v, ok := m["v4_domain"]; ok { if v == nil { data.V4Domain = types.StringNull() } else if s, ok := v.(string); ok { data.V4Domain = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.V4Domain = types.StringValue(string(j)) } } }
		if v, ok := m["mountd_port"]; ok { if v == nil { data.MountdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.MountdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["rpcstatd_port"]; ok { if v == nil { data.RpcstatdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.RpcstatdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["rpclockd_port"]; ok { if v == nil { data.RpclockdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.RpclockdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["mountd_log"]; ok { if v == nil { data.MountdLog = types.BoolNull() } else if b, ok := v.(bool); ok { data.MountdLog = types.BoolValue(b) } }
		if v, ok := m["statd_lockd_log"]; ok { if v == nil { data.StatdLockdLog = types.BoolNull() } else if b, ok := v.(bool); ok { data.StatdLockdLog = types.BoolValue(b) } }
		if v, ok := m["userd_manage_gids"]; ok { if v == nil { data.UserdManageGids = types.BoolNull() } else if b, ok := v.(bool); ok { data.UserdManageGids = types.BoolValue(b) } }
		if v, ok := m["rdma"]; ok { if v == nil { data.Rdma = types.BoolNull() } else if b, ok := v.(bool); ok { data.Rdma = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NfsConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NfsConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Servers.IsNull() && !data.Servers.IsUnknown() { params["servers"] = data.Servers.ValueInt64() }
	if !data.AllowNonroot.IsNull() && !data.AllowNonroot.IsUnknown() { params["allow_nonroot"] = data.AllowNonroot.ValueBool() }
	if !data.V4Krb.IsNull() && !data.V4Krb.IsUnknown() { params["v4_krb"] = data.V4Krb.ValueBool() }
	if !data.V4Domain.IsNull() && !data.V4Domain.IsUnknown() { params["v4_domain"] = data.V4Domain.ValueString() }
	if !data.MountdPort.IsNull() && !data.MountdPort.IsUnknown() { params["mountd_port"] = data.MountdPort.ValueInt64() }
	if !data.RpcstatdPort.IsNull() && !data.RpcstatdPort.IsUnknown() { params["rpcstatd_port"] = data.RpcstatdPort.ValueInt64() }
	if !data.RpclockdPort.IsNull() && !data.RpclockdPort.IsUnknown() { params["rpclockd_port"] = data.RpclockdPort.ValueInt64() }
	if !data.MountdLog.IsNull() && !data.MountdLog.IsUnknown() { params["mountd_log"] = data.MountdLog.ValueBool() }
	if !data.StatdLockdLog.IsNull() && !data.StatdLockdLog.IsUnknown() { params["statd_lockd_log"] = data.StatdLockdLog.ValueBool() }
	if !data.UserdManageGids.IsNull() && !data.UserdManageGids.IsUnknown() { params["userd_manage_gids"] = data.UserdManageGids.ValueBool() }
	if !data.Rdma.IsNull() && !data.Rdma.IsUnknown() { params["rdma"] = data.Rdma.ValueBool() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("nfs.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update nfs config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("nfs.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read nfs config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["servers"]; ok { if v == nil { data.Servers = types.Int64Null() } else if f, ok := v.(float64); ok { data.Servers = types.Int64Value(int64(f)) } }
		if v, ok := m["allow_nonroot"]; ok { if v == nil { data.AllowNonroot = types.BoolNull() } else if b, ok := v.(bool); ok { data.AllowNonroot = types.BoolValue(b) } }
		if v, ok := m["v4_krb"]; ok { if v == nil { data.V4Krb = types.BoolNull() } else if b, ok := v.(bool); ok { data.V4Krb = types.BoolValue(b) } }
		if v, ok := m["v4_domain"]; ok { if v == nil { data.V4Domain = types.StringNull() } else if s, ok := v.(string); ok { data.V4Domain = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.V4Domain = types.StringValue(string(j)) } } }
		if v, ok := m["mountd_port"]; ok { if v == nil { data.MountdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.MountdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["rpcstatd_port"]; ok { if v == nil { data.RpcstatdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.RpcstatdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["rpclockd_port"]; ok { if v == nil { data.RpclockdPort = types.Int64Null() } else if f, ok := v.(float64); ok { data.RpclockdPort = types.Int64Value(int64(f)) } }
		if v, ok := m["mountd_log"]; ok { if v == nil { data.MountdLog = types.BoolNull() } else if b, ok := v.(bool); ok { data.MountdLog = types.BoolValue(b) } }
		if v, ok := m["statd_lockd_log"]; ok { if v == nil { data.StatdLockdLog = types.BoolNull() } else if b, ok := v.(bool); ok { data.StatdLockdLog = types.BoolValue(b) } }
		if v, ok := m["userd_manage_gids"]; ok { if v == nil { data.UserdManageGids = types.BoolNull() } else if b, ok := v.(bool); ok { data.UserdManageGids = types.BoolValue(b) } }
		if v, ok := m["rdma"]; ok { if v == nil { data.Rdma = types.BoolNull() } else if b, ok := v.(bool); ok { data.Rdma = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NfsConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
