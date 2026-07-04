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

type SnmpConfigResource struct {
	client *client.Client
}

type SnmpConfigResourceModel struct {
	Location         types.String `tfsdk:"location"`
	Contact          types.String `tfsdk:"contact"`
	Traps            types.Bool   `tfsdk:"traps"`
	V3               types.Bool   `tfsdk:"v3"`
	Community        types.String `tfsdk:"community"`
	V3Username       types.String `tfsdk:"v3_username"`
	V3Authtype       types.String `tfsdk:"v3_authtype"`
	V3Password       types.String `tfsdk:"v3_password"`
	V3Privproto      types.String `tfsdk:"v3_privproto"`
	V3Privpassphrase types.String `tfsdk:"v3_privpassphrase"`
	Loglevel         types.Int64  `tfsdk:"loglevel"`
	Options          types.String `tfsdk:"options"`
	Zilstat          types.Bool   `tfsdk:"zilstat"`
}

func NewSnmpConfigResource() resource.Resource {
	return &SnmpConfigResource{}
}

func (r *SnmpConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmp_config"
}

func (r *SnmpConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SNMP service configuration",
		Attributes: map[string]schema.Attribute{
			"location":          schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "A comment describing the physical location of the server."},
			"contact":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Contact information for the system administrator (email or name)."},
			"traps":             schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether SNMP traps are enabled."},
			"v3":                schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether SNMP version 3 is enabled.  Enabling version 3 also requires username, authtype and password."},
			"community":         schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "SNMP community string for v1/v2c access. Allows         letters and numbers: a-zA-Z0-9          special characters: !$%&()+-_={}[]<>,.?          and spaces. Notable excluded characters: # / \\ @"},
			"v3_username":       schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Username for SNMP version 3 authentication."},
			"v3_authtype":       schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Authentication type for SNMP version 3 (empty string means no authentication)."},
			"v3_password":       schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Password for SNMP version 3 authentication."},
			"v3_privproto":      schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Privacy protocol for SNMP version 3 encryption. `null` means no encryption.      If set, ['AES'|'DES'], a `privpassphrase` must be supplied."},
			"v3_privpassphrase": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Privacy passphrase for SNMP version 3 encryption. This field is required when `privproto` is set."},
			"loglevel":          schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Logging level for SNMP daemon (0=emergency to 7=debug)."},
			"options":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional SNMP daemon configuration options.     Manual settings should be used with caution as they may render the SNMP service non-functional."},
			"zilstat":           schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to enable ZFS dataset statistics collection for SNMP."},
		},
	}
}

func (r *SnmpConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SnmpConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SnmpConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		params["location"] = data.Location.ValueString()
	}
	if !data.Contact.IsNull() && !data.Contact.IsUnknown() {
		params["contact"] = data.Contact.ValueString()
	}
	if !data.Traps.IsNull() && !data.Traps.IsUnknown() {
		params["traps"] = data.Traps.ValueBool()
	}
	if !data.V3.IsNull() && !data.V3.IsUnknown() {
		params["v3"] = data.V3.ValueBool()
	}
	if !data.Community.IsNull() && !data.Community.IsUnknown() {
		params["community"] = data.Community.ValueString()
	}
	if !data.V3Username.IsNull() && !data.V3Username.IsUnknown() {
		params["v3_username"] = data.V3Username.ValueString()
	}
	if !data.V3Authtype.IsNull() && !data.V3Authtype.IsUnknown() {
		params["v3_authtype"] = data.V3Authtype.ValueString()
	}
	if !data.V3Password.IsNull() && !data.V3Password.IsUnknown() {
		params["v3_password"] = data.V3Password.ValueString()
	}
	if !data.V3Privproto.IsNull() && !data.V3Privproto.IsUnknown() {
		params["v3_privproto"] = data.V3Privproto.ValueString()
	}
	if !data.V3Privpassphrase.IsNull() && !data.V3Privpassphrase.IsUnknown() {
		params["v3_privpassphrase"] = data.V3Privpassphrase.ValueString()
	}
	if !data.Loglevel.IsNull() && !data.Loglevel.IsUnknown() {
		params["loglevel"] = data.Loglevel.ValueInt64()
	}
	if !data.Options.IsNull() && !data.Options.IsUnknown() {
		params["options"] = data.Options.ValueString()
	}
	if !data.Zilstat.IsNull() && !data.Zilstat.IsUnknown() {
		params["zilstat"] = data.Zilstat.ValueBool()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("snmp.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply snmp config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("snmp.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read snmp config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["location"]; ok {
			if v == nil {
				data.Location = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Location = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Location = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["contact"]; ok {
			if v == nil {
				data.Contact = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Contact = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Contact = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["traps"]; ok {
			if v == nil {
				data.Traps = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Traps = types.BoolValue(b)
			}
		}
		if v, ok := m["v3"]; ok {
			if v == nil {
				data.V3 = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.V3 = types.BoolValue(b)
			}
		}
		if v, ok := m["community"]; ok {
			if v == nil {
				data.Community = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Community = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Community = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_username"]; ok {
			if v == nil {
				data.V3Username = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Username = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Username = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_authtype"]; ok {
			if v == nil {
				data.V3Authtype = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Authtype = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Authtype = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_password"]; ok {
			if v == nil {
				data.V3Password = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Password = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Password = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_privproto"]; ok {
			if v == nil {
				data.V3Privproto = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Privproto = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Privproto = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_privpassphrase"]; ok {
			if v == nil {
				data.V3Privpassphrase = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Privpassphrase = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Privpassphrase = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["loglevel"]; ok {
			if v == nil {
				data.Loglevel = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Loglevel = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["options"]; ok {
			if v == nil {
				data.Options = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Options = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Options = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["zilstat"]; ok {
			if v == nil {
				data.Zilstat = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Zilstat = types.BoolValue(b)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SnmpConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SnmpConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("snmp.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read snmp config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["location"]; ok {
			if v == nil {
				data.Location = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Location = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Location = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["contact"]; ok {
			if v == nil {
				data.Contact = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Contact = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Contact = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["traps"]; ok {
			if v == nil {
				data.Traps = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Traps = types.BoolValue(b)
			}
		}
		if v, ok := m["v3"]; ok {
			if v == nil {
				data.V3 = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.V3 = types.BoolValue(b)
			}
		}
		if v, ok := m["community"]; ok {
			if v == nil {
				data.Community = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Community = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Community = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_username"]; ok {
			if v == nil {
				data.V3Username = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Username = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Username = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_authtype"]; ok {
			if v == nil {
				data.V3Authtype = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Authtype = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Authtype = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_password"]; ok {
			if v == nil {
				data.V3Password = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Password = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Password = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_privproto"]; ok {
			if v == nil {
				data.V3Privproto = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Privproto = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Privproto = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_privpassphrase"]; ok {
			if v == nil {
				data.V3Privpassphrase = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Privpassphrase = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Privpassphrase = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["loglevel"]; ok {
			if v == nil {
				data.Loglevel = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Loglevel = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["options"]; ok {
			if v == nil {
				data.Options = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Options = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Options = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["zilstat"]; ok {
			if v == nil {
				data.Zilstat = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Zilstat = types.BoolValue(b)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SnmpConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SnmpConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		params["location"] = data.Location.ValueString()
	}
	if !data.Contact.IsNull() && !data.Contact.IsUnknown() {
		params["contact"] = data.Contact.ValueString()
	}
	if !data.Traps.IsNull() && !data.Traps.IsUnknown() {
		params["traps"] = data.Traps.ValueBool()
	}
	if !data.V3.IsNull() && !data.V3.IsUnknown() {
		params["v3"] = data.V3.ValueBool()
	}
	if !data.Community.IsNull() && !data.Community.IsUnknown() {
		params["community"] = data.Community.ValueString()
	}
	if !data.V3Username.IsNull() && !data.V3Username.IsUnknown() {
		params["v3_username"] = data.V3Username.ValueString()
	}
	if !data.V3Authtype.IsNull() && !data.V3Authtype.IsUnknown() {
		params["v3_authtype"] = data.V3Authtype.ValueString()
	}
	if !data.V3Password.IsNull() && !data.V3Password.IsUnknown() {
		params["v3_password"] = data.V3Password.ValueString()
	}
	if !data.V3Privproto.IsNull() && !data.V3Privproto.IsUnknown() {
		params["v3_privproto"] = data.V3Privproto.ValueString()
	}
	if !data.V3Privpassphrase.IsNull() && !data.V3Privpassphrase.IsUnknown() {
		params["v3_privpassphrase"] = data.V3Privpassphrase.ValueString()
	}
	if !data.Loglevel.IsNull() && !data.Loglevel.IsUnknown() {
		params["loglevel"] = data.Loglevel.ValueInt64()
	}
	if !data.Options.IsNull() && !data.Options.IsUnknown() {
		params["options"] = data.Options.ValueString()
	}
	if !data.Zilstat.IsNull() && !data.Zilstat.IsUnknown() {
		params["zilstat"] = data.Zilstat.ValueBool()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("snmp.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update snmp config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("snmp.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read snmp config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["location"]; ok {
			if v == nil {
				data.Location = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Location = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Location = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["contact"]; ok {
			if v == nil {
				data.Contact = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Contact = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Contact = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["traps"]; ok {
			if v == nil {
				data.Traps = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Traps = types.BoolValue(b)
			}
		}
		if v, ok := m["v3"]; ok {
			if v == nil {
				data.V3 = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.V3 = types.BoolValue(b)
			}
		}
		if v, ok := m["community"]; ok {
			if v == nil {
				data.Community = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Community = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Community = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_username"]; ok {
			if v == nil {
				data.V3Username = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Username = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Username = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_authtype"]; ok {
			if v == nil {
				data.V3Authtype = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Authtype = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Authtype = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_password"]; ok {
			if v == nil {
				data.V3Password = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Password = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Password = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_privproto"]; ok {
			if v == nil {
				data.V3Privproto = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Privproto = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Privproto = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["v3_privpassphrase"]; ok {
			if v == nil {
				data.V3Privpassphrase = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.V3Privpassphrase = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.V3Privpassphrase = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["loglevel"]; ok {
			if v == nil {
				data.Loglevel = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Loglevel = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["options"]; ok {
			if v == nil {
				data.Options = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Options = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Options = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["zilstat"]; ok {
			if v == nil {
				data.Zilstat = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Zilstat = types.BoolValue(b)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SnmpConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
