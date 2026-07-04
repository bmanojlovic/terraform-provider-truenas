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

type SshConfigResource struct {
	client *client.Client
}

type SshConfigResourceModel struct {
	Tcpport         types.Int64  `tfsdk:"tcpport"`
	Passwordauth    types.Bool   `tfsdk:"passwordauth"`
	Kerberosauth    types.Bool   `tfsdk:"kerberosauth"`
	Tcpfwd          types.Bool   `tfsdk:"tcpfwd"`
	Compression     types.Bool   `tfsdk:"compression"`
	SftpLogLevel    types.String `tfsdk:"sftp_log_level"`
	SftpLogFacility types.String `tfsdk:"sftp_log_facility"`
	Options         types.String `tfsdk:"options"`
}

func NewSshConfigResource() resource.Resource {
	return &SshConfigResource{}
}

func (r *SshConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_config"
}

func (r *SshConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH service configuration",
		Attributes: map[string]schema.Attribute{
			"tcpport":           schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "TCP port number for SSH connections."},
			"passwordauth":      schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether password authentication is enabled."},
			"kerberosauth":      schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether Kerberos authentication is enabled."},
			"tcpfwd":            schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether TCP forwarding is enabled."},
			"compression":       schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether compression is enabled for SSH connections."},
			"sftp_log_level":    schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Logging level for SFTP subsystem (empty string means default)."},
			"sftp_log_facility": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Syslog facility for SFTP logging (empty string means default)."},
			"options":           schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional SSH daemon configuration options."},
		},
	}
}

func (r *SshConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SshConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SshConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Tcpport.IsNull() && !data.Tcpport.IsUnknown() {
		params["tcpport"] = data.Tcpport.ValueInt64()
	}
	if !data.Passwordauth.IsNull() && !data.Passwordauth.IsUnknown() {
		params["passwordauth"] = data.Passwordauth.ValueBool()
	}
	if !data.Kerberosauth.IsNull() && !data.Kerberosauth.IsUnknown() {
		params["kerberosauth"] = data.Kerberosauth.ValueBool()
	}
	if !data.Tcpfwd.IsNull() && !data.Tcpfwd.IsUnknown() {
		params["tcpfwd"] = data.Tcpfwd.ValueBool()
	}
	if !data.Compression.IsNull() && !data.Compression.IsUnknown() {
		params["compression"] = data.Compression.ValueBool()
	}
	if !data.SftpLogLevel.IsNull() && !data.SftpLogLevel.IsUnknown() {
		params["sftp_log_level"] = data.SftpLogLevel.ValueString()
	}
	if !data.SftpLogFacility.IsNull() && !data.SftpLogFacility.IsUnknown() {
		params["sftp_log_facility"] = data.SftpLogFacility.ValueString()
	}
	if !data.Options.IsNull() && !data.Options.IsUnknown() {
		params["options"] = data.Options.ValueString()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("ssh.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply ssh config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("ssh.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ssh config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["tcpport"]; ok {
			if v == nil {
				data.Tcpport = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Tcpport = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["passwordauth"]; ok {
			if v == nil {
				data.Passwordauth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Passwordauth = types.BoolValue(b)
			}
		}
		if v, ok := m["kerberosauth"]; ok {
			if v == nil {
				data.Kerberosauth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Kerberosauth = types.BoolValue(b)
			}
		}
		if v, ok := m["tcpfwd"]; ok {
			if v == nil {
				data.Tcpfwd = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Tcpfwd = types.BoolValue(b)
			}
		}
		if v, ok := m["compression"]; ok {
			if v == nil {
				data.Compression = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Compression = types.BoolValue(b)
			}
		}
		if v, ok := m["sftp_log_level"]; ok {
			if v == nil {
				data.SftpLogLevel = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SftpLogLevel = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SftpLogLevel = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sftp_log_facility"]; ok {
			if v == nil {
				data.SftpLogFacility = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SftpLogFacility = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SftpLogFacility = types.StringValue(string(j))
				}
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
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SshConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("ssh.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ssh config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["tcpport"]; ok {
			if v == nil {
				data.Tcpport = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Tcpport = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["passwordauth"]; ok {
			if v == nil {
				data.Passwordauth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Passwordauth = types.BoolValue(b)
			}
		}
		if v, ok := m["kerberosauth"]; ok {
			if v == nil {
				data.Kerberosauth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Kerberosauth = types.BoolValue(b)
			}
		}
		if v, ok := m["tcpfwd"]; ok {
			if v == nil {
				data.Tcpfwd = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Tcpfwd = types.BoolValue(b)
			}
		}
		if v, ok := m["compression"]; ok {
			if v == nil {
				data.Compression = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Compression = types.BoolValue(b)
			}
		}
		if v, ok := m["sftp_log_level"]; ok {
			if v == nil {
				data.SftpLogLevel = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SftpLogLevel = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SftpLogLevel = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sftp_log_facility"]; ok {
			if v == nil {
				data.SftpLogFacility = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SftpLogFacility = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SftpLogFacility = types.StringValue(string(j))
				}
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
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SshConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Tcpport.IsNull() && !data.Tcpport.IsUnknown() {
		params["tcpport"] = data.Tcpport.ValueInt64()
	}
	if !data.Passwordauth.IsNull() && !data.Passwordauth.IsUnknown() {
		params["passwordauth"] = data.Passwordauth.ValueBool()
	}
	if !data.Kerberosauth.IsNull() && !data.Kerberosauth.IsUnknown() {
		params["kerberosauth"] = data.Kerberosauth.ValueBool()
	}
	if !data.Tcpfwd.IsNull() && !data.Tcpfwd.IsUnknown() {
		params["tcpfwd"] = data.Tcpfwd.ValueBool()
	}
	if !data.Compression.IsNull() && !data.Compression.IsUnknown() {
		params["compression"] = data.Compression.ValueBool()
	}
	if !data.SftpLogLevel.IsNull() && !data.SftpLogLevel.IsUnknown() {
		params["sftp_log_level"] = data.SftpLogLevel.ValueString()
	}
	if !data.SftpLogFacility.IsNull() && !data.SftpLogFacility.IsUnknown() {
		params["sftp_log_facility"] = data.SftpLogFacility.ValueString()
	}
	if !data.Options.IsNull() && !data.Options.IsUnknown() {
		params["options"] = data.Options.ValueString()
	}
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("ssh.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update ssh config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("ssh.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ssh config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["tcpport"]; ok {
			if v == nil {
				data.Tcpport = types.Int64Null()
			} else if f, ok := v.(float64); ok {
				data.Tcpport = types.Int64Value(int64(f))
			}
		}
		if v, ok := m["passwordauth"]; ok {
			if v == nil {
				data.Passwordauth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Passwordauth = types.BoolValue(b)
			}
		}
		if v, ok := m["kerberosauth"]; ok {
			if v == nil {
				data.Kerberosauth = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Kerberosauth = types.BoolValue(b)
			}
		}
		if v, ok := m["tcpfwd"]; ok {
			if v == nil {
				data.Tcpfwd = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Tcpfwd = types.BoolValue(b)
			}
		}
		if v, ok := m["compression"]; ok {
			if v == nil {
				data.Compression = types.BoolNull()
			} else if b, ok := v.(bool); ok {
				data.Compression = types.BoolValue(b)
			}
		}
		if v, ok := m["sftp_log_level"]; ok {
			if v == nil {
				data.SftpLogLevel = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SftpLogLevel = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SftpLogLevel = types.StringValue(string(j))
				}
			}
		}
		if v, ok := m["sftp_log_facility"]; ok {
			if v == nil {
				data.SftpLogFacility = types.StringNull()
			} else if s, ok := v.(string); ok {
				data.SftpLogFacility = types.StringValue(s)
			} else {
				if j, e := json.Marshal(v); e == nil {
					data.SftpLogFacility = types.StringValue(string(j))
				}
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
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
