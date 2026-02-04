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

type System_GeneralConfigResource struct {
	client *client.Client
}

type System_GeneralConfigResourceModel struct {
	UiCertificate   types.Int64  `tfsdk:"ui_certificate"`
	UiHttpsport     types.Int64  `tfsdk:"ui_httpsport"`
	UiHttpsredirect types.Bool   `tfsdk:"ui_httpsredirect"`
	UiPort          types.Int64  `tfsdk:"ui_port"`
	UiConsolemsg    types.Bool   `tfsdk:"ui_consolemsg"`
	UiXFrameOptions types.String `tfsdk:"ui_x_frame_options"`
	Kbdmap          types.String `tfsdk:"kbdmap"`
	Timezone        types.String `tfsdk:"timezone"`
	UsageCollection types.Bool   `tfsdk:"usage_collection"`
	DsAuth          types.Bool   `tfsdk:"ds_auth"`
	UiRestartDelay  types.Int64  `tfsdk:"ui_restart_delay"`
	RollbackTimeout types.Int64  `tfsdk:"rollback_timeout"`
}

func NewSystem_GeneralConfigResource() resource.Resource {
	return &System_GeneralConfigResource{}
}

func (r *System_GeneralConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_general_config"
}

func (r *System_GeneralConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "System general settings",
		Attributes: map[string]schema.Attribute{
			"ui_certificate":     schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Used to enable HTTPS access to the system. If `ui_certificate` is not configured on boot, it is automatically     created by the system."},
			"ui_httpsport":       schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "HTTPS port for the web UI."},
			"ui_httpsredirect":   schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "When set, makes sure that all HTTP requests are converted to HTTPS requests to better enhance security."},
			"ui_port":            schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "HTTP port for the web UI."},
			"ui_consolemsg":      schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to show console messages on the web UI."},
			"ui_x_frame_options": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "X-Frame-Options header policy for web UI security."},
			"kbdmap":             schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "System keyboard layout mapping."},
			"timezone":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "System timezone identifier."},
			"usage_collection":   schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether usage data collection is enabled. `null` if not set."},
			"ds_auth":            schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Controls whether configured Directory Service users that are granted with Privileges are allowed to log in to     the Web UI or use TrueNAS API."},
			"ui_restart_delay":   schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Delay in seconds before restarting the UI after configuration changes. `null` to use default."},
			"rollback_timeout":   schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Timeout in seconds for automatic rollback of UI changes. `null` for no timeout."},
		},
	}
}

func (r *System_GeneralConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *System_GeneralConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data System_GeneralConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.UiCertificate.IsNull() && !data.UiCertificate.IsUnknown() {
		params["ui_certificate"] = data.UiCertificate.ValueInt64()
	}
	if !data.UiHttpsport.IsNull() && !data.UiHttpsport.IsUnknown() {
		params["ui_httpsport"] = data.UiHttpsport.ValueInt64()
	}
	if !data.UiHttpsredirect.IsNull() && !data.UiHttpsredirect.IsUnknown() {
		params["ui_httpsredirect"] = data.UiHttpsredirect.ValueBool()
	}
	if !data.UiPort.IsNull() && !data.UiPort.IsUnknown() {
		params["ui_port"] = data.UiPort.ValueInt64()
	}
	if !data.UiConsolemsg.IsNull() && !data.UiConsolemsg.IsUnknown() {
		params["ui_consolemsg"] = data.UiConsolemsg.ValueBool()
	}
	if !data.UiXFrameOptions.IsNull() && !data.UiXFrameOptions.IsUnknown() {
		params["ui_x_frame_options"] = data.UiXFrameOptions.ValueString()
	}
	if !data.Kbdmap.IsNull() && !data.Kbdmap.IsUnknown() {
		params["kbdmap"] = data.Kbdmap.ValueString()
	}
	if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
		params["timezone"] = data.Timezone.ValueString()
	}
	if !data.UsageCollection.IsNull() && !data.UsageCollection.IsUnknown() {
		params["usage_collection"] = data.UsageCollection.ValueBool()
	}
	if !data.DsAuth.IsNull() && !data.DsAuth.IsUnknown() {
		params["ds_auth"] = data.DsAuth.ValueBool()
	}
	if !data.UiRestartDelay.IsNull() && !data.UiRestartDelay.IsUnknown() {
		params["ui_restart_delay"] = data.UiRestartDelay.ValueInt64()
	}
	if !data.RollbackTimeout.IsNull() && !data.RollbackTimeout.IsUnknown() {
		params["rollback_timeout"] = data.RollbackTimeout.ValueInt64()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("system.general.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply system_general config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("system.general.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read system_general config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["ui_certificate"]; ok {
			if v == nil {
				data.UiCertificate = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiCertificate = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_httpsport"]; ok {
			if v == nil {
				data.UiHttpsport = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiHttpsport = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_httpsredirect"]; ok {
			if v == nil {
				data.UiHttpsredirect = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UiHttpsredirect = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_port"]; ok {
			if v == nil {
				data.UiPort = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiPort = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_consolemsg"]; ok {
			if v == nil {
				data.UiConsolemsg = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UiConsolemsg = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_x_frame_options"]; ok {
			if v == nil {
				data.UiXFrameOptions = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.UiXFrameOptions = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.UiXFrameOptions = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["kbdmap"]; ok {
			if v == nil {
				data.Kbdmap = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Kbdmap = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Kbdmap = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["timezone"]; ok {
			if v == nil {
				data.Timezone = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Timezone = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Timezone = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["usage_collection"]; ok {
			if v == nil {
				data.UsageCollection = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UsageCollection = types.BoolValue(b)
			}
		}
		if v, ok := m["ds_auth"]; ok {
			if v == nil {
				data.DsAuth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.DsAuth = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_restart_delay"]; ok {
			if v == nil {
				data.UiRestartDelay = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiRestartDelay = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["rollback_timeout"]; ok {
			if v == nil {
				data.RollbackTimeout = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.RollbackTimeout = types.Int64Value(int64(f))
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *System_GeneralConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data System_GeneralConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("system.general.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read system_general config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["ui_certificate"]; ok {
			if v == nil {
				data.UiCertificate = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiCertificate = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_httpsport"]; ok {
			if v == nil {
				data.UiHttpsport = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiHttpsport = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_httpsredirect"]; ok {
			if v == nil {
				data.UiHttpsredirect = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UiHttpsredirect = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_port"]; ok {
			if v == nil {
				data.UiPort = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiPort = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_consolemsg"]; ok {
			if v == nil {
				data.UiConsolemsg = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UiConsolemsg = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_x_frame_options"]; ok {
			if v == nil {
				data.UiXFrameOptions = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.UiXFrameOptions = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.UiXFrameOptions = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["kbdmap"]; ok {
			if v == nil {
				data.Kbdmap = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Kbdmap = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Kbdmap = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["timezone"]; ok {
			if v == nil {
				data.Timezone = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Timezone = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Timezone = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["usage_collection"]; ok {
			if v == nil {
				data.UsageCollection = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UsageCollection = types.BoolValue(b)
			}
		}
		if v, ok := m["ds_auth"]; ok {
			if v == nil {
				data.DsAuth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.DsAuth = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_restart_delay"]; ok {
			if v == nil {
				data.UiRestartDelay = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiRestartDelay = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["rollback_timeout"]; ok {
			if v == nil {
				data.RollbackTimeout = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.RollbackTimeout = types.Int64Value(int64(f))
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *System_GeneralConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data System_GeneralConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.UiCertificate.IsNull() && !data.UiCertificate.IsUnknown() {
		params["ui_certificate"] = data.UiCertificate.ValueInt64()
	}
	if !data.UiHttpsport.IsNull() && !data.UiHttpsport.IsUnknown() {
		params["ui_httpsport"] = data.UiHttpsport.ValueInt64()
	}
	if !data.UiHttpsredirect.IsNull() && !data.UiHttpsredirect.IsUnknown() {
		params["ui_httpsredirect"] = data.UiHttpsredirect.ValueBool()
	}
	if !data.UiPort.IsNull() && !data.UiPort.IsUnknown() {
		params["ui_port"] = data.UiPort.ValueInt64()
	}
	if !data.UiConsolemsg.IsNull() && !data.UiConsolemsg.IsUnknown() {
		params["ui_consolemsg"] = data.UiConsolemsg.ValueBool()
	}
	if !data.UiXFrameOptions.IsNull() && !data.UiXFrameOptions.IsUnknown() {
		params["ui_x_frame_options"] = data.UiXFrameOptions.ValueString()
	}
	if !data.Kbdmap.IsNull() && !data.Kbdmap.IsUnknown() {
		params["kbdmap"] = data.Kbdmap.ValueString()
	}
	if !data.Timezone.IsNull() && !data.Timezone.IsUnknown() {
		params["timezone"] = data.Timezone.ValueString()
	}
	if !data.UsageCollection.IsNull() && !data.UsageCollection.IsUnknown() {
		params["usage_collection"] = data.UsageCollection.ValueBool()
	}
	if !data.DsAuth.IsNull() && !data.DsAuth.IsUnknown() {
		params["ds_auth"] = data.DsAuth.ValueBool()
	}
	if !data.UiRestartDelay.IsNull() && !data.UiRestartDelay.IsUnknown() {
		params["ui_restart_delay"] = data.UiRestartDelay.ValueInt64()
	}
	if !data.RollbackTimeout.IsNull() && !data.RollbackTimeout.IsUnknown() {
		params["rollback_timeout"] = data.RollbackTimeout.ValueInt64()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("system.general.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update system_general config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("system.general.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read system_general config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["ui_certificate"]; ok {
			if v == nil {
				data.UiCertificate = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiCertificate = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_httpsport"]; ok {
			if v == nil {
				data.UiHttpsport = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiHttpsport = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_httpsredirect"]; ok {
			if v == nil {
				data.UiHttpsredirect = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UiHttpsredirect = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_port"]; ok {
			if v == nil {
				data.UiPort = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiPort = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["ui_consolemsg"]; ok {
			if v == nil {
				data.UiConsolemsg = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UiConsolemsg = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_x_frame_options"]; ok {
			if v == nil {
				data.UiXFrameOptions = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.UiXFrameOptions = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.UiXFrameOptions = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["kbdmap"]; ok {
			if v == nil {
				data.Kbdmap = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Kbdmap = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Kbdmap = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["timezone"]; ok {
			if v == nil {
				data.Timezone = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Timezone = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Timezone = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["usage_collection"]; ok {
			if v == nil {
				data.UsageCollection = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.UsageCollection = types.BoolValue(b)
			}
		}
		if v, ok := m["ds_auth"]; ok {
			if v == nil {
				data.DsAuth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.DsAuth = types.BoolValue(b)
			}
		}
		if v, ok := m["ui_restart_delay"]; ok {
			if v == nil {
				data.UiRestartDelay = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.UiRestartDelay = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["rollback_timeout"]; ok {
			if v == nil {
				data.RollbackTimeout = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.RollbackTimeout = types.Int64Value(int64(f))
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *System_GeneralConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
