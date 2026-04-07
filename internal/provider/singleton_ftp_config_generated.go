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

type FtpConfigResource struct {
	client *client.Client
}

type FtpConfigResourceModel struct {
	Port types.Int64 `tfsdk:"port"`
	Clients types.Int64 `tfsdk:"clients"`
	Ipconnections types.Int64 `tfsdk:"ipconnections"`
	Loginattempt types.Int64 `tfsdk:"loginattempt"`
	Timeout types.Int64 `tfsdk:"timeout"`
	TimeoutNotransfer types.Int64 `tfsdk:"timeout_notransfer"`
	Onlyanonymous types.Bool `tfsdk:"onlyanonymous"`
	Anonpath types.String `tfsdk:"anonpath"`
	Onlylocal types.Bool `tfsdk:"onlylocal"`
	Banner types.String `tfsdk:"banner"`
	Filemask types.String `tfsdk:"filemask"`
	Dirmask types.String `tfsdk:"dirmask"`
	Fxp types.Bool `tfsdk:"fxp"`
	Resume types.Bool `tfsdk:"resume"`
	Defaultroot types.Bool `tfsdk:"defaultroot"`
	Ident types.Bool `tfsdk:"ident"`
	Reversedns types.Bool `tfsdk:"reversedns"`
	Masqaddress types.String `tfsdk:"masqaddress"`
	Passiveportsmin types.Int64 `tfsdk:"passiveportsmin"`
	Passiveportsmax types.Int64 `tfsdk:"passiveportsmax"`
	Localuserbw types.Int64 `tfsdk:"localuserbw"`
	Localuserdlbw types.Int64 `tfsdk:"localuserdlbw"`
	Anonuserbw types.Int64 `tfsdk:"anonuserbw"`
	Anonuserdlbw types.Int64 `tfsdk:"anonuserdlbw"`
	Tls types.Bool `tfsdk:"tls"`
	TlsPolicy types.String `tfsdk:"tls_policy"`
	TlsOptAllowClientRenegotiations types.Bool `tfsdk:"tls_opt_allow_client_renegotiations"`
	TlsOptAllowDotLogin types.Bool `tfsdk:"tls_opt_allow_dot_login"`
	TlsOptAllowPerUser types.Bool `tfsdk:"tls_opt_allow_per_user"`
	TlsOptCommonNameRequired types.Bool `tfsdk:"tls_opt_common_name_required"`
	TlsOptEnableDiags types.Bool `tfsdk:"tls_opt_enable_diags"`
	TlsOptExportCertData types.Bool `tfsdk:"tls_opt_export_cert_data"`
	TlsOptNoEmptyFragments types.Bool `tfsdk:"tls_opt_no_empty_fragments"`
	TlsOptNoSessionReuseRequired types.Bool `tfsdk:"tls_opt_no_session_reuse_required"`
	TlsOptStdenvvars types.Bool `tfsdk:"tls_opt_stdenvvars"`
	TlsOptDnsNameRequired types.Bool `tfsdk:"tls_opt_dns_name_required"`
	TlsOptIpAddressRequired types.Bool `tfsdk:"tls_opt_ip_address_required"`
	SsltlsCertificate types.Int64 `tfsdk:"ssltls_certificate"`
	Options types.String `tfsdk:"options"`
}

func NewFtpConfigResource() resource.Resource {
	return &FtpConfigResource{}
}

func (r *FtpConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ftp_config"
}

func (r *FtpConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "FTP service configuration",
		Attributes: map[string]schema.Attribute{
			"port": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "TCP port number on which the FTP service listens for incoming connections."},
			"clients": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum number of simultaneous client connections allowed."},
			"ipconnections": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum number of connections allowed from a single IP address. 0 means unlimited."},
			"loginattempt": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum number of failed login attempts before blocking an IP address. 0 disables this limit."},
			"timeout": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Idle timeout in seconds before disconnecting inactive clients. 0 disables timeout."},
			"timeout_notransfer": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Timeout in seconds for clients that connect but do not transfer data. 0 disables timeout."},
			"onlyanonymous": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to allow only anonymous FTP access, disabling authenticated user login."},
			"anonpath": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Filesystem path for anonymous FTP users. `null` to use the default anonymous FTP directory."},
			"onlylocal": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to allow only local system users to login, disabling anonymous access."},
			"banner": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Welcome message displayed to FTP clients upon connection."},
			"filemask": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Default Unix permissions (umask) for files created by FTP users."},
			"dirmask": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Default Unix permissions (umask) for directories created by FTP users."},
			"fxp": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to enable File eXchange Protocol (FXP) for server-to-server transfers."},
			"resume": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to allow clients to resume interrupted file transfers."},
			"defaultroot": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to restrict users to their home directories (chroot jail)."},
			"ident": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to perform RFC 1413 ident lookups on connecting clients."},
			"reversedns": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to perform reverse DNS lookups on client IP addresses for logging."},
			"masqaddress": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Public IP address to advertise to clients for passive mode connections when behind NAT."},
			"passiveportsmin": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Minimum port number for passive mode data connections. Must be 0 or between 1024-65535."},
			"passiveportsmax": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum port number for passive mode data connections. Must be 0 or between 1024-65535."},
			"localuserbw": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum upload bandwidth in KiB/s for local users. 0 means unlimited."},
			"localuserdlbw": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum download bandwidth in KiB/s for local users. 0 means unlimited."},
			"anonuserbw": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum upload bandwidth in KiB/s for anonymous users. 0 means unlimited."},
			"anonuserdlbw": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "Maximum download bandwidth in KiB/s for anonymous users. 0 means unlimited."},
			"tls": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to enable TLS/SSL encryption for FTP connections."},
			"tls_policy": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "TLS policy for connections. Values include: `\"on\"` (required), `\"off\"` (disabled), `\"data\"` (data only),     `\"auth\"` (authentication only), `\"ctrl\"` (control only), or combinations with `+`"},
			"tls_opt_allow_client_renegotiations": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to allow TLS clients to initiate renegotiation of the TLS connection."},
			"tls_opt_allow_dot_login": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to allow .ftpaccess files to override TLS requirements for specific users."},
			"tls_opt_allow_per_user": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to allow per-user TLS configuration overrides."},
			"tls_opt_common_name_required": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to require client certificates to have a Common Name field."},
			"tls_opt_enable_diags": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to enable detailed TLS diagnostic logging."},
			"tls_opt_export_cert_data": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to export client certificate data to environment variables."},
			"tls_opt_no_empty_fragments": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to disable empty TLS record fragments to improve compatibility with some clients.      Disabling increases vulnerability to some attack vectors."},
			"tls_opt_no_session_reuse_required": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to disable the requirement for TLS session reuse."},
			"tls_opt_stdenvvars": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to export standard TLS environment variables for use by external programs."},
			"tls_opt_dns_name_required": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to require client certificates to contain a DNS name in the Subject Alternative Name extension.     The `reversedns` setting must also be enabled."},
			"tls_opt_ip_address_required": schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether to require client certificates to contain an IP address in the Subject Alternative Name extension."},
			"ssltls_certificate": schema.Int64Attribute{Optional: true, Computed: true, MarkdownDescription: "ID of the certificate to use for TLS/SSL connections. `null` to use the default system certificate."},
			"options": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Additional ProFTPD configuration directives to include in the server configuration.     Manual directives may render the FTP service non-functional and should be used with caution."},
		},
	}
}

func (r *FtpConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FtpConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FtpConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Port.IsNull() && !data.Port.IsUnknown() { params["port"] = data.Port.ValueInt64() }
	if !data.Clients.IsNull() && !data.Clients.IsUnknown() { params["clients"] = data.Clients.ValueInt64() }
	if !data.Ipconnections.IsNull() && !data.Ipconnections.IsUnknown() { params["ipconnections"] = data.Ipconnections.ValueInt64() }
	if !data.Loginattempt.IsNull() && !data.Loginattempt.IsUnknown() { params["loginattempt"] = data.Loginattempt.ValueInt64() }
	if !data.Timeout.IsNull() && !data.Timeout.IsUnknown() { params["timeout"] = data.Timeout.ValueInt64() }
	if !data.TimeoutNotransfer.IsNull() && !data.TimeoutNotransfer.IsUnknown() { params["timeout_notransfer"] = data.TimeoutNotransfer.ValueInt64() }
	if !data.Onlyanonymous.IsNull() && !data.Onlyanonymous.IsUnknown() { params["onlyanonymous"] = data.Onlyanonymous.ValueBool() }
	if !data.Anonpath.IsNull() && !data.Anonpath.IsUnknown() { params["anonpath"] = data.Anonpath.ValueString() }
	if !data.Onlylocal.IsNull() && !data.Onlylocal.IsUnknown() { params["onlylocal"] = data.Onlylocal.ValueBool() }
	if !data.Banner.IsNull() && !data.Banner.IsUnknown() { params["banner"] = data.Banner.ValueString() }
	if !data.Filemask.IsNull() && !data.Filemask.IsUnknown() { params["filemask"] = data.Filemask.ValueString() }
	if !data.Dirmask.IsNull() && !data.Dirmask.IsUnknown() { params["dirmask"] = data.Dirmask.ValueString() }
	if !data.Fxp.IsNull() && !data.Fxp.IsUnknown() { params["fxp"] = data.Fxp.ValueBool() }
	if !data.Resume.IsNull() && !data.Resume.IsUnknown() { params["resume"] = data.Resume.ValueBool() }
	if !data.Defaultroot.IsNull() && !data.Defaultroot.IsUnknown() { params["defaultroot"] = data.Defaultroot.ValueBool() }
	if !data.Ident.IsNull() && !data.Ident.IsUnknown() { params["ident"] = data.Ident.ValueBool() }
	if !data.Reversedns.IsNull() && !data.Reversedns.IsUnknown() { params["reversedns"] = data.Reversedns.ValueBool() }
	if !data.Masqaddress.IsNull() && !data.Masqaddress.IsUnknown() { params["masqaddress"] = data.Masqaddress.ValueString() }
	if !data.Passiveportsmin.IsNull() && !data.Passiveportsmin.IsUnknown() { params["passiveportsmin"] = data.Passiveportsmin.ValueInt64() }
	if !data.Passiveportsmax.IsNull() && !data.Passiveportsmax.IsUnknown() { params["passiveportsmax"] = data.Passiveportsmax.ValueInt64() }
	if !data.Localuserbw.IsNull() && !data.Localuserbw.IsUnknown() { params["localuserbw"] = data.Localuserbw.ValueInt64() }
	if !data.Localuserdlbw.IsNull() && !data.Localuserdlbw.IsUnknown() { params["localuserdlbw"] = data.Localuserdlbw.ValueInt64() }
	if !data.Anonuserbw.IsNull() && !data.Anonuserbw.IsUnknown() { params["anonuserbw"] = data.Anonuserbw.ValueInt64() }
	if !data.Anonuserdlbw.IsNull() && !data.Anonuserdlbw.IsUnknown() { params["anonuserdlbw"] = data.Anonuserdlbw.ValueInt64() }
	if !data.Tls.IsNull() && !data.Tls.IsUnknown() { params["tls"] = data.Tls.ValueBool() }
	if !data.TlsPolicy.IsNull() && !data.TlsPolicy.IsUnknown() { params["tls_policy"] = data.TlsPolicy.ValueString() }
	if !data.TlsOptAllowClientRenegotiations.IsNull() && !data.TlsOptAllowClientRenegotiations.IsUnknown() { params["tls_opt_allow_client_renegotiations"] = data.TlsOptAllowClientRenegotiations.ValueBool() }
	if !data.TlsOptAllowDotLogin.IsNull() && !data.TlsOptAllowDotLogin.IsUnknown() { params["tls_opt_allow_dot_login"] = data.TlsOptAllowDotLogin.ValueBool() }
	if !data.TlsOptAllowPerUser.IsNull() && !data.TlsOptAllowPerUser.IsUnknown() { params["tls_opt_allow_per_user"] = data.TlsOptAllowPerUser.ValueBool() }
	if !data.TlsOptCommonNameRequired.IsNull() && !data.TlsOptCommonNameRequired.IsUnknown() { params["tls_opt_common_name_required"] = data.TlsOptCommonNameRequired.ValueBool() }
	if !data.TlsOptEnableDiags.IsNull() && !data.TlsOptEnableDiags.IsUnknown() { params["tls_opt_enable_diags"] = data.TlsOptEnableDiags.ValueBool() }
	if !data.TlsOptExportCertData.IsNull() && !data.TlsOptExportCertData.IsUnknown() { params["tls_opt_export_cert_data"] = data.TlsOptExportCertData.ValueBool() }
	if !data.TlsOptNoEmptyFragments.IsNull() && !data.TlsOptNoEmptyFragments.IsUnknown() { params["tls_opt_no_empty_fragments"] = data.TlsOptNoEmptyFragments.ValueBool() }
	if !data.TlsOptNoSessionReuseRequired.IsNull() && !data.TlsOptNoSessionReuseRequired.IsUnknown() { params["tls_opt_no_session_reuse_required"] = data.TlsOptNoSessionReuseRequired.ValueBool() }
	if !data.TlsOptStdenvvars.IsNull() && !data.TlsOptStdenvvars.IsUnknown() { params["tls_opt_stdenvvars"] = data.TlsOptStdenvvars.ValueBool() }
	if !data.TlsOptDnsNameRequired.IsNull() && !data.TlsOptDnsNameRequired.IsUnknown() { params["tls_opt_dns_name_required"] = data.TlsOptDnsNameRequired.ValueBool() }
	if !data.TlsOptIpAddressRequired.IsNull() && !data.TlsOptIpAddressRequired.IsUnknown() { params["tls_opt_ip_address_required"] = data.TlsOptIpAddressRequired.ValueBool() }
	if !data.SsltlsCertificate.IsNull() && !data.SsltlsCertificate.IsUnknown() { params["ssltls_certificate"] = data.SsltlsCertificate.ValueInt64() }
	if !data.Options.IsNull() && !data.Options.IsUnknown() { params["options"] = data.Options.ValueString() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("ftp.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Configuration Failed", fmt.Sprintf("Failed to apply ftp config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("ftp.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ftp config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["port"]; ok { if v == nil { data.Port = types.Int64Null() } else if f, ok := v.(float64); ok { data.Port = types.Int64Value(int64(f)) } }
		if v, ok := m["clients"]; ok { if v == nil { data.Clients = types.Int64Null() } else if f, ok := v.(float64); ok { data.Clients = types.Int64Value(int64(f)) } }
		if v, ok := m["ipconnections"]; ok { if v == nil { data.Ipconnections = types.Int64Null() } else if f, ok := v.(float64); ok { data.Ipconnections = types.Int64Value(int64(f)) } }
		if v, ok := m["loginattempt"]; ok { if v == nil { data.Loginattempt = types.Int64Null() } else if f, ok := v.(float64); ok { data.Loginattempt = types.Int64Value(int64(f)) } }
		if v, ok := m["timeout"]; ok { if v == nil { data.Timeout = types.Int64Null() } else if f, ok := v.(float64); ok { data.Timeout = types.Int64Value(int64(f)) } }
		if v, ok := m["timeout_notransfer"]; ok { if v == nil { data.TimeoutNotransfer = types.Int64Null() } else if f, ok := v.(float64); ok { data.TimeoutNotransfer = types.Int64Value(int64(f)) } }
		if v, ok := m["onlyanonymous"]; ok { if v == nil { data.Onlyanonymous = types.BoolNull() } else if b, ok := v.(bool); ok { data.Onlyanonymous = types.BoolValue(b) } }
		if v, ok := m["anonpath"]; ok { if v == nil { data.Anonpath = types.StringNull() } else if s, ok := v.(string); ok { data.Anonpath = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Anonpath = types.StringValue(string(j)) } } }
		if v, ok := m["onlylocal"]; ok { if v == nil { data.Onlylocal = types.BoolNull() } else if b, ok := v.(bool); ok { data.Onlylocal = types.BoolValue(b) } }
		if v, ok := m["banner"]; ok { if v == nil { data.Banner = types.StringNull() } else if s, ok := v.(string); ok { data.Banner = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Banner = types.StringValue(string(j)) } } }
		if v, ok := m["filemask"]; ok { if v == nil { data.Filemask = types.StringNull() } else if s, ok := v.(string); ok { data.Filemask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Filemask = types.StringValue(string(j)) } } }
		if v, ok := m["dirmask"]; ok { if v == nil { data.Dirmask = types.StringNull() } else if s, ok := v.(string); ok { data.Dirmask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Dirmask = types.StringValue(string(j)) } } }
		if v, ok := m["fxp"]; ok { if v == nil { data.Fxp = types.BoolNull() } else if b, ok := v.(bool); ok { data.Fxp = types.BoolValue(b) } }
		if v, ok := m["resume"]; ok { if v == nil { data.Resume = types.BoolNull() } else if b, ok := v.(bool); ok { data.Resume = types.BoolValue(b) } }
		if v, ok := m["defaultroot"]; ok { if v == nil { data.Defaultroot = types.BoolNull() } else if b, ok := v.(bool); ok { data.Defaultroot = types.BoolValue(b) } }
		if v, ok := m["ident"]; ok { if v == nil { data.Ident = types.BoolNull() } else if b, ok := v.(bool); ok { data.Ident = types.BoolValue(b) } }
		if v, ok := m["reversedns"]; ok { if v == nil { data.Reversedns = types.BoolNull() } else if b, ok := v.(bool); ok { data.Reversedns = types.BoolValue(b) } }
		if v, ok := m["masqaddress"]; ok { if v == nil { data.Masqaddress = types.StringNull() } else if s, ok := v.(string); ok { data.Masqaddress = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Masqaddress = types.StringValue(string(j)) } } }
		if v, ok := m["passiveportsmin"]; ok { if v == nil { data.Passiveportsmin = types.Int64Null() } else if f, ok := v.(float64); ok { data.Passiveportsmin = types.Int64Value(int64(f)) } }
		if v, ok := m["passiveportsmax"]; ok { if v == nil { data.Passiveportsmax = types.Int64Null() } else if f, ok := v.(float64); ok { data.Passiveportsmax = types.Int64Value(int64(f)) } }
		if v, ok := m["localuserbw"]; ok { if v == nil { data.Localuserbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Localuserbw = types.Int64Value(int64(f)) } }
		if v, ok := m["localuserdlbw"]; ok { if v == nil { data.Localuserdlbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Localuserdlbw = types.Int64Value(int64(f)) } }
		if v, ok := m["anonuserbw"]; ok { if v == nil { data.Anonuserbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Anonuserbw = types.Int64Value(int64(f)) } }
		if v, ok := m["anonuserdlbw"]; ok { if v == nil { data.Anonuserdlbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Anonuserdlbw = types.Int64Value(int64(f)) } }
		if v, ok := m["tls"]; ok { if v == nil { data.Tls = types.BoolNull() } else if b, ok := v.(bool); ok { data.Tls = types.BoolValue(b) } }
		if v, ok := m["tls_policy"]; ok { if v == nil { data.TlsPolicy = types.StringNull() } else if s, ok := v.(string); ok { data.TlsPolicy = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.TlsPolicy = types.StringValue(string(j)) } } }
		if v, ok := m["tls_opt_allow_client_renegotiations"]; ok { if v == nil { data.TlsOptAllowClientRenegotiations = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowClientRenegotiations = types.BoolValue(b) } }
		if v, ok := m["tls_opt_allow_dot_login"]; ok { if v == nil { data.TlsOptAllowDotLogin = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowDotLogin = types.BoolValue(b) } }
		if v, ok := m["tls_opt_allow_per_user"]; ok { if v == nil { data.TlsOptAllowPerUser = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowPerUser = types.BoolValue(b) } }
		if v, ok := m["tls_opt_common_name_required"]; ok { if v == nil { data.TlsOptCommonNameRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptCommonNameRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_enable_diags"]; ok { if v == nil { data.TlsOptEnableDiags = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptEnableDiags = types.BoolValue(b) } }
		if v, ok := m["tls_opt_export_cert_data"]; ok { if v == nil { data.TlsOptExportCertData = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptExportCertData = types.BoolValue(b) } }
		if v, ok := m["tls_opt_no_empty_fragments"]; ok { if v == nil { data.TlsOptNoEmptyFragments = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptNoEmptyFragments = types.BoolValue(b) } }
		if v, ok := m["tls_opt_no_session_reuse_required"]; ok { if v == nil { data.TlsOptNoSessionReuseRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptNoSessionReuseRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_stdenvvars"]; ok { if v == nil { data.TlsOptStdenvvars = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptStdenvvars = types.BoolValue(b) } }
		if v, ok := m["tls_opt_dns_name_required"]; ok { if v == nil { data.TlsOptDnsNameRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptDnsNameRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_ip_address_required"]; ok { if v == nil { data.TlsOptIpAddressRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptIpAddressRequired = types.BoolValue(b) } }
		if v, ok := m["ssltls_certificate"]; ok { if v == nil { data.SsltlsCertificate = types.Int64Null() } else if f, ok := v.(float64); ok { data.SsltlsCertificate = types.Int64Value(int64(f)) } }
		if v, ok := m["options"]; ok { if v == nil { data.Options = types.StringNull() } else if s, ok := v.(string); ok { data.Options = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Options = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FtpConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FtpConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current configuration
	result, err := r.client.Call("ftp.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ftp config: %s", err.Error()))
		return
	}

	// Map result to state
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["port"]; ok { if v == nil { data.Port = types.Int64Null() } else if f, ok := v.(float64); ok { data.Port = types.Int64Value(int64(f)) } }
		if v, ok := m["clients"]; ok { if v == nil { data.Clients = types.Int64Null() } else if f, ok := v.(float64); ok { data.Clients = types.Int64Value(int64(f)) } }
		if v, ok := m["ipconnections"]; ok { if v == nil { data.Ipconnections = types.Int64Null() } else if f, ok := v.(float64); ok { data.Ipconnections = types.Int64Value(int64(f)) } }
		if v, ok := m["loginattempt"]; ok { if v == nil { data.Loginattempt = types.Int64Null() } else if f, ok := v.(float64); ok { data.Loginattempt = types.Int64Value(int64(f)) } }
		if v, ok := m["timeout"]; ok { if v == nil { data.Timeout = types.Int64Null() } else if f, ok := v.(float64); ok { data.Timeout = types.Int64Value(int64(f)) } }
		if v, ok := m["timeout_notransfer"]; ok { if v == nil { data.TimeoutNotransfer = types.Int64Null() } else if f, ok := v.(float64); ok { data.TimeoutNotransfer = types.Int64Value(int64(f)) } }
		if v, ok := m["onlyanonymous"]; ok { if v == nil { data.Onlyanonymous = types.BoolNull() } else if b, ok := v.(bool); ok { data.Onlyanonymous = types.BoolValue(b) } }
		if v, ok := m["anonpath"]; ok { if v == nil { data.Anonpath = types.StringNull() } else if s, ok := v.(string); ok { data.Anonpath = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Anonpath = types.StringValue(string(j)) } } }
		if v, ok := m["onlylocal"]; ok { if v == nil { data.Onlylocal = types.BoolNull() } else if b, ok := v.(bool); ok { data.Onlylocal = types.BoolValue(b) } }
		if v, ok := m["banner"]; ok { if v == nil { data.Banner = types.StringNull() } else if s, ok := v.(string); ok { data.Banner = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Banner = types.StringValue(string(j)) } } }
		if v, ok := m["filemask"]; ok { if v == nil { data.Filemask = types.StringNull() } else if s, ok := v.(string); ok { data.Filemask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Filemask = types.StringValue(string(j)) } } }
		if v, ok := m["dirmask"]; ok { if v == nil { data.Dirmask = types.StringNull() } else if s, ok := v.(string); ok { data.Dirmask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Dirmask = types.StringValue(string(j)) } } }
		if v, ok := m["fxp"]; ok { if v == nil { data.Fxp = types.BoolNull() } else if b, ok := v.(bool); ok { data.Fxp = types.BoolValue(b) } }
		if v, ok := m["resume"]; ok { if v == nil { data.Resume = types.BoolNull() } else if b, ok := v.(bool); ok { data.Resume = types.BoolValue(b) } }
		if v, ok := m["defaultroot"]; ok { if v == nil { data.Defaultroot = types.BoolNull() } else if b, ok := v.(bool); ok { data.Defaultroot = types.BoolValue(b) } }
		if v, ok := m["ident"]; ok { if v == nil { data.Ident = types.BoolNull() } else if b, ok := v.(bool); ok { data.Ident = types.BoolValue(b) } }
		if v, ok := m["reversedns"]; ok { if v == nil { data.Reversedns = types.BoolNull() } else if b, ok := v.(bool); ok { data.Reversedns = types.BoolValue(b) } }
		if v, ok := m["masqaddress"]; ok { if v == nil { data.Masqaddress = types.StringNull() } else if s, ok := v.(string); ok { data.Masqaddress = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Masqaddress = types.StringValue(string(j)) } } }
		if v, ok := m["passiveportsmin"]; ok { if v == nil { data.Passiveportsmin = types.Int64Null() } else if f, ok := v.(float64); ok { data.Passiveportsmin = types.Int64Value(int64(f)) } }
		if v, ok := m["passiveportsmax"]; ok { if v == nil { data.Passiveportsmax = types.Int64Null() } else if f, ok := v.(float64); ok { data.Passiveportsmax = types.Int64Value(int64(f)) } }
		if v, ok := m["localuserbw"]; ok { if v == nil { data.Localuserbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Localuserbw = types.Int64Value(int64(f)) } }
		if v, ok := m["localuserdlbw"]; ok { if v == nil { data.Localuserdlbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Localuserdlbw = types.Int64Value(int64(f)) } }
		if v, ok := m["anonuserbw"]; ok { if v == nil { data.Anonuserbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Anonuserbw = types.Int64Value(int64(f)) } }
		if v, ok := m["anonuserdlbw"]; ok { if v == nil { data.Anonuserdlbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Anonuserdlbw = types.Int64Value(int64(f)) } }
		if v, ok := m["tls"]; ok { if v == nil { data.Tls = types.BoolNull() } else if b, ok := v.(bool); ok { data.Tls = types.BoolValue(b) } }
		if v, ok := m["tls_policy"]; ok { if v == nil { data.TlsPolicy = types.StringNull() } else if s, ok := v.(string); ok { data.TlsPolicy = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.TlsPolicy = types.StringValue(string(j)) } } }
		if v, ok := m["tls_opt_allow_client_renegotiations"]; ok { if v == nil { data.TlsOptAllowClientRenegotiations = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowClientRenegotiations = types.BoolValue(b) } }
		if v, ok := m["tls_opt_allow_dot_login"]; ok { if v == nil { data.TlsOptAllowDotLogin = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowDotLogin = types.BoolValue(b) } }
		if v, ok := m["tls_opt_allow_per_user"]; ok { if v == nil { data.TlsOptAllowPerUser = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowPerUser = types.BoolValue(b) } }
		if v, ok := m["tls_opt_common_name_required"]; ok { if v == nil { data.TlsOptCommonNameRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptCommonNameRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_enable_diags"]; ok { if v == nil { data.TlsOptEnableDiags = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptEnableDiags = types.BoolValue(b) } }
		if v, ok := m["tls_opt_export_cert_data"]; ok { if v == nil { data.TlsOptExportCertData = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptExportCertData = types.BoolValue(b) } }
		if v, ok := m["tls_opt_no_empty_fragments"]; ok { if v == nil { data.TlsOptNoEmptyFragments = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptNoEmptyFragments = types.BoolValue(b) } }
		if v, ok := m["tls_opt_no_session_reuse_required"]; ok { if v == nil { data.TlsOptNoSessionReuseRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptNoSessionReuseRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_stdenvvars"]; ok { if v == nil { data.TlsOptStdenvvars = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptStdenvvars = types.BoolValue(b) } }
		if v, ok := m["tls_opt_dns_name_required"]; ok { if v == nil { data.TlsOptDnsNameRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptDnsNameRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_ip_address_required"]; ok { if v == nil { data.TlsOptIpAddressRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptIpAddressRequired = types.BoolValue(b) } }
		if v, ok := m["ssltls_certificate"]; ok { if v == nil { data.SsltlsCertificate = types.Int64Null() } else if f, ok := v.(float64); ok { data.SsltlsCertificate = types.Int64Value(int64(f)) } }
		if v, ok := m["options"]; ok { if v == nil { data.Options = types.StringNull() } else if s, ok := v.(string); ok { data.Options = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Options = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FtpConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FtpConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	params := map[string]interface{}{}
	if !data.Port.IsNull() && !data.Port.IsUnknown() { params["port"] = data.Port.ValueInt64() }
	if !data.Clients.IsNull() && !data.Clients.IsUnknown() { params["clients"] = data.Clients.ValueInt64() }
	if !data.Ipconnections.IsNull() && !data.Ipconnections.IsUnknown() { params["ipconnections"] = data.Ipconnections.ValueInt64() }
	if !data.Loginattempt.IsNull() && !data.Loginattempt.IsUnknown() { params["loginattempt"] = data.Loginattempt.ValueInt64() }
	if !data.Timeout.IsNull() && !data.Timeout.IsUnknown() { params["timeout"] = data.Timeout.ValueInt64() }
	if !data.TimeoutNotransfer.IsNull() && !data.TimeoutNotransfer.IsUnknown() { params["timeout_notransfer"] = data.TimeoutNotransfer.ValueInt64() }
	if !data.Onlyanonymous.IsNull() && !data.Onlyanonymous.IsUnknown() { params["onlyanonymous"] = data.Onlyanonymous.ValueBool() }
	if !data.Anonpath.IsNull() && !data.Anonpath.IsUnknown() { params["anonpath"] = data.Anonpath.ValueString() }
	if !data.Onlylocal.IsNull() && !data.Onlylocal.IsUnknown() { params["onlylocal"] = data.Onlylocal.ValueBool() }
	if !data.Banner.IsNull() && !data.Banner.IsUnknown() { params["banner"] = data.Banner.ValueString() }
	if !data.Filemask.IsNull() && !data.Filemask.IsUnknown() { params["filemask"] = data.Filemask.ValueString() }
	if !data.Dirmask.IsNull() && !data.Dirmask.IsUnknown() { params["dirmask"] = data.Dirmask.ValueString() }
	if !data.Fxp.IsNull() && !data.Fxp.IsUnknown() { params["fxp"] = data.Fxp.ValueBool() }
	if !data.Resume.IsNull() && !data.Resume.IsUnknown() { params["resume"] = data.Resume.ValueBool() }
	if !data.Defaultroot.IsNull() && !data.Defaultroot.IsUnknown() { params["defaultroot"] = data.Defaultroot.ValueBool() }
	if !data.Ident.IsNull() && !data.Ident.IsUnknown() { params["ident"] = data.Ident.ValueBool() }
	if !data.Reversedns.IsNull() && !data.Reversedns.IsUnknown() { params["reversedns"] = data.Reversedns.ValueBool() }
	if !data.Masqaddress.IsNull() && !data.Masqaddress.IsUnknown() { params["masqaddress"] = data.Masqaddress.ValueString() }
	if !data.Passiveportsmin.IsNull() && !data.Passiveportsmin.IsUnknown() { params["passiveportsmin"] = data.Passiveportsmin.ValueInt64() }
	if !data.Passiveportsmax.IsNull() && !data.Passiveportsmax.IsUnknown() { params["passiveportsmax"] = data.Passiveportsmax.ValueInt64() }
	if !data.Localuserbw.IsNull() && !data.Localuserbw.IsUnknown() { params["localuserbw"] = data.Localuserbw.ValueInt64() }
	if !data.Localuserdlbw.IsNull() && !data.Localuserdlbw.IsUnknown() { params["localuserdlbw"] = data.Localuserdlbw.ValueInt64() }
	if !data.Anonuserbw.IsNull() && !data.Anonuserbw.IsUnknown() { params["anonuserbw"] = data.Anonuserbw.ValueInt64() }
	if !data.Anonuserdlbw.IsNull() && !data.Anonuserdlbw.IsUnknown() { params["anonuserdlbw"] = data.Anonuserdlbw.ValueInt64() }
	if !data.Tls.IsNull() && !data.Tls.IsUnknown() { params["tls"] = data.Tls.ValueBool() }
	if !data.TlsPolicy.IsNull() && !data.TlsPolicy.IsUnknown() { params["tls_policy"] = data.TlsPolicy.ValueString() }
	if !data.TlsOptAllowClientRenegotiations.IsNull() && !data.TlsOptAllowClientRenegotiations.IsUnknown() { params["tls_opt_allow_client_renegotiations"] = data.TlsOptAllowClientRenegotiations.ValueBool() }
	if !data.TlsOptAllowDotLogin.IsNull() && !data.TlsOptAllowDotLogin.IsUnknown() { params["tls_opt_allow_dot_login"] = data.TlsOptAllowDotLogin.ValueBool() }
	if !data.TlsOptAllowPerUser.IsNull() && !data.TlsOptAllowPerUser.IsUnknown() { params["tls_opt_allow_per_user"] = data.TlsOptAllowPerUser.ValueBool() }
	if !data.TlsOptCommonNameRequired.IsNull() && !data.TlsOptCommonNameRequired.IsUnknown() { params["tls_opt_common_name_required"] = data.TlsOptCommonNameRequired.ValueBool() }
	if !data.TlsOptEnableDiags.IsNull() && !data.TlsOptEnableDiags.IsUnknown() { params["tls_opt_enable_diags"] = data.TlsOptEnableDiags.ValueBool() }
	if !data.TlsOptExportCertData.IsNull() && !data.TlsOptExportCertData.IsUnknown() { params["tls_opt_export_cert_data"] = data.TlsOptExportCertData.ValueBool() }
	if !data.TlsOptNoEmptyFragments.IsNull() && !data.TlsOptNoEmptyFragments.IsUnknown() { params["tls_opt_no_empty_fragments"] = data.TlsOptNoEmptyFragments.ValueBool() }
	if !data.TlsOptNoSessionReuseRequired.IsNull() && !data.TlsOptNoSessionReuseRequired.IsUnknown() { params["tls_opt_no_session_reuse_required"] = data.TlsOptNoSessionReuseRequired.ValueBool() }
	if !data.TlsOptStdenvvars.IsNull() && !data.TlsOptStdenvvars.IsUnknown() { params["tls_opt_stdenvvars"] = data.TlsOptStdenvvars.ValueBool() }
	if !data.TlsOptDnsNameRequired.IsNull() && !data.TlsOptDnsNameRequired.IsUnknown() { params["tls_opt_dns_name_required"] = data.TlsOptDnsNameRequired.ValueBool() }
	if !data.TlsOptIpAddressRequired.IsNull() && !data.TlsOptIpAddressRequired.IsUnknown() { params["tls_opt_ip_address_required"] = data.TlsOptIpAddressRequired.ValueBool() }
	if !data.SsltlsCertificate.IsNull() && !data.SsltlsCertificate.IsUnknown() { params["ssltls_certificate"] = data.SsltlsCertificate.ValueInt64() }
	if !data.Options.IsNull() && !data.Options.IsUnknown() { params["options"] = data.Options.ValueString() }
	paramsArr := []interface{}{params}

	// Apply configuration
	_, err := r.client.Call("ftp.update", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to update ftp config: %s", err.Error()))
		return
	}

	// Read back to populate computed fields
	result, err := r.client.Call("ftp.config", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", fmt.Sprintf("Failed to read ftp config: %s", err.Error()))
		return
	}
	if m, ok := result.(map[string]interface{}); ok {
		if v, ok := m["port"]; ok { if v == nil { data.Port = types.Int64Null() } else if f, ok := v.(float64); ok { data.Port = types.Int64Value(int64(f)) } }
		if v, ok := m["clients"]; ok { if v == nil { data.Clients = types.Int64Null() } else if f, ok := v.(float64); ok { data.Clients = types.Int64Value(int64(f)) } }
		if v, ok := m["ipconnections"]; ok { if v == nil { data.Ipconnections = types.Int64Null() } else if f, ok := v.(float64); ok { data.Ipconnections = types.Int64Value(int64(f)) } }
		if v, ok := m["loginattempt"]; ok { if v == nil { data.Loginattempt = types.Int64Null() } else if f, ok := v.(float64); ok { data.Loginattempt = types.Int64Value(int64(f)) } }
		if v, ok := m["timeout"]; ok { if v == nil { data.Timeout = types.Int64Null() } else if f, ok := v.(float64); ok { data.Timeout = types.Int64Value(int64(f)) } }
		if v, ok := m["timeout_notransfer"]; ok { if v == nil { data.TimeoutNotransfer = types.Int64Null() } else if f, ok := v.(float64); ok { data.TimeoutNotransfer = types.Int64Value(int64(f)) } }
		if v, ok := m["onlyanonymous"]; ok { if v == nil { data.Onlyanonymous = types.BoolNull() } else if b, ok := v.(bool); ok { data.Onlyanonymous = types.BoolValue(b) } }
		if v, ok := m["anonpath"]; ok { if v == nil { data.Anonpath = types.StringNull() } else if s, ok := v.(string); ok { data.Anonpath = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Anonpath = types.StringValue(string(j)) } } }
		if v, ok := m["onlylocal"]; ok { if v == nil { data.Onlylocal = types.BoolNull() } else if b, ok := v.(bool); ok { data.Onlylocal = types.BoolValue(b) } }
		if v, ok := m["banner"]; ok { if v == nil { data.Banner = types.StringNull() } else if s, ok := v.(string); ok { data.Banner = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Banner = types.StringValue(string(j)) } } }
		if v, ok := m["filemask"]; ok { if v == nil { data.Filemask = types.StringNull() } else if s, ok := v.(string); ok { data.Filemask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Filemask = types.StringValue(string(j)) } } }
		if v, ok := m["dirmask"]; ok { if v == nil { data.Dirmask = types.StringNull() } else if s, ok := v.(string); ok { data.Dirmask = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Dirmask = types.StringValue(string(j)) } } }
		if v, ok := m["fxp"]; ok { if v == nil { data.Fxp = types.BoolNull() } else if b, ok := v.(bool); ok { data.Fxp = types.BoolValue(b) } }
		if v, ok := m["resume"]; ok { if v == nil { data.Resume = types.BoolNull() } else if b, ok := v.(bool); ok { data.Resume = types.BoolValue(b) } }
		if v, ok := m["defaultroot"]; ok { if v == nil { data.Defaultroot = types.BoolNull() } else if b, ok := v.(bool); ok { data.Defaultroot = types.BoolValue(b) } }
		if v, ok := m["ident"]; ok { if v == nil { data.Ident = types.BoolNull() } else if b, ok := v.(bool); ok { data.Ident = types.BoolValue(b) } }
		if v, ok := m["reversedns"]; ok { if v == nil { data.Reversedns = types.BoolNull() } else if b, ok := v.(bool); ok { data.Reversedns = types.BoolValue(b) } }
		if v, ok := m["masqaddress"]; ok { if v == nil { data.Masqaddress = types.StringNull() } else if s, ok := v.(string); ok { data.Masqaddress = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Masqaddress = types.StringValue(string(j)) } } }
		if v, ok := m["passiveportsmin"]; ok { if v == nil { data.Passiveportsmin = types.Int64Null() } else if f, ok := v.(float64); ok { data.Passiveportsmin = types.Int64Value(int64(f)) } }
		if v, ok := m["passiveportsmax"]; ok { if v == nil { data.Passiveportsmax = types.Int64Null() } else if f, ok := v.(float64); ok { data.Passiveportsmax = types.Int64Value(int64(f)) } }
		if v, ok := m["localuserbw"]; ok { if v == nil { data.Localuserbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Localuserbw = types.Int64Value(int64(f)) } }
		if v, ok := m["localuserdlbw"]; ok { if v == nil { data.Localuserdlbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Localuserdlbw = types.Int64Value(int64(f)) } }
		if v, ok := m["anonuserbw"]; ok { if v == nil { data.Anonuserbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Anonuserbw = types.Int64Value(int64(f)) } }
		if v, ok := m["anonuserdlbw"]; ok { if v == nil { data.Anonuserdlbw = types.Int64Null() } else if f, ok := v.(float64); ok { data.Anonuserdlbw = types.Int64Value(int64(f)) } }
		if v, ok := m["tls"]; ok { if v == nil { data.Tls = types.BoolNull() } else if b, ok := v.(bool); ok { data.Tls = types.BoolValue(b) } }
		if v, ok := m["tls_policy"]; ok { if v == nil { data.TlsPolicy = types.StringNull() } else if s, ok := v.(string); ok { data.TlsPolicy = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.TlsPolicy = types.StringValue(string(j)) } } }
		if v, ok := m["tls_opt_allow_client_renegotiations"]; ok { if v == nil { data.TlsOptAllowClientRenegotiations = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowClientRenegotiations = types.BoolValue(b) } }
		if v, ok := m["tls_opt_allow_dot_login"]; ok { if v == nil { data.TlsOptAllowDotLogin = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowDotLogin = types.BoolValue(b) } }
		if v, ok := m["tls_opt_allow_per_user"]; ok { if v == nil { data.TlsOptAllowPerUser = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptAllowPerUser = types.BoolValue(b) } }
		if v, ok := m["tls_opt_common_name_required"]; ok { if v == nil { data.TlsOptCommonNameRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptCommonNameRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_enable_diags"]; ok { if v == nil { data.TlsOptEnableDiags = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptEnableDiags = types.BoolValue(b) } }
		if v, ok := m["tls_opt_export_cert_data"]; ok { if v == nil { data.TlsOptExportCertData = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptExportCertData = types.BoolValue(b) } }
		if v, ok := m["tls_opt_no_empty_fragments"]; ok { if v == nil { data.TlsOptNoEmptyFragments = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptNoEmptyFragments = types.BoolValue(b) } }
		if v, ok := m["tls_opt_no_session_reuse_required"]; ok { if v == nil { data.TlsOptNoSessionReuseRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptNoSessionReuseRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_stdenvvars"]; ok { if v == nil { data.TlsOptStdenvvars = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptStdenvvars = types.BoolValue(b) } }
		if v, ok := m["tls_opt_dns_name_required"]; ok { if v == nil { data.TlsOptDnsNameRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptDnsNameRequired = types.BoolValue(b) } }
		if v, ok := m["tls_opt_ip_address_required"]; ok { if v == nil { data.TlsOptIpAddressRequired = types.BoolNull() } else if b, ok := v.(bool); ok { data.TlsOptIpAddressRequired = types.BoolValue(b) } }
		if v, ok := m["ssltls_certificate"]; ok { if v == nil { data.SsltlsCertificate = types.Int64Null() } else if f, ok := v.(float64); ok { data.SsltlsCertificate = types.Int64Value(int64(f)) } }
		if v, ok := m["options"]; ok { if v == nil { data.Options = types.StringNull() } else if s, ok := v.(string); ok { data.Options = types.StringValue(s) } else { if j, e := json.Marshal(v); e == nil { data.Options = types.StringValue(string(j)) } } }
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FtpConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton configs cannot be deleted - just remove from state
}
