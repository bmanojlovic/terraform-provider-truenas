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

type System_AdvancedConfigResource struct {
	client *client.Client
}

type System_AdvancedConfigResourceModel struct {
	Advancedmode       types.Bool   `tfsdk:"advancedmode"`
	Autotune           types.Bool   `tfsdk:"autotune"`
	KdumpEnabled       types.Bool   `tfsdk:"kdump_enabled"`
	BootScrub          types.Int64  `tfsdk:"boot_scrub"`
	Consolemenu        types.Bool   `tfsdk:"consolemenu"`
	Consolemsg         types.Bool   `tfsdk:"consolemsg"`
	Debugkernel        types.Bool   `tfsdk:"debugkernel"`
	FqdnSyslog         types.Bool   `tfsdk:"fqdn_syslog"`
	Motd               types.String `tfsdk:"motd"`
	LoginBanner        types.String `tfsdk:"login_banner"`
	Powerdaemon        types.Bool   `tfsdk:"powerdaemon"`
	Serialconsole      types.Bool   `tfsdk:"serialconsole"`
	Serialport         types.String `tfsdk:"serialport"`
	Serialspeed        types.String `tfsdk:"serialspeed"`
	Overprovision      types.Int64  `tfsdk:"overprovision"`
	Traceback          types.Bool   `tfsdk:"traceback"`
	Uploadcrash        types.Bool   `tfsdk:"uploadcrash"`
	Anonstats          types.Bool   `tfsdk:"anonstats"`
	SedUser            types.String `tfsdk:"sed_user"`
	Sysloglevel        types.String `tfsdk:"sysloglevel"`
	SyslogAudit        types.Bool   `tfsdk:"syslog_audit"`
	KernelExtraOptions types.String `tfsdk:"kernel_extra_options"`
	SedPasswd          types.String `tfsdk:"sed_passwd"`
}

func NewSystem_AdvancedConfigResource() resource.Resource {
	return &System_AdvancedConfigResource{}
}

func (r *System_AdvancedConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_advanced_config"
}

func (r *System_AdvancedConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "System advanced settings",
		Attributes: map[string]schema.Attribute{
			"advancedmode":         schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable advanced mode to show additional configuration options in the web interface."},
			"autotune":             schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Execute autotune script which attempts to optimize the system based on the installed hardware."},
			"kdump_enabled":        schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable kernel crash dumps for debugging system crashes."},
			"boot_scrub":           schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Number of days between automatic boot pool scrubs."},
			"consolemenu":          schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable console menu. Default to standard login in the console if disabled."},
			"consolemsg":           schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Deprecated: Please use `consolemsg` attribute in the `system.general` plugin instead."},
			"debugkernel":          schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable debug kernel for additional logging and debugging capabilities."},
			"fqdn_syslog":          schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Include the full domain name in syslog messages."},
			"motd":                 schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Message of the day displayed after login."},
			"login_banner":         schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Banner message displayed before login prompt."},
			"powerdaemon":          schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable the power management daemon for automatic power management."},
			"serialconsole":        schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable serial console access."},
			"serialport":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Serial port device for console access."},
			"serialspeed":          schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Baud rate for serial console communication."},
			"overprovision":        schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Percentage of SSD overprovisioning to reserve for wear leveling."},
			"traceback":            schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable generation and saving of tracebacks for debugging."},
			"uploadcrash":          schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Automatically upload crash reports to iXsystems for analysis."},
			"anonstats":            schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable anonymous usage statistics reporting to help improve TrueNAS."},
			"sed_user":             schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "SED (Self-Encrypting Drive) user type for drive encryption."},
			"sysloglevel":          schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Minimum log level for syslog messages. F_EMERG is most critical, F_DEBUG is least critical."},
			"syslog_audit":         schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "The remote syslog server(s) will also receive audit messages."},
			"kernel_extra_options": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional kernel boot parameters to pass to the Linux kernel."},
			"sed_passwd":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Password for SED (Self-Encrypting Drive) global unlock."},
		},
	}
}

func (r *System_AdvancedConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *System_AdvancedConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data System_AdvancedConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Advancedmode.IsNull() && !data.Advancedmode.IsUnknown() {
		params["advancedmode"] = data.Advancedmode.ValueBool()
	}
	if !data.Autotune.IsNull() && !data.Autotune.IsUnknown() {
		params["autotune"] = data.Autotune.ValueBool()
	}
	if !data.KdumpEnabled.IsNull() && !data.KdumpEnabled.IsUnknown() {
		params["kdump_enabled"] = data.KdumpEnabled.ValueBool()
	}
	if !data.BootScrub.IsNull() && !data.BootScrub.IsUnknown() {
		params["boot_scrub"] = data.BootScrub.ValueInt64()
	}
	if !data.Consolemenu.IsNull() && !data.Consolemenu.IsUnknown() {
		params["consolemenu"] = data.Consolemenu.ValueBool()
	}
	if !data.Consolemsg.IsNull() && !data.Consolemsg.IsUnknown() {
		params["consolemsg"] = data.Consolemsg.ValueBool()
	}
	if !data.Debugkernel.IsNull() && !data.Debugkernel.IsUnknown() {
		params["debugkernel"] = data.Debugkernel.ValueBool()
	}
	if !data.FqdnSyslog.IsNull() && !data.FqdnSyslog.IsUnknown() {
		params["fqdn_syslog"] = data.FqdnSyslog.ValueBool()
	}
	if !data.Motd.IsNull() && !data.Motd.IsUnknown() {
		params["motd"] = data.Motd.ValueString()
	}
	if !data.LoginBanner.IsNull() && !data.LoginBanner.IsUnknown() {
		params["login_banner"] = data.LoginBanner.ValueString()
	}
	if !data.Powerdaemon.IsNull() && !data.Powerdaemon.IsUnknown() {
		params["powerdaemon"] = data.Powerdaemon.ValueBool()
	}
	if !data.Serialconsole.IsNull() && !data.Serialconsole.IsUnknown() {
		params["serialconsole"] = data.Serialconsole.ValueBool()
	}
	if !data.Serialport.IsNull() && !data.Serialport.IsUnknown() {
		params["serialport"] = data.Serialport.ValueString()
	}
	if !data.Serialspeed.IsNull() && !data.Serialspeed.IsUnknown() {
		params["serialspeed"] = data.Serialspeed.ValueString()
	}
	if !data.Overprovision.IsNull() && !data.Overprovision.IsUnknown() {
		params["overprovision"] = data.Overprovision.ValueInt64()
	}
	if !data.Traceback.IsNull() && !data.Traceback.IsUnknown() {
		params["traceback"] = data.Traceback.ValueBool()
	}
	if !data.Uploadcrash.IsNull() && !data.Uploadcrash.IsUnknown() {
		params["uploadcrash"] = data.Uploadcrash.ValueBool()
	}
	if !data.Anonstats.IsNull() && !data.Anonstats.IsUnknown() {
		params["anonstats"] = data.Anonstats.ValueBool()
	}
	if !data.SedUser.IsNull() && !data.SedUser.IsUnknown() {
		params["sed_user"] = data.SedUser.ValueString()
	}
	if !data.Sysloglevel.IsNull() && !data.Sysloglevel.IsUnknown() {
		params["sysloglevel"] = data.Sysloglevel.ValueString()
	}
	if !data.SyslogAudit.IsNull() && !data.SyslogAudit.IsUnknown() {
		params["syslog_audit"] = data.SyslogAudit.ValueBool()
	}
	if !data.KernelExtraOptions.IsNull() && !data.KernelExtraOptions.IsUnknown() {
		params["kernel_extra_options"] = data.KernelExtraOptions.ValueString()
	}
	if !data.SedPasswd.IsNull() && !data.SedPasswd.IsUnknown() {
		params["sed_passwd"] = data.SedPasswd.ValueString()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("system.advanced.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply system_advanced config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("system.advanced.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read system_advanced config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["advancedmode"]; ok {
			if v == nil {
				data.Advancedmode = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Advancedmode = types.BoolValue(b)
			}
		}
		if v, ok := m["autotune"]; ok {
			if v == nil {
				data.Autotune = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Autotune = types.BoolValue(b)
			}
		}
		if v, ok := m["kdump_enabled"]; ok {
			if v == nil {
				data.KdumpEnabled = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.KdumpEnabled = types.BoolValue(b)
			}
		}
		if v, ok := m["boot_scrub"]; ok {
			if v == nil {
				data.BootScrub = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.BootScrub = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["consolemenu"]; ok {
			if v == nil {
				data.Consolemenu = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Consolemenu = types.BoolValue(b)
			}
		}
		if v, ok := m["consolemsg"]; ok {
			if v == nil {
				data.Consolemsg = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Consolemsg = types.BoolValue(b)
			}
		}
		if v, ok := m["debugkernel"]; ok {
			if v == nil {
				data.Debugkernel = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Debugkernel = types.BoolValue(b)
			}
		}
		if v, ok := m["fqdn_syslog"]; ok {
			if v == nil {
				data.FqdnSyslog = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.FqdnSyslog = types.BoolValue(b)
			}
		}
		if v, ok := m["motd"]; ok {
			if v == nil {
				data.Motd = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Motd = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Motd = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["login_banner"]; ok {
			if v == nil {
				data.LoginBanner = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.LoginBanner = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.LoginBanner = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["powerdaemon"]; ok {
			if v == nil {
				data.Powerdaemon = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Powerdaemon = types.BoolValue(b)
			}
		}
		if v, ok := m["serialconsole"]; ok {
			if v == nil {
				data.Serialconsole = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Serialconsole = types.BoolValue(b)
			}
		}
		if v, ok := m["serialport"]; ok {
			if v == nil {
				data.Serialport = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Serialport = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Serialport = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["serialspeed"]; ok {
			if v == nil {
				data.Serialspeed = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Serialspeed = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Serialspeed = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["overprovision"]; ok {
			if v == nil {
				data.Overprovision = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Overprovision = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["traceback"]; ok {
			if v == nil {
				data.Traceback = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Traceback = types.BoolValue(b)
			}
		}
		if v, ok := m["uploadcrash"]; ok {
			if v == nil {
				data.Uploadcrash = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Uploadcrash = types.BoolValue(b)
			}
		}
		if v, ok := m["anonstats"]; ok {
			if v == nil {
				data.Anonstats = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Anonstats = types.BoolValue(b)
			}
		}
		if v, ok := m["sed_user"]; ok {
			if v == nil {
				data.SedUser = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SedUser = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SedUser = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sysloglevel"]; ok {
			if v == nil {
				data.Sysloglevel = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Sysloglevel = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Sysloglevel = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["syslog_audit"]; ok {
			if v == nil {
				data.SyslogAudit = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.SyslogAudit = types.BoolValue(b)
			}
		}
		if v, ok := m["kernel_extra_options"]; ok {
			if v == nil {
				data.KernelExtraOptions = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.KernelExtraOptions = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.KernelExtraOptions = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sed_passwd"]; ok {
			if v == nil {
				data.SedPasswd = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SedPasswd = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SedPasswd = types.StringValue(string(j))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *System_AdvancedConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data System_AdvancedConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("system.advanced.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read system_advanced config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["advancedmode"]; ok {
			if v == nil {
				data.Advancedmode = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Advancedmode = types.BoolValue(b)
			}
		}
		if v, ok := m["autotune"]; ok {
			if v == nil {
				data.Autotune = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Autotune = types.BoolValue(b)
			}
		}
		if v, ok := m["kdump_enabled"]; ok {
			if v == nil {
				data.KdumpEnabled = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.KdumpEnabled = types.BoolValue(b)
			}
		}
		if v, ok := m["boot_scrub"]; ok {
			if v == nil {
				data.BootScrub = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.BootScrub = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["consolemenu"]; ok {
			if v == nil {
				data.Consolemenu = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Consolemenu = types.BoolValue(b)
			}
		}
		if v, ok := m["consolemsg"]; ok {
			if v == nil {
				data.Consolemsg = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Consolemsg = types.BoolValue(b)
			}
		}
		if v, ok := m["debugkernel"]; ok {
			if v == nil {
				data.Debugkernel = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Debugkernel = types.BoolValue(b)
			}
		}
		if v, ok := m["fqdn_syslog"]; ok {
			if v == nil {
				data.FqdnSyslog = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.FqdnSyslog = types.BoolValue(b)
			}
		}
		if v, ok := m["motd"]; ok {
			if v == nil {
				data.Motd = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Motd = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Motd = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["login_banner"]; ok {
			if v == nil {
				data.LoginBanner = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.LoginBanner = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.LoginBanner = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["powerdaemon"]; ok {
			if v == nil {
				data.Powerdaemon = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Powerdaemon = types.BoolValue(b)
			}
		}
		if v, ok := m["serialconsole"]; ok {
			if v == nil {
				data.Serialconsole = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Serialconsole = types.BoolValue(b)
			}
		}
		if v, ok := m["serialport"]; ok {
			if v == nil {
				data.Serialport = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Serialport = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Serialport = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["serialspeed"]; ok {
			if v == nil {
				data.Serialspeed = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Serialspeed = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Serialspeed = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["overprovision"]; ok {
			if v == nil {
				data.Overprovision = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Overprovision = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["traceback"]; ok {
			if v == nil {
				data.Traceback = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Traceback = types.BoolValue(b)
			}
		}
		if v, ok := m["uploadcrash"]; ok {
			if v == nil {
				data.Uploadcrash = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Uploadcrash = types.BoolValue(b)
			}
		}
		if v, ok := m["anonstats"]; ok {
			if v == nil {
				data.Anonstats = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Anonstats = types.BoolValue(b)
			}
		}
		if v, ok := m["sed_user"]; ok {
			if v == nil {
				data.SedUser = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SedUser = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SedUser = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sysloglevel"]; ok {
			if v == nil {
				data.Sysloglevel = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Sysloglevel = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Sysloglevel = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["syslog_audit"]; ok {
			if v == nil {
				data.SyslogAudit = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.SyslogAudit = types.BoolValue(b)
			}
		}
		if v, ok := m["kernel_extra_options"]; ok {
			if v == nil {
				data.KernelExtraOptions = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.KernelExtraOptions = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.KernelExtraOptions = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sed_passwd"]; ok {
			if v == nil {
				data.SedPasswd = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SedPasswd = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SedPasswd = types.StringValue(string(j))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *System_AdvancedConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data System_AdvancedConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Advancedmode.IsNull() && !data.Advancedmode.IsUnknown() {
		params["advancedmode"] = data.Advancedmode.ValueBool()
	}
	if !data.Autotune.IsNull() && !data.Autotune.IsUnknown() {
		params["autotune"] = data.Autotune.ValueBool()
	}
	if !data.KdumpEnabled.IsNull() && !data.KdumpEnabled.IsUnknown() {
		params["kdump_enabled"] = data.KdumpEnabled.ValueBool()
	}
	if !data.BootScrub.IsNull() && !data.BootScrub.IsUnknown() {
		params["boot_scrub"] = data.BootScrub.ValueInt64()
	}
	if !data.Consolemenu.IsNull() && !data.Consolemenu.IsUnknown() {
		params["consolemenu"] = data.Consolemenu.ValueBool()
	}
	if !data.Consolemsg.IsNull() && !data.Consolemsg.IsUnknown() {
		params["consolemsg"] = data.Consolemsg.ValueBool()
	}
	if !data.Debugkernel.IsNull() && !data.Debugkernel.IsUnknown() {
		params["debugkernel"] = data.Debugkernel.ValueBool()
	}
	if !data.FqdnSyslog.IsNull() && !data.FqdnSyslog.IsUnknown() {
		params["fqdn_syslog"] = data.FqdnSyslog.ValueBool()
	}
	if !data.Motd.IsNull() && !data.Motd.IsUnknown() {
		params["motd"] = data.Motd.ValueString()
	}
	if !data.LoginBanner.IsNull() && !data.LoginBanner.IsUnknown() {
		params["login_banner"] = data.LoginBanner.ValueString()
	}
	if !data.Powerdaemon.IsNull() && !data.Powerdaemon.IsUnknown() {
		params["powerdaemon"] = data.Powerdaemon.ValueBool()
	}
	if !data.Serialconsole.IsNull() && !data.Serialconsole.IsUnknown() {
		params["serialconsole"] = data.Serialconsole.ValueBool()
	}
	if !data.Serialport.IsNull() && !data.Serialport.IsUnknown() {
		params["serialport"] = data.Serialport.ValueString()
	}
	if !data.Serialspeed.IsNull() && !data.Serialspeed.IsUnknown() {
		params["serialspeed"] = data.Serialspeed.ValueString()
	}
	if !data.Overprovision.IsNull() && !data.Overprovision.IsUnknown() {
		params["overprovision"] = data.Overprovision.ValueInt64()
	}
	if !data.Traceback.IsNull() && !data.Traceback.IsUnknown() {
		params["traceback"] = data.Traceback.ValueBool()
	}
	if !data.Uploadcrash.IsNull() && !data.Uploadcrash.IsUnknown() {
		params["uploadcrash"] = data.Uploadcrash.ValueBool()
	}
	if !data.Anonstats.IsNull() && !data.Anonstats.IsUnknown() {
		params["anonstats"] = data.Anonstats.ValueBool()
	}
	if !data.SedUser.IsNull() && !data.SedUser.IsUnknown() {
		params["sed_user"] = data.SedUser.ValueString()
	}
	if !data.Sysloglevel.IsNull() && !data.Sysloglevel.IsUnknown() {
		params["sysloglevel"] = data.Sysloglevel.ValueString()
	}
	if !data.SyslogAudit.IsNull() && !data.SyslogAudit.IsUnknown() {
		params["syslog_audit"] = data.SyslogAudit.ValueBool()
	}
	if !data.KernelExtraOptions.IsNull() && !data.KernelExtraOptions.IsUnknown() {
		params["kernel_extra_options"] = data.KernelExtraOptions.ValueString()
	}
	if !data.SedPasswd.IsNull() && !data.SedPasswd.IsUnknown() {
		params["sed_passwd"] = data.SedPasswd.ValueString()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("system.advanced.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update system_advanced config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("system.advanced.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read system_advanced config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["advancedmode"]; ok {
			if v == nil {
				data.Advancedmode = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Advancedmode = types.BoolValue(b)
			}
		}
		if v, ok := m["autotune"]; ok {
			if v == nil {
				data.Autotune = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Autotune = types.BoolValue(b)
			}
		}
		if v, ok := m["kdump_enabled"]; ok {
			if v == nil {
				data.KdumpEnabled = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.KdumpEnabled = types.BoolValue(b)
			}
		}
		if v, ok := m["boot_scrub"]; ok {
			if v == nil {
				data.BootScrub = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.BootScrub = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["consolemenu"]; ok {
			if v == nil {
				data.Consolemenu = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Consolemenu = types.BoolValue(b)
			}
		}
		if v, ok := m["consolemsg"]; ok {
			if v == nil {
				data.Consolemsg = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Consolemsg = types.BoolValue(b)
			}
		}
		if v, ok := m["debugkernel"]; ok {
			if v == nil {
				data.Debugkernel = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Debugkernel = types.BoolValue(b)
			}
		}
		if v, ok := m["fqdn_syslog"]; ok {
			if v == nil {
				data.FqdnSyslog = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.FqdnSyslog = types.BoolValue(b)
			}
		}
		if v, ok := m["motd"]; ok {
			if v == nil {
				data.Motd = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Motd = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Motd = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["login_banner"]; ok {
			if v == nil {
				data.LoginBanner = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.LoginBanner = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.LoginBanner = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["powerdaemon"]; ok {
			if v == nil {
				data.Powerdaemon = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Powerdaemon = types.BoolValue(b)
			}
		}
		if v, ok := m["serialconsole"]; ok {
			if v == nil {
				data.Serialconsole = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Serialconsole = types.BoolValue(b)
			}
		}
		if v, ok := m["serialport"]; ok {
			if v == nil {
				data.Serialport = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Serialport = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Serialport = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["serialspeed"]; ok {
			if v == nil {
				data.Serialspeed = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Serialspeed = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Serialspeed = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["overprovision"]; ok {
			if v == nil {
				data.Overprovision = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Overprovision = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["traceback"]; ok {
			if v == nil {
				data.Traceback = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Traceback = types.BoolValue(b)
			}
		}
		if v, ok := m["uploadcrash"]; ok {
			if v == nil {
				data.Uploadcrash = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Uploadcrash = types.BoolValue(b)
			}
		}
		if v, ok := m["anonstats"]; ok {
			if v == nil {
				data.Anonstats = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Anonstats = types.BoolValue(b)
			}
		}
		if v, ok := m["sed_user"]; ok {
			if v == nil {
				data.SedUser = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SedUser = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SedUser = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sysloglevel"]; ok {
			if v == nil {
				data.Sysloglevel = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Sysloglevel = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Sysloglevel = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["syslog_audit"]; ok {
			if v == nil {
				data.SyslogAudit = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.SyslogAudit = types.BoolValue(b)
			}
		}
		if v, ok := m["kernel_extra_options"]; ok {
			if v == nil {
				data.KernelExtraOptions = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.KernelExtraOptions = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.KernelExtraOptions = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sed_passwd"]; ok {
			if v == nil {
				data.SedPasswd = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SedPasswd = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SedPasswd = types.StringValue(string(j))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *System_AdvancedConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
