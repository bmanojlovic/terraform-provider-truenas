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

type SmbConfigResource struct {
	client *client.Client
}

type SmbConfigResourceModel struct {
	Netbiosname types.String `tfsdk:"netbiosname"`
	Workgroup types.String `tfsdk:"workgroup"`
	Description types.String `tfsdk:"description"`
	EnableSmb1 types.Bool `tfsdk:"enable_smb1"`
	Unixcharset types.String `tfsdk:"unixcharset"`
	Localmaster types.Bool `tfsdk:"localmaster"`
	Syslog types.Bool `tfsdk:"syslog"`
	AaplExtensions types.Bool `tfsdk:"aapl_extensions"`
	AdminGroup types.String `tfsdk:"admin_group"`
	Guest types.String `tfsdk:"guest"`
	Filemask types.String `tfsdk:"filemask"`
	Dirmask types.String `tfsdk:"dirmask"`
	Ntlmv1Auth types.Bool `tfsdk:"ntlmv1_auth"`
	Multichannel types.Bool `tfsdk:"multichannel"`
	Encryption types.String `tfsdk:"encryption"`
	ServerSid types.String `tfsdk:"server_sid"`
	SmbOptions types.String `tfsdk:"smb_options"`
	Debug types.Bool `tfsdk:"debug"`
}

func NewSmbConfigResource() resource.Resource {
	return &SmbConfigResource{}
}

func (r *SmbConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_config"
}

func (r *SmbConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SMB service configuration",
		Attributes: map[string]schema.Attribute{
			"netbiosname": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The NetBIOS name of this server. "},
			"workgroup": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Workgroup name. When TrueNAS joins active directory, it automatically changes this value to match the NetBIOS     domain of the Active Directory domain. "},
			"description": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Description of the SMB server. SMB clients may see this description during some operations. "},
			"enable_smb1": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable SMB1 support on the server. WARNING: using the SMB1 protocol is not recommended. "},
			"unixcharset": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Select character set for file names on local filesystem. Use this option only if you know the names are not     UTF-8. "},
			"localmaster": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "When set to `true` the NetBIOS name server in TrueNAS participates in elections for the local master browser. When set to `false` the NetBIOS name server does not attempt to become a local master brow"},
			"syslog": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Send log messages to syslog. Enable this option if you want SMB server error logs to be included in     information sent to a remote syslog server. NOTE: This requires that remote syslog is globally c"},
			"aapl_extensions": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable support for SMB2/3 AAPL protocol extensions. This setting makes the TrueNAS server advertise support     for Apple protocol extensions as a MacOS server. Enabling this is required for Time Mach"},
			"admin_group": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The selected group has full administrator privileges on TrueNAS via the SMB protocol. "},
			"guest": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "SMB guest account username. This username provides access to legacy SMB shares with guest access enabled.     It must be a valid, existing local user account. "},
			"filemask": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "`smb.conf` create mask. DEFAULT applies current server default which is 664. "},
			"dirmask": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "`smb.conf` directory mask. DEFAULT applies current server default which is 775. "},
			"ntlmv1_auth": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable legacy and very insecure NTLMv1 authentication. This should never be done except     in extreme edge cases and may be against regulations in non-home environments. "},
			"multichannel": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Enable SMB3 multi-channel support. "},
			"encryption": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "SMB2/3 transport encryption setting for the TrueNAS SMB server.  * `NEGOTIATE`: Enable negotiation of data encryption. Encrypt data only if the client explicitly requests it. * `DESIRED`: Enable negot"},
			"server_sid": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "The unique identifier for the TrueNAS SMB server. It also serves as the domain SID for all local SMB user and     group accounts. "},
			"smb_options": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional unvalidated and unsupported configuration options for the SMB server. WARNING: Using `smb_options` may produce unexpected server behavior. "},
			"debug": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Set SMB log levels to debug. Use this setting only when troubleshooting a specific SMB issue. Do not use it     in production environments. "},
		},
	}
}

func (r *SmbConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SmbConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SmbConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Netbiosname.IsNull() && !data.Netbiosname.IsUnknown() { params["netbiosname"] = data.Netbiosname.ValueString() }
	if !data.Workgroup.IsNull() && !data.Workgroup.IsUnknown() { params["workgroup"] = data.Workgroup.ValueString() }
	if !data.Description.IsNull() && !data.Description.IsUnknown() { params["description"] = data.Description.ValueString() }
	if !data.EnableSmb1.IsNull() && !data.EnableSmb1.IsUnknown() { params["enable_smb1"] = data.EnableSmb1.ValueBool() }
	if !data.Unixcharset.IsNull() && !data.Unixcharset.IsUnknown() { params["unixcharset"] = data.Unixcharset.ValueString() }
	if !data.Localmaster.IsNull() && !data.Localmaster.IsUnknown() { params["localmaster"] = data.Localmaster.ValueBool() }
	if !data.Syslog.IsNull() && !data.Syslog.IsUnknown() { params["syslog"] = data.Syslog.ValueBool() }
	if !data.AaplExtensions.IsNull() && !data.AaplExtensions.IsUnknown() { params["aapl_extensions"] = data.AaplExtensions.ValueBool() }
	if !data.AdminGroup.IsNull() && !data.AdminGroup.IsUnknown() { params["admin_group"] = data.AdminGroup.ValueString() }
	if !data.Guest.IsNull() && !data.Guest.IsUnknown() { params["guest"] = data.Guest.ValueString() }
	if !data.Filemask.IsNull() && !data.Filemask.IsUnknown() { params["filemask"] = data.Filemask.ValueString() }
	if !data.Dirmask.IsNull() && !data.Dirmask.IsUnknown() { params["dirmask"] = data.Dirmask.ValueString() }
	if !data.Ntlmv1Auth.IsNull() && !data.Ntlmv1Auth.IsUnknown() { params["ntlmv1_auth"] = data.Ntlmv1Auth.ValueBool() }
	if !data.Multichannel.IsNull() && !data.Multichannel.IsUnknown() { params["multichannel"] = data.Multichannel.ValueBool() }
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() { params["encryption"] = data.Encryption.ValueString() }
	if !data.ServerSid.IsNull() && !data.ServerSid.IsUnknown() { params["server_sid"] = data.ServerSid.ValueString() }
	if !data.SmbOptions.IsNull() && !data.SmbOptions.IsUnknown() { params["smb_options"] = data.SmbOptions.ValueString() }
	if !data.Debug.IsNull() && !data.Debug.IsUnknown() { params["debug"] = data.Debug.ValueBool() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("smb.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply smb config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("smb.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read smb config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["netbiosname"]; ok { if v == nil { data.Netbiosname = types.StringNull() } else if s, ok := v.(string); ok { data.Netbiosname = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Netbiosname = types.StringValue(string(j)) } } }
		if v, ok := m["workgroup"]; ok { if v == nil { data.Workgroup = types.StringNull() } else if s, ok := v.(string); ok { data.Workgroup = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Workgroup = types.StringValue(string(j)) } } }
		if v, ok := m["description"]; ok { if v == nil { data.Description = types.StringNull() } else if s, ok := v.(string); ok { data.Description = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Description = types.StringValue(string(j)) } } }
		if v, ok := m["enable_smb1"]; ok { if v == nil { data.EnableSmb1 = types.BoolNull() } else if b, ok := v.(bool); ok { data.EnableSmb1 = types.BoolValue(b) } }
		if v, ok := m["unixcharset"]; ok { if v == nil { data.Unixcharset = types.StringNull() } else if s, ok := v.(string); ok { data.Unixcharset = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Unixcharset = types.StringValue(string(j)) } } }
		if v, ok := m["localmaster"]; ok { if v == nil { data.Localmaster = types.BoolNull() } else if b, ok := v.(bool); ok { data.Localmaster = types.BoolValue(b) } }
		if v, ok := m["syslog"]; ok { if v == nil { data.Syslog = types.BoolNull() } else if b, ok := v.(bool); ok { data.Syslog = types.BoolValue(b) } }
		if v, ok := m["aapl_extensions"]; ok { if v == nil { data.AaplExtensions = types.BoolNull() } else if b, ok := v.(bool); ok { data.AaplExtensions = types.BoolValue(b) } }
		if v, ok := m["admin_group"]; ok { if v == nil { data.AdminGroup = types.StringNull() } else if s, ok := v.(string); ok { data.AdminGroup = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.AdminGroup = types.StringValue(string(j)) } } }
		if v, ok := m["guest"]; ok { if v == nil { data.Guest = types.StringNull() } else if s, ok := v.(string); ok { data.Guest = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Guest = types.StringValue(string(j)) } } }
		if v, ok := m["filemask"]; ok { if v == nil { data.Filemask = types.StringNull() } else if s, ok := v.(string); ok { data.Filemask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Filemask = types.StringValue(string(j)) } } }
		if v, ok := m["dirmask"]; ok { if v == nil { data.Dirmask = types.StringNull() } else if s, ok := v.(string); ok { data.Dirmask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Dirmask = types.StringValue(string(j)) } } }
		if v, ok := m["ntlmv1_auth"]; ok { if v == nil { data.Ntlmv1Auth = types.BoolNull() } else if b, ok := v.(bool); ok { data.Ntlmv1Auth = types.BoolValue(b) } }
		if v, ok := m["multichannel"]; ok { if v == nil { data.Multichannel = types.BoolNull() } else if b, ok := v.(bool); ok { data.Multichannel = types.BoolValue(b) } }
		if v, ok := m["encryption"]; ok { if v == nil { data.Encryption = types.StringNull() } else if s, ok := v.(string); ok { data.Encryption = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Encryption = types.StringValue(string(j)) } } }
		if v, ok := m["server_sid"]; ok { if v == nil { data.ServerSid = types.StringNull() } else if s, ok := v.(string); ok { data.ServerSid = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.ServerSid = types.StringValue(string(j)) } } }
		if v, ok := m["smb_options"]; ok { if v == nil { data.SmbOptions = types.StringNull() } else if s, ok := v.(string); ok { data.SmbOptions = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.SmbOptions = types.StringValue(string(j)) } } }
		if v, ok := m["debug"]; ok { if v == nil { data.Debug = types.BoolNull() } else if b, ok := v.(bool); ok { data.Debug = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SmbConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SmbConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("smb.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read smb config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["netbiosname"]; ok { if v == nil { data.Netbiosname = types.StringNull() } else if s, ok := v.(string); ok { data.Netbiosname = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Netbiosname = types.StringValue(string(j)) } } }
		if v, ok := m["workgroup"]; ok { if v == nil { data.Workgroup = types.StringNull() } else if s, ok := v.(string); ok { data.Workgroup = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Workgroup = types.StringValue(string(j)) } } }
		if v, ok := m["description"]; ok { if v == nil { data.Description = types.StringNull() } else if s, ok := v.(string); ok { data.Description = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Description = types.StringValue(string(j)) } } }
		if v, ok := m["enable_smb1"]; ok { if v == nil { data.EnableSmb1 = types.BoolNull() } else if b, ok := v.(bool); ok { data.EnableSmb1 = types.BoolValue(b) } }
		if v, ok := m["unixcharset"]; ok { if v == nil { data.Unixcharset = types.StringNull() } else if s, ok := v.(string); ok { data.Unixcharset = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Unixcharset = types.StringValue(string(j)) } } }
		if v, ok := m["localmaster"]; ok { if v == nil { data.Localmaster = types.BoolNull() } else if b, ok := v.(bool); ok { data.Localmaster = types.BoolValue(b) } }
		if v, ok := m["syslog"]; ok { if v == nil { data.Syslog = types.BoolNull() } else if b, ok := v.(bool); ok { data.Syslog = types.BoolValue(b) } }
		if v, ok := m["aapl_extensions"]; ok { if v == nil { data.AaplExtensions = types.BoolNull() } else if b, ok := v.(bool); ok { data.AaplExtensions = types.BoolValue(b) } }
		if v, ok := m["admin_group"]; ok { if v == nil { data.AdminGroup = types.StringNull() } else if s, ok := v.(string); ok { data.AdminGroup = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.AdminGroup = types.StringValue(string(j)) } } }
		if v, ok := m["guest"]; ok { if v == nil { data.Guest = types.StringNull() } else if s, ok := v.(string); ok { data.Guest = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Guest = types.StringValue(string(j)) } } }
		if v, ok := m["filemask"]; ok { if v == nil { data.Filemask = types.StringNull() } else if s, ok := v.(string); ok { data.Filemask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Filemask = types.StringValue(string(j)) } } }
		if v, ok := m["dirmask"]; ok { if v == nil { data.Dirmask = types.StringNull() } else if s, ok := v.(string); ok { data.Dirmask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Dirmask = types.StringValue(string(j)) } } }
		if v, ok := m["ntlmv1_auth"]; ok { if v == nil { data.Ntlmv1Auth = types.BoolNull() } else if b, ok := v.(bool); ok { data.Ntlmv1Auth = types.BoolValue(b) } }
		if v, ok := m["multichannel"]; ok { if v == nil { data.Multichannel = types.BoolNull() } else if b, ok := v.(bool); ok { data.Multichannel = types.BoolValue(b) } }
		if v, ok := m["encryption"]; ok { if v == nil { data.Encryption = types.StringNull() } else if s, ok := v.(string); ok { data.Encryption = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Encryption = types.StringValue(string(j)) } } }
		if v, ok := m["server_sid"]; ok { if v == nil { data.ServerSid = types.StringNull() } else if s, ok := v.(string); ok { data.ServerSid = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.ServerSid = types.StringValue(string(j)) } } }
		if v, ok := m["smb_options"]; ok { if v == nil { data.SmbOptions = types.StringNull() } else if s, ok := v.(string); ok { data.SmbOptions = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.SmbOptions = types.StringValue(string(j)) } } }
		if v, ok := m["debug"]; ok { if v == nil { data.Debug = types.BoolNull() } else if b, ok := v.(bool); ok { data.Debug = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SmbConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SmbConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Netbiosname.IsNull() && !data.Netbiosname.IsUnknown() { params["netbiosname"] = data.Netbiosname.ValueString() }
	if !data.Workgroup.IsNull() && !data.Workgroup.IsUnknown() { params["workgroup"] = data.Workgroup.ValueString() }
	if !data.Description.IsNull() && !data.Description.IsUnknown() { params["description"] = data.Description.ValueString() }
	if !data.EnableSmb1.IsNull() && !data.EnableSmb1.IsUnknown() { params["enable_smb1"] = data.EnableSmb1.ValueBool() }
	if !data.Unixcharset.IsNull() && !data.Unixcharset.IsUnknown() { params["unixcharset"] = data.Unixcharset.ValueString() }
	if !data.Localmaster.IsNull() && !data.Localmaster.IsUnknown() { params["localmaster"] = data.Localmaster.ValueBool() }
	if !data.Syslog.IsNull() && !data.Syslog.IsUnknown() { params["syslog"] = data.Syslog.ValueBool() }
	if !data.AaplExtensions.IsNull() && !data.AaplExtensions.IsUnknown() { params["aapl_extensions"] = data.AaplExtensions.ValueBool() }
	if !data.AdminGroup.IsNull() && !data.AdminGroup.IsUnknown() { params["admin_group"] = data.AdminGroup.ValueString() }
	if !data.Guest.IsNull() && !data.Guest.IsUnknown() { params["guest"] = data.Guest.ValueString() }
	if !data.Filemask.IsNull() && !data.Filemask.IsUnknown() { params["filemask"] = data.Filemask.ValueString() }
	if !data.Dirmask.IsNull() && !data.Dirmask.IsUnknown() { params["dirmask"] = data.Dirmask.ValueString() }
	if !data.Ntlmv1Auth.IsNull() && !data.Ntlmv1Auth.IsUnknown() { params["ntlmv1_auth"] = data.Ntlmv1Auth.ValueBool() }
	if !data.Multichannel.IsNull() && !data.Multichannel.IsUnknown() { params["multichannel"] = data.Multichannel.ValueBool() }
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() { params["encryption"] = data.Encryption.ValueString() }
	if !data.ServerSid.IsNull() && !data.ServerSid.IsUnknown() { params["server_sid"] = data.ServerSid.ValueString() }
	if !data.SmbOptions.IsNull() && !data.SmbOptions.IsUnknown() { params["smb_options"] = data.SmbOptions.ValueString() }
	if !data.Debug.IsNull() && !data.Debug.IsUnknown() { params["debug"] = data.Debug.ValueBool() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("smb.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update smb config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("smb.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read smb config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["netbiosname"]; ok { if v == nil { data.Netbiosname = types.StringNull() } else if s, ok := v.(string); ok { data.Netbiosname = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Netbiosname = types.StringValue(string(j)) } } }
		if v, ok := m["workgroup"]; ok { if v == nil { data.Workgroup = types.StringNull() } else if s, ok := v.(string); ok { data.Workgroup = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Workgroup = types.StringValue(string(j)) } } }
		if v, ok := m["description"]; ok { if v == nil { data.Description = types.StringNull() } else if s, ok := v.(string); ok { data.Description = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Description = types.StringValue(string(j)) } } }
		if v, ok := m["enable_smb1"]; ok { if v == nil { data.EnableSmb1 = types.BoolNull() } else if b, ok := v.(bool); ok { data.EnableSmb1 = types.BoolValue(b) } }
		if v, ok := m["unixcharset"]; ok { if v == nil { data.Unixcharset = types.StringNull() } else if s, ok := v.(string); ok { data.Unixcharset = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Unixcharset = types.StringValue(string(j)) } } }
		if v, ok := m["localmaster"]; ok { if v == nil { data.Localmaster = types.BoolNull() } else if b, ok := v.(bool); ok { data.Localmaster = types.BoolValue(b) } }
		if v, ok := m["syslog"]; ok { if v == nil { data.Syslog = types.BoolNull() } else if b, ok := v.(bool); ok { data.Syslog = types.BoolValue(b) } }
		if v, ok := m["aapl_extensions"]; ok { if v == nil { data.AaplExtensions = types.BoolNull() } else if b, ok := v.(bool); ok { data.AaplExtensions = types.BoolValue(b) } }
		if v, ok := m["admin_group"]; ok { if v == nil { data.AdminGroup = types.StringNull() } else if s, ok := v.(string); ok { data.AdminGroup = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.AdminGroup = types.StringValue(string(j)) } } }
		if v, ok := m["guest"]; ok { if v == nil { data.Guest = types.StringNull() } else if s, ok := v.(string); ok { data.Guest = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Guest = types.StringValue(string(j)) } } }
		if v, ok := m["filemask"]; ok { if v == nil { data.Filemask = types.StringNull() } else if s, ok := v.(string); ok { data.Filemask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Filemask = types.StringValue(string(j)) } } }
		if v, ok := m["dirmask"]; ok { if v == nil { data.Dirmask = types.StringNull() } else if s, ok := v.(string); ok { data.Dirmask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Dirmask = types.StringValue(string(j)) } } }
		if v, ok := m["ntlmv1_auth"]; ok { if v == nil { data.Ntlmv1Auth = types.BoolNull() } else if b, ok := v.(bool); ok { data.Ntlmv1Auth = types.BoolValue(b) } }
		if v, ok := m["multichannel"]; ok { if v == nil { data.Multichannel = types.BoolNull() } else if b, ok := v.(bool); ok { data.Multichannel = types.BoolValue(b) } }
		if v, ok := m["encryption"]; ok { if v == nil { data.Encryption = types.StringNull() } else if s, ok := v.(string); ok { data.Encryption = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Encryption = types.StringValue(string(j)) } } }
		if v, ok := m["server_sid"]; ok { if v == nil { data.ServerSid = types.StringNull() } else if s, ok := v.(string); ok { data.ServerSid = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.ServerSid = types.StringValue(string(j)) } } }
		if v, ok := m["smb_options"]; ok { if v == nil { data.SmbOptions = types.StringNull() } else if s, ok := v.(string); ok { data.SmbOptions = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.SmbOptions = types.StringValue(string(j)) } } }
		if v, ok := m["debug"]; ok { if v == nil { data.Debug = types.BoolNull() } else if b, ok := v.(bool); ok { data.Debug = types.BoolValue(b) } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SmbConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
