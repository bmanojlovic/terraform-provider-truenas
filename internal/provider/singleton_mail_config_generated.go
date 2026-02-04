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

type MailConfigResource struct {
	client *client.Client
}

type MailConfigResourceModel struct {
	Fromemail      types.String `tfsdk:"fromemail"`
	Fromname       types.String `tfsdk:"fromname"`
	Outgoingserver types.String `tfsdk:"outgoingserver"`
	Port           types.Int64  `tfsdk:"port"`
	Security       types.String `tfsdk:"security"`
	Smtp           types.Bool   `tfsdk:"smtp"`
	User           types.String `tfsdk:"user"`
	Pass           types.String `tfsdk:"pass"`
	Oauth          types.String `tfsdk:"oauth"`
}

func NewMailConfigResource() resource.Resource {
	return &MailConfigResource{}
}

func (r *MailConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_config"
}

func (r *MailConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Mail/SMTP configuration",
		Attributes: map[string]schema.Attribute{
			"fromemail":      schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The sending address that the mail server will use for sending emails."},
			"fromname":       schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Display name that will appear as the sender name in outgoing emails."},
			"outgoingserver": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Hostname or IP address of the SMTP server used for sending emails."},
			"port":           schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "TCP port number for the SMTP server connection."},
			"security":       schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Type of encryption."},
			"smtp":           schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether SMTP authentication is enabled and `user`, `pass` are required."},
			"user":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "SMTP username."},
			"pass":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "SMTP password."},
			"oauth":          schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "OAuth configuration for email providers that support it or `null` for basic authentication."},
		},
	}
}

func (r *MailConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MailConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Fromemail.IsNull() && !data.Fromemail.IsUnknown() {
		params["fromemail"] = data.Fromemail.ValueString()
	}
	if !data.Fromname.IsNull() && !data.Fromname.IsUnknown() {
		params["fromname"] = data.Fromname.ValueString()
	}
	if !data.Outgoingserver.IsNull() && !data.Outgoingserver.IsUnknown() {
		params["outgoingserver"] = data.Outgoingserver.ValueString()
	}
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		params["port"] = data.Port.ValueInt64()
	}
	if !data.Security.IsNull() && !data.Security.IsUnknown() {
		params["security"] = data.Security.ValueString()
	}
	if !data.Smtp.IsNull() && !data.Smtp.IsUnknown() {
		params["smtp"] = data.Smtp.ValueBool()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Pass.IsNull() && !data.Pass.IsUnknown() {
		params["pass"] = data.Pass.ValueString()
	}
	if !data.Oauth.IsNull() && !data.Oauth.IsUnknown() {
		params["oauth"] = data.Oauth.ValueString()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("mail.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply mail config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("mail.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read mail config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["fromemail"]; ok {
			if v == nil {
				data.Fromemail = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Fromemail = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Fromemail = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["fromname"]; ok {
			if v == nil {
				data.Fromname = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Fromname = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Fromname = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["outgoingserver"]; ok {
			if v == nil {
				data.Outgoingserver = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Outgoingserver = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Outgoingserver = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["port"]; ok {
			if v == nil {
				data.Port = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Port = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["security"]; ok {
			if v == nil {
				data.Security = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Security = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Security = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["smtp"]; ok {
			if v == nil {
				data.Smtp = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Smtp = types.BoolValue(b)
			}
		}
		if v, ok := m["user"]; ok {
			if v == nil {
				data.User = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.User = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.User = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["pass"]; ok {
			if v == nil {
				data.Pass = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Pass = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Pass = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["oauth"]; ok {
			if v == nil {
				data.Oauth = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Oauth = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Oauth = types.StringValue(string(j))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MailConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MailConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("mail.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read mail config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["fromemail"]; ok {
			if v == nil {
				data.Fromemail = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Fromemail = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Fromemail = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["fromname"]; ok {
			if v == nil {
				data.Fromname = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Fromname = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Fromname = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["outgoingserver"]; ok {
			if v == nil {
				data.Outgoingserver = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Outgoingserver = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Outgoingserver = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["port"]; ok {
			if v == nil {
				data.Port = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Port = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["security"]; ok {
			if v == nil {
				data.Security = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Security = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Security = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["smtp"]; ok {
			if v == nil {
				data.Smtp = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Smtp = types.BoolValue(b)
			}
		}
		if v, ok := m["user"]; ok {
			if v == nil {
				data.User = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.User = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.User = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["pass"]; ok {
			if v == nil {
				data.Pass = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Pass = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Pass = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["oauth"]; ok {
			if v == nil {
				data.Oauth = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Oauth = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Oauth = types.StringValue(string(j))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MailConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MailConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Fromemail.IsNull() && !data.Fromemail.IsUnknown() {
		params["fromemail"] = data.Fromemail.ValueString()
	}
	if !data.Fromname.IsNull() && !data.Fromname.IsUnknown() {
		params["fromname"] = data.Fromname.ValueString()
	}
	if !data.Outgoingserver.IsNull() && !data.Outgoingserver.IsUnknown() {
		params["outgoingserver"] = data.Outgoingserver.ValueString()
	}
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		params["port"] = data.Port.ValueInt64()
	}
	if !data.Security.IsNull() && !data.Security.IsUnknown() {
		params["security"] = data.Security.ValueString()
	}
	if !data.Smtp.IsNull() && !data.Smtp.IsUnknown() {
		params["smtp"] = data.Smtp.ValueBool()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Pass.IsNull() && !data.Pass.IsUnknown() {
		params["pass"] = data.Pass.ValueString()
	}
	if !data.Oauth.IsNull() && !data.Oauth.IsUnknown() {
		params["oauth"] = data.Oauth.ValueString()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("mail.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update mail config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("mail.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read mail config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["fromemail"]; ok {
			if v == nil {
				data.Fromemail = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Fromemail = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Fromemail = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["fromname"]; ok {
			if v == nil {
				data.Fromname = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Fromname = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Fromname = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["outgoingserver"]; ok {
			if v == nil {
				data.Outgoingserver = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Outgoingserver = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Outgoingserver = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["port"]; ok {
			if v == nil {
				data.Port = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Port = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["security"]; ok {
			if v == nil {
				data.Security = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Security = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Security = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["smtp"]; ok {
			if v == nil {
				data.Smtp = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Smtp = types.BoolValue(b)
			}
		}
		if v, ok := m["user"]; ok {
			if v == nil {
				data.User = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.User = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.User = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["pass"]; ok {
			if v == nil {
				data.Pass = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Pass = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Pass = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["oauth"]; ok {
			if v == nil {
				data.Oauth = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.Oauth = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.Oauth = types.StringValue(string(j))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MailConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
