package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type CertificateResource struct {
	client *client.Client
}

type CertificateResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	CreateType types.String `tfsdk:"create_type"`
	AddToTrustedStore types.Bool `tfsdk:"add_to_trusted_store"`
	Certificate types.String `tfsdk:"certificate"`
	Privatekey types.String `tfsdk:"privatekey"`
	Csr types.String `tfsdk:"csr"`
	KeyLength types.Int64 `tfsdk:"key_length"`
	KeyType types.String `tfsdk:"key_type"`
	EcCurve types.String `tfsdk:"ec_curve"`
	Passphrase types.String `tfsdk:"passphrase"`
	City types.String `tfsdk:"city"`
	Common types.String `tfsdk:"common"`
	Country types.String `tfsdk:"country"`
	Email types.String `tfsdk:"email"`
	Organization types.String `tfsdk:"organization"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
	State types.String `tfsdk:"state"`
	DigestAlgorithm types.String `tfsdk:"digest_algorithm"`
	San types.List `tfsdk:"san"`
	CertExtensions types.String `tfsdk:"cert_extensions"`
	AcmeDirectoryUri types.String `tfsdk:"acme_directory_uri"`
	CsrId types.Int64 `tfsdk:"csr_id"`
	Tos types.String `tfsdk:"tos"`
	DnsMapping types.String `tfsdk:"dns_mapping"`
	RenewDays types.Int64 `tfsdk:"renew_days"`
}

func NewCertificateResource() resource.Resource {
	return &CertificateResource{}
}

func (r *CertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *CertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new Certificate",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Certificate name.",
			},
			"create_type": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Type of certificate creation operation.",
			},
			"add_to_trusted_store": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to add this certificate to the trusted certificate store.",
			},
			"certificate": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "PEM-encoded certificate to import or `null`.",
			},
			"privatekey": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "PEM-encoded private key to import or `null`.",
			},
			"csr": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "PEM-encoded certificate signing request to import or `null`.",
			},
			"key_length": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "RSA key length in bits or `null`.",
			},
			"key_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Type of cryptographic key to generate.",
			},
			"ec_curve": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Elliptic curve to use for EC keys.",
			},
			"passphrase": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Passphrase to protect the private key or `null`.",
			},
			"city": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "City or locality name for certificate subject or `null`.",
			},
			"common": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Common name for certificate subject or `null`.",
			},
			"country": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Country name for certificate subject or `null`.",
			},
			"email": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Email address for certificate subject or `null`.",
			},
			"organization": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Organization name for certificate subject or `null`.",
			},
			"organizational_unit": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Organizational unit for certificate subject or `null`.",
			},
			"state": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "State or province name for certificate subject or `null`.",
			},
			"digest_algorithm": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Hash algorithm for certificate signing.",
			},
			"san": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Subject alternative names for the certificate.",
			},
			"cert_extensions": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Certificate extensions configuration.",
			},
			"acme_directory_uri": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "ACME directory URI to be used for ACME certificate creation.",
			},
			"csr_id": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "CSR to be used for ACME certificate creation.",
			},
			"tos": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Set this when creating an ACME certificate to accept terms of service of the ACME service.",
			},
			"dns_mapping": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "A mapping of domain to ACME DNS Authenticator ID for each domain listed in SAN or common name of the",
			},
			"renew_days": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of days before the certificate expiration date to attempt certificate renewal. If certificate",
			},
		},
	}
}

func (r *CertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.CreateType.IsNull() {
		params["create_type"] = data.CreateType.ValueString()
	}
	if !data.AddToTrustedStore.IsNull() {
		params["add_to_trusted_store"] = data.AddToTrustedStore.ValueBool()
	}
	if !data.Certificate.IsNull() {
		params["certificate"] = data.Certificate.ValueString()
	}
	if !data.Privatekey.IsNull() {
		params["privatekey"] = data.Privatekey.ValueString()
	}
	if !data.Csr.IsNull() {
		params["CSR"] = data.Csr.ValueString()
	}
	if !data.KeyLength.IsNull() {
		params["key_length"] = data.KeyLength.ValueInt64()
	}
	if !data.KeyType.IsNull() {
		params["key_type"] = data.KeyType.ValueString()
	}
	if !data.EcCurve.IsNull() {
		params["ec_curve"] = data.EcCurve.ValueString()
	}
	if !data.Passphrase.IsNull() {
		params["passphrase"] = data.Passphrase.ValueString()
	}
	if !data.City.IsNull() {
		params["city"] = data.City.ValueString()
	}
	if !data.Common.IsNull() {
		params["common"] = data.Common.ValueString()
	}
	if !data.Country.IsNull() {
		params["country"] = data.Country.ValueString()
	}
	if !data.Email.IsNull() {
		params["email"] = data.Email.ValueString()
	}
	if !data.Organization.IsNull() {
		params["organization"] = data.Organization.ValueString()
	}
	if !data.OrganizationalUnit.IsNull() {
		params["organizational_unit"] = data.OrganizationalUnit.ValueString()
	}
	if !data.State.IsNull() {
		params["state"] = data.State.ValueString()
	}
	if !data.DigestAlgorithm.IsNull() {
		params["digest_algorithm"] = data.DigestAlgorithm.ValueString()
	}
	if !data.San.IsNull() {
		var sanList []string
		data.San.ElementsAs(ctx, &sanList, false)
		params["san"] = sanList
	}
	if !data.CertExtensions.IsNull() {
		params["cert_extensions"] = data.CertExtensions.ValueString()
	}
	if !data.AcmeDirectoryUri.IsNull() {
		params["acme_directory_uri"] = data.AcmeDirectoryUri.ValueString()
	}
	if !data.CsrId.IsNull() {
		params["csr_id"] = data.CsrId.ValueInt64()
	}
	if !data.Tos.IsNull() {
		params["tos"] = data.Tos.ValueString()
	}
	if !data.DnsMapping.IsNull() {
		params["dns_mapping"] = data.DnsMapping.ValueString()
	}
	if !data.RenewDays.IsNull() {
		params["renew_days"] = data.RenewDays.ValueInt64()
	}

	result, err := r.client.CallWithJob("certificate.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create certificate: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("certificate.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read certificate: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["create_type"]; ok && v != nil {
			data.CreateType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["add_to_trusted_store"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AddToTrustedStore = types.BoolValue(bv) }
		}
		if v, ok := resultMap["certificate"]; ok && v != nil {
			data.Certificate = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["privatekey"]; ok && v != nil {
			data.Privatekey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["CSR"]; ok && v != nil {
			data.Csr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["key_length"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.KeyLength = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["key_type"]; ok && v != nil {
			data.KeyType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ec_curve"]; ok && v != nil {
			data.EcCurve = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["passphrase"]; ok && v != nil {
			data.Passphrase = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["city"]; ok && v != nil {
			data.City = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["common"]; ok && v != nil {
			data.Common = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["country"]; ok && v != nil {
			data.Country = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["email"]; ok && v != nil {
			data.Email = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["organization"]; ok && v != nil {
			data.Organization = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["organizational_unit"]; ok && v != nil {
			data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			data.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["digest_algorithm"]; ok && v != nil {
			data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["san"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.San, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["cert_extensions"]; ok && v != nil {
			data.CertExtensions = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["acme_directory_uri"]; ok && v != nil {
			data.AcmeDirectoryUri = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["csr_id"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.CsrId = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["tos"]; ok && v != nil {
			data.Tos = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dns_mapping"]; ok && v != nil {
			data.DnsMapping = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["renew_days"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.RenewDays = types.Int64Value(int64(fv)) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.CreateType.IsNull() {
		params["create_type"] = data.CreateType.ValueString()
	}
	if !data.AddToTrustedStore.IsNull() {
		params["add_to_trusted_store"] = data.AddToTrustedStore.ValueBool()
	}
	if !data.Certificate.IsNull() {
		params["certificate"] = data.Certificate.ValueString()
	}
	if !data.Privatekey.IsNull() {
		params["privatekey"] = data.Privatekey.ValueString()
	}
	if !data.Csr.IsNull() {
		params["CSR"] = data.Csr.ValueString()
	}
	if !data.KeyLength.IsNull() {
		params["key_length"] = data.KeyLength.ValueInt64()
	}
	if !data.KeyType.IsNull() {
		params["key_type"] = data.KeyType.ValueString()
	}
	if !data.EcCurve.IsNull() {
		params["ec_curve"] = data.EcCurve.ValueString()
	}
	if !data.Passphrase.IsNull() {
		params["passphrase"] = data.Passphrase.ValueString()
	}
	if !data.City.IsNull() {
		params["city"] = data.City.ValueString()
	}
	if !data.Common.IsNull() {
		params["common"] = data.Common.ValueString()
	}
	if !data.Country.IsNull() {
		params["country"] = data.Country.ValueString()
	}
	if !data.Email.IsNull() {
		params["email"] = data.Email.ValueString()
	}
	if !data.Organization.IsNull() {
		params["organization"] = data.Organization.ValueString()
	}
	if !data.OrganizationalUnit.IsNull() {
		params["organizational_unit"] = data.OrganizationalUnit.ValueString()
	}
	if !data.State.IsNull() {
		params["state"] = data.State.ValueString()
	}
	if !data.DigestAlgorithm.IsNull() {
		params["digest_algorithm"] = data.DigestAlgorithm.ValueString()
	}
	if !data.San.IsNull() {
		var sanList []string
		data.San.ElementsAs(ctx, &sanList, false)
		params["san"] = sanList
	}
	if !data.CertExtensions.IsNull() {
		params["cert_extensions"] = data.CertExtensions.ValueString()
	}
	if !data.AcmeDirectoryUri.IsNull() {
		params["acme_directory_uri"] = data.AcmeDirectoryUri.ValueString()
	}
	if !data.CsrId.IsNull() {
		params["csr_id"] = data.CsrId.ValueInt64()
	}
	if !data.Tos.IsNull() {
		params["tos"] = data.Tos.ValueString()
	}
	if !data.DnsMapping.IsNull() {
		params["dns_mapping"] = data.DnsMapping.ValueString()
	}
	if !data.RenewDays.IsNull() {
		params["renew_days"] = data.RenewDays.ValueInt64()
	}

	_, err = r.client.CallWithJob("certificate.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update certificate: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.CallWithJob("certificate.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete certificate: %s", err))
		return
	}
}
