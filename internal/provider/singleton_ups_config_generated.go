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

type UpsConfigResource struct {
	client *client.Client
}

type UpsConfigResourceModel struct {
	Powerdown types.Bool `tfsdk:"powerdown"`
	Rmonitor types.Bool `tfsdk:"rmonitor"`
	Nocommwarntime types.Int64 `tfsdk:"nocommwarntime"`
	Remoteport types.Int64 `tfsdk:"remoteport"`
	Shutdowntimer types.Int64 `tfsdk:"shutdowntimer"`
	Hostsync types.Int64 `tfsdk:"hostsync"`
	Description types.String `tfsdk:"description"`
	Driver types.String `tfsdk:"driver"`
	Extrausers types.String `tfsdk:"extrausers"`
	Identifier types.String `tfsdk:"identifier"`
	Mode types.String `tfsdk:"mode"`
	Monpwd types.String `tfsdk:"monpwd"`
	Monuser types.String `tfsdk:"monuser"`
	Options types.String `tfsdk:"options"`
	Optionsupsd types.String `tfsdk:"optionsupsd"`
	Port types.String `tfsdk:"port"`
	Remotehost types.String `tfsdk:"remotehost"`
	Shutdown types.String `tfsdk:"shutdown"`
	Shutdowncmd types.String `tfsdk:"shutdowncmd"`
}

func NewUpsConfigResource() resource.Resource {
	return &UpsConfigResource{}
}

func (r *UpsConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ups_config"
}

func (r *UpsConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "UPS service configuration",
		Attributes: map[string]schema.Attribute{
			"powerdown": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether the UPS should power down after completing the shutdown sequence."},
			"rmonitor": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to enable remote monitoring of the UPS status over the network."},
			"nocommwarntime": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Seconds to wait before warning about communication loss with UPS. `null` for default."},
			"remoteport": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Network port for communicating with remote UPS monitoring systems."},
			"shutdowntimer": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Seconds to wait after initiating shutdown before forcing power off."},
			"hostsync": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum seconds to wait for other systems to shutdown before continuing."},
			"description": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Human-readable description of this UPS configuration."},
			"driver": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "UPS driver name that handles communication with the specific UPS hardware model."},
			"extrausers": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional user configurations for UPS monitoring access."},
			"identifier": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Unique identifier name for this UPS device within the monitoring system."},
			"mode": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Operating mode. * `MASTER` controls the UPS directly * `SLAVE` monitors remotely"},
			"monpwd": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Password for UPS monitoring authentication (required for updates)."},
			"monuser": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Username for UPS monitoring authentication."},
			"options": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional configuration options passed to the UPS driver."},
			"optionsupsd": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional configuration options for the UPS daemon."},
			"port": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Serial port or device path for UPS communication."},
			"remotehost": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Hostname or IP address of remote UPS server when operating in SLAVE mode."},
			"shutdown": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Shutdown trigger condition: LOWBATT on low battery, BATT when on battery power."},
			"shutdowncmd": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Custom command to execute during UPS shutdown sequence. `null` for default."},
		},
	}
}

func (r *UpsConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UpsConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UpsConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Powerdown.IsNull() && !data.Powerdown.IsUnknown() { params["powerdown"] = data.Powerdown.ValueBool() }
	if !data.Rmonitor.IsNull() && !data.Rmonitor.IsUnknown() { params["rmonitor"] = data.Rmonitor.ValueBool() }
	if !data.Nocommwarntime.IsNull() && !data.Nocommwarntime.IsUnknown() { params["nocommwarntime"] = data.Nocommwarntime.ValueInt64() }
	if !data.Remoteport.IsNull() && !data.Remoteport.IsUnknown() { params["remoteport"] = data.Remoteport.ValueInt64() }
	if !data.Shutdowntimer.IsNull() && !data.Shutdowntimer.IsUnknown() { params["shutdowntimer"] = data.Shutdowntimer.ValueInt64() }
	if !data.Hostsync.IsNull() && !data.Hostsync.IsUnknown() { params["hostsync"] = data.Hostsync.ValueInt64() }
	if !data.Description.IsNull() && !data.Description.IsUnknown() { params["description"] = data.Description.ValueString() }
	if !data.Driver.IsNull() && !data.Driver.IsUnknown() { params["driver"] = data.Driver.ValueString() }
	if !data.Extrausers.IsNull() && !data.Extrausers.IsUnknown() { params["extrausers"] = data.Extrausers.ValueString() }
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() { params["identifier"] = data.Identifier.ValueString() }
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() { params["mode"] = data.Mode.ValueString() }
	if !data.Monpwd.IsNull() && !data.Monpwd.IsUnknown() { params["monpwd"] = data.Monpwd.ValueString() }
	if !data.Monuser.IsNull() && !data.Monuser.IsUnknown() { params["monuser"] = data.Monuser.ValueString() }
	if !data.Options.IsNull() && !data.Options.IsUnknown() { params["options"] = data.Options.ValueString() }
	if !data.Optionsupsd.IsNull() && !data.Optionsupsd.IsUnknown() { params["optionsupsd"] = data.Optionsupsd.ValueString() }
	if !data.Port.IsNull() && !data.Port.IsUnknown() { params["port"] = data.Port.ValueString() }
	if !data.Remotehost.IsNull() && !data.Remotehost.IsUnknown() { params["remotehost"] = data.Remotehost.ValueString() }
	if !data.Shutdown.IsNull() && !data.Shutdown.IsUnknown() { params["shutdown"] = data.Shutdown.ValueString() }
	if !data.Shutdowncmd.IsNull() && !data.Shutdowncmd.IsUnknown() { params["shutdowncmd"] = data.Shutdowncmd.ValueString() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("ups.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply ups config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("ups.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ups config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["powerdown"]; ok { if v == nil { data.Powerdown = types.BoolNull() } else if b, ok := v.(bool); ok { data.Powerdown = types.BoolValue(b) } }
		if v, ok := m["rmonitor"]; ok { if v == nil { data.Rmonitor = types.BoolNull() } else if b, ok := v.(bool); ok { data.Rmonitor = types.BoolValue(b) } }
		if v, ok := m["nocommwarntime"]; ok { if v == nil { data.Nocommwarntime = types.Int64Null() } else if f, ok := v.(float64); ok { data.Nocommwarntime = types.Int64Value(int64(f)) } }
		if v, ok := m["remoteport"]; ok { if v == nil { data.Remoteport = types.Int64Null() } else if f, ok := v.(float64); ok { data.Remoteport = types.Int64Value(int64(f)) } }
		if v, ok := m["shutdowntimer"]; ok { if v == nil { data.Shutdowntimer = types.Int64Null() } else if f, ok := v.(float64); ok { data.Shutdowntimer = types.Int64Value(int64(f)) } }
		if v, ok := m["hostsync"]; ok { if v == nil { data.Hostsync = types.Int64Null() } else if f, ok := v.(float64); ok { data.Hostsync = types.Int64Value(int64(f)) } }
		if v, ok := m["description"]; ok { if v == nil { data.Description = types.StringNull() } else if s, ok := v.(string); ok { data.Description = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Description = types.StringValue(string(j)) } } }
		if v, ok := m["driver"]; ok { if v == nil { data.Driver = types.StringNull() } else if s, ok := v.(string); ok { data.Driver = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Driver = types.StringValue(string(j)) } } }
		if v, ok := m["extrausers"]; ok { if v == nil { data.Extrausers = types.StringNull() } else if s, ok := v.(string); ok { data.Extrausers = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Extrausers = types.StringValue(string(j)) } } }
		if v, ok := m["identifier"]; ok { if v == nil { data.Identifier = types.StringNull() } else if s, ok := v.(string); ok { data.Identifier = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Identifier = types.StringValue(string(j)) } } }
		if v, ok := m["mode"]; ok { if v == nil { data.Mode = types.StringNull() } else if s, ok := v.(string); ok { data.Mode = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Mode = types.StringValue(string(j)) } } }
		if v, ok := m["monpwd"]; ok { if v == nil { data.Monpwd = types.StringNull() } else if s, ok := v.(string); ok { data.Monpwd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Monpwd = types.StringValue(string(j)) } } }
		if v, ok := m["monuser"]; ok { if v == nil { data.Monuser = types.StringNull() } else if s, ok := v.(string); ok { data.Monuser = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Monuser = types.StringValue(string(j)) } } }
		if v, ok := m["options"]; ok { if v == nil { data.Options = types.StringNull() } else if s, ok := v.(string); ok { data.Options = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Options = types.StringValue(string(j)) } } }
		if v, ok := m["optionsupsd"]; ok { if v == nil { data.Optionsupsd = types.StringNull() } else if s, ok := v.(string); ok { data.Optionsupsd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Optionsupsd = types.StringValue(string(j)) } } }
		if v, ok := m["port"]; ok { if v == nil { data.Port = types.StringNull() } else if s, ok := v.(string); ok { data.Port = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Port = types.StringValue(string(j)) } } }
		if v, ok := m["remotehost"]; ok { if v == nil { data.Remotehost = types.StringNull() } else if s, ok := v.(string); ok { data.Remotehost = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Remotehost = types.StringValue(string(j)) } } }
		if v, ok := m["shutdown"]; ok { if v == nil { data.Shutdown = types.StringNull() } else if s, ok := v.(string); ok { data.Shutdown = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Shutdown = types.StringValue(string(j)) } } }
		if v, ok := m["shutdowncmd"]; ok { if v == nil { data.Shutdowncmd = types.StringNull() } else if s, ok := v.(string); ok { data.Shutdowncmd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Shutdowncmd = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UpsConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UpsConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("ups.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ups config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["powerdown"]; ok { if v == nil { data.Powerdown = types.BoolNull() } else if b, ok := v.(bool); ok { data.Powerdown = types.BoolValue(b) } }
		if v, ok := m["rmonitor"]; ok { if v == nil { data.Rmonitor = types.BoolNull() } else if b, ok := v.(bool); ok { data.Rmonitor = types.BoolValue(b) } }
		if v, ok := m["nocommwarntime"]; ok { if v == nil { data.Nocommwarntime = types.Int64Null() } else if f, ok := v.(float64); ok { data.Nocommwarntime = types.Int64Value(int64(f)) } }
		if v, ok := m["remoteport"]; ok { if v == nil { data.Remoteport = types.Int64Null() } else if f, ok := v.(float64); ok { data.Remoteport = types.Int64Value(int64(f)) } }
		if v, ok := m["shutdowntimer"]; ok { if v == nil { data.Shutdowntimer = types.Int64Null() } else if f, ok := v.(float64); ok { data.Shutdowntimer = types.Int64Value(int64(f)) } }
		if v, ok := m["hostsync"]; ok { if v == nil { data.Hostsync = types.Int64Null() } else if f, ok := v.(float64); ok { data.Hostsync = types.Int64Value(int64(f)) } }
		if v, ok := m["description"]; ok { if v == nil { data.Description = types.StringNull() } else if s, ok := v.(string); ok { data.Description = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Description = types.StringValue(string(j)) } } }
		if v, ok := m["driver"]; ok { if v == nil { data.Driver = types.StringNull() } else if s, ok := v.(string); ok { data.Driver = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Driver = types.StringValue(string(j)) } } }
		if v, ok := m["extrausers"]; ok { if v == nil { data.Extrausers = types.StringNull() } else if s, ok := v.(string); ok { data.Extrausers = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Extrausers = types.StringValue(string(j)) } } }
		if v, ok := m["identifier"]; ok { if v == nil { data.Identifier = types.StringNull() } else if s, ok := v.(string); ok { data.Identifier = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Identifier = types.StringValue(string(j)) } } }
		if v, ok := m["mode"]; ok { if v == nil { data.Mode = types.StringNull() } else if s, ok := v.(string); ok { data.Mode = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Mode = types.StringValue(string(j)) } } }
		if v, ok := m["monpwd"]; ok { if v == nil { data.Monpwd = types.StringNull() } else if s, ok := v.(string); ok { data.Monpwd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Monpwd = types.StringValue(string(j)) } } }
		if v, ok := m["monuser"]; ok { if v == nil { data.Monuser = types.StringNull() } else if s, ok := v.(string); ok { data.Monuser = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Monuser = types.StringValue(string(j)) } } }
		if v, ok := m["options"]; ok { if v == nil { data.Options = types.StringNull() } else if s, ok := v.(string); ok { data.Options = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Options = types.StringValue(string(j)) } } }
		if v, ok := m["optionsupsd"]; ok { if v == nil { data.Optionsupsd = types.StringNull() } else if s, ok := v.(string); ok { data.Optionsupsd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Optionsupsd = types.StringValue(string(j)) } } }
		if v, ok := m["port"]; ok { if v == nil { data.Port = types.StringNull() } else if s, ok := v.(string); ok { data.Port = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Port = types.StringValue(string(j)) } } }
		if v, ok := m["remotehost"]; ok { if v == nil { data.Remotehost = types.StringNull() } else if s, ok := v.(string); ok { data.Remotehost = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Remotehost = types.StringValue(string(j)) } } }
		if v, ok := m["shutdown"]; ok { if v == nil { data.Shutdown = types.StringNull() } else if s, ok := v.(string); ok { data.Shutdown = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Shutdown = types.StringValue(string(j)) } } }
		if v, ok := m["shutdowncmd"]; ok { if v == nil { data.Shutdowncmd = types.StringNull() } else if s, ok := v.(string); ok { data.Shutdowncmd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Shutdowncmd = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UpsConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UpsConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Powerdown.IsNull() && !data.Powerdown.IsUnknown() { params["powerdown"] = data.Powerdown.ValueBool() }
	if !data.Rmonitor.IsNull() && !data.Rmonitor.IsUnknown() { params["rmonitor"] = data.Rmonitor.ValueBool() }
	if !data.Nocommwarntime.IsNull() && !data.Nocommwarntime.IsUnknown() { params["nocommwarntime"] = data.Nocommwarntime.ValueInt64() }
	if !data.Remoteport.IsNull() && !data.Remoteport.IsUnknown() { params["remoteport"] = data.Remoteport.ValueInt64() }
	if !data.Shutdowntimer.IsNull() && !data.Shutdowntimer.IsUnknown() { params["shutdowntimer"] = data.Shutdowntimer.ValueInt64() }
	if !data.Hostsync.IsNull() && !data.Hostsync.IsUnknown() { params["hostsync"] = data.Hostsync.ValueInt64() }
	if !data.Description.IsNull() && !data.Description.IsUnknown() { params["description"] = data.Description.ValueString() }
	if !data.Driver.IsNull() && !data.Driver.IsUnknown() { params["driver"] = data.Driver.ValueString() }
	if !data.Extrausers.IsNull() && !data.Extrausers.IsUnknown() { params["extrausers"] = data.Extrausers.ValueString() }
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() { params["identifier"] = data.Identifier.ValueString() }
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() { params["mode"] = data.Mode.ValueString() }
	if !data.Monpwd.IsNull() && !data.Monpwd.IsUnknown() { params["monpwd"] = data.Monpwd.ValueString() }
	if !data.Monuser.IsNull() && !data.Monuser.IsUnknown() { params["monuser"] = data.Monuser.ValueString() }
	if !data.Options.IsNull() && !data.Options.IsUnknown() { params["options"] = data.Options.ValueString() }
	if !data.Optionsupsd.IsNull() && !data.Optionsupsd.IsUnknown() { params["optionsupsd"] = data.Optionsupsd.ValueString() }
	if !data.Port.IsNull() && !data.Port.IsUnknown() { params["port"] = data.Port.ValueString() }
	if !data.Remotehost.IsNull() && !data.Remotehost.IsUnknown() { params["remotehost"] = data.Remotehost.ValueString() }
	if !data.Shutdown.IsNull() && !data.Shutdown.IsUnknown() { params["shutdown"] = data.Shutdown.ValueString() }
	if !data.Shutdowncmd.IsNull() && !data.Shutdowncmd.IsUnknown() { params["shutdowncmd"] = data.Shutdowncmd.ValueString() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("ups.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update ups config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("ups.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ups config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["powerdown"]; ok { if v == nil { data.Powerdown = types.BoolNull() } else if b, ok := v.(bool); ok { data.Powerdown = types.BoolValue(b) } }
		if v, ok := m["rmonitor"]; ok { if v == nil { data.Rmonitor = types.BoolNull() } else if b, ok := v.(bool); ok { data.Rmonitor = types.BoolValue(b) } }
		if v, ok := m["nocommwarntime"]; ok { if v == nil { data.Nocommwarntime = types.Int64Null() } else if f, ok := v.(float64); ok { data.Nocommwarntime = types.Int64Value(int64(f)) } }
		if v, ok := m["remoteport"]; ok { if v == nil { data.Remoteport = types.Int64Null() } else if f, ok := v.(float64); ok { data.Remoteport = types.Int64Value(int64(f)) } }
		if v, ok := m["shutdowntimer"]; ok { if v == nil { data.Shutdowntimer = types.Int64Null() } else if f, ok := v.(float64); ok { data.Shutdowntimer = types.Int64Value(int64(f)) } }
		if v, ok := m["hostsync"]; ok { if v == nil { data.Hostsync = types.Int64Null() } else if f, ok := v.(float64); ok { data.Hostsync = types.Int64Value(int64(f)) } }
		if v, ok := m["description"]; ok { if v == nil { data.Description = types.StringNull() } else if s, ok := v.(string); ok { data.Description = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Description = types.StringValue(string(j)) } } }
		if v, ok := m["driver"]; ok { if v == nil { data.Driver = types.StringNull() } else if s, ok := v.(string); ok { data.Driver = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Driver = types.StringValue(string(j)) } } }
		if v, ok := m["extrausers"]; ok { if v == nil { data.Extrausers = types.StringNull() } else if s, ok := v.(string); ok { data.Extrausers = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Extrausers = types.StringValue(string(j)) } } }
		if v, ok := m["identifier"]; ok { if v == nil { data.Identifier = types.StringNull() } else if s, ok := v.(string); ok { data.Identifier = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Identifier = types.StringValue(string(j)) } } }
		if v, ok := m["mode"]; ok { if v == nil { data.Mode = types.StringNull() } else if s, ok := v.(string); ok { data.Mode = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Mode = types.StringValue(string(j)) } } }
		if v, ok := m["monpwd"]; ok { if v == nil { data.Monpwd = types.StringNull() } else if s, ok := v.(string); ok { data.Monpwd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Monpwd = types.StringValue(string(j)) } } }
		if v, ok := m["monuser"]; ok { if v == nil { data.Monuser = types.StringNull() } else if s, ok := v.(string); ok { data.Monuser = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Monuser = types.StringValue(string(j)) } } }
		if v, ok := m["options"]; ok { if v == nil { data.Options = types.StringNull() } else if s, ok := v.(string); ok { data.Options = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Options = types.StringValue(string(j)) } } }
		if v, ok := m["optionsupsd"]; ok { if v == nil { data.Optionsupsd = types.StringNull() } else if s, ok := v.(string); ok { data.Optionsupsd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Optionsupsd = types.StringValue(string(j)) } } }
		if v, ok := m["port"]; ok { if v == nil { data.Port = types.StringNull() } else if s, ok := v.(string); ok { data.Port = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Port = types.StringValue(string(j)) } } }
		if v, ok := m["remotehost"]; ok { if v == nil { data.Remotehost = types.StringNull() } else if s, ok := v.(string); ok { data.Remotehost = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Remotehost = types.StringValue(string(j)) } } }
		if v, ok := m["shutdown"]; ok { if v == nil { data.Shutdown = types.StringNull() } else if s, ok := v.(string); ok { data.Shutdown = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Shutdown = types.StringValue(string(j)) } } }
		if v, ok := m["shutdowncmd"]; ok { if v == nil { data.Shutdowncmd = types.StringNull() } else if s, ok := v.(string); ok { data.Shutdowncmd = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Shutdowncmd = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UpsConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
