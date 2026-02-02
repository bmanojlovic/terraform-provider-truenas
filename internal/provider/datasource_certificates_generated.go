package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &CertificatesDataSource{}

func NewCertificatesDataSource() datasource.DataSource {
	return &CertificatesDataSource{}
}

type CertificatesDataSource struct {
	client *client.Client
}

type CertificatesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type CertificatesItemModel struct {
	ID                    types.String `tfsdk:"id"`
	Type                  types.Int64  `tfsdk:"type"`
	Name                  types.String `tfsdk:"name"`
	Certificate           types.String `tfsdk:"certificate"`
	Privatekey            types.String `tfsdk:"privatekey"`
	Csr                   types.String `tfsdk:"csr"`
	AcmeUri               types.String `tfsdk:"acme_uri"`
	DomainsAuthenticators types.String `tfsdk:"domains_authenticators"`
	RenewDays             types.Int64  `tfsdk:"renew_days"`
	Acme                  types.String `tfsdk:"acme"`
	AddToTrustedStore     types.Bool   `tfsdk:"add_to_trusted_store"`
	RootPath              types.String `tfsdk:"root_path"`
	CertificatePath       types.String `tfsdk:"certificate_path"`
	PrivatekeyPath        types.String `tfsdk:"privatekey_path"`
	CsrPath               types.String `tfsdk:"csr_path"`
	CertType              types.String `tfsdk:"cert_type"`
	CertTypeExisting      types.Bool   `tfsdk:"cert_type_existing"`
	CertTypeCsr           types.Bool   `tfsdk:"cert_type_csr"`
	CertTypeCa            types.Bool   `tfsdk:"cert_type_ca"`
	KeyLength             types.Int64  `tfsdk:"key_length"`
	KeyType               types.String `tfsdk:"key_type"`
	Country               types.String `tfsdk:"country"`
	State                 types.String `tfsdk:"state"`
	City                  types.String `tfsdk:"city"`
	Organization          types.String `tfsdk:"organization"`
	OrganizationalUnit    types.String `tfsdk:"organizational_unit"`
	Common                types.String `tfsdk:"common"`
	Email                 types.String `tfsdk:"email"`
	Dn                    types.String `tfsdk:"dn"`
	SubjectNameHash       types.Int64  `tfsdk:"subject_name_hash"`
	Extensions            types.String `tfsdk:"extensions"`
	DigestAlgorithm       types.String `tfsdk:"digest_algorithm"`
	Lifetime              types.Int64  `tfsdk:"lifetime"`
	From                  types.String `tfsdk:"from"`
	Until                 types.String `tfsdk:"until"`
	Serial                types.Int64  `tfsdk:"serial"`
	Chain                 types.Bool   `tfsdk:"chain"`
	Fingerprint           types.String `tfsdk:"fingerprint"`
	Expired               types.Bool   `tfsdk:"expired"`
	Parsed                types.Bool   `tfsdk:"parsed"`
}

func (d *CertificatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

func (d *CertificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query certificates",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of certificates resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"type": schema.Int64Attribute{
							Computed:    true,
							Description: "Internal certificate type identifier used to determine certificate capabilities.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable name for this certificate. Must be unique and contain only alphanumeric characters,  ",
						},
						"certificate": schema.StringAttribute{
							Computed:    true,
							Description: "PEM-encoded X.509 certificate data. `null` for certificate signing requests (CSR) that have not yet ",
						},
						"privatekey": schema.StringAttribute{
							Computed:    true,
							Description: "PEM-encoded private key corresponding to the certificate. `null` if no private key is available or f",
						},
						"csr": schema.StringAttribute{
							Computed:    true,
							Description: "PEM-encoded Certificate Signing Request (CSR) data. `null` for imported certificates or completed   ",
						},
						"acme_uri": schema.StringAttribute{
							Computed:    true,
							Description: "ACME directory server URI used for automated certificate management. `null` for non-ACME certificate",
						},
						"domains_authenticators": schema.StringAttribute{
							Computed:    true,
							Description: "Mapping of domain names to ACME DNS authenticator IDs for domain validation. `null` for non-ACME    ",
						},
						"renew_days": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of days before expiration to attempt automatic renewal. Only applicable for ACME certificates",
						},
						"acme": schema.StringAttribute{
							Computed:    true,
							Description: "ACME registration and account information used for certificate lifecycle management. `null` for     ",
						},
						"add_to_trusted_store": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this certificate should be added to the system's trusted certificate store.",
						},
						"root_path": schema.StringAttribute{
							Computed:    true,
							Description: "Filesystem path where certificate-related files are stored.",
						},
						"certificate_path": schema.StringAttribute{
							Computed:    true,
							Description: "Filesystem path to the certificate file (.crt). `null` if no certificate is available.",
						},
						"privatekey_path": schema.StringAttribute{
							Computed:    true,
							Description: "Filesystem path to the private key file (.key). `null` if no private key is available.",
						},
						"csr_path": schema.StringAttribute{
							Computed:    true,
							Description: "Filesystem path to the certificate signing request file (.csr). `null` if no CSR is available.",
						},
						"cert_type": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable certificate type, typically 'CERTIFICATE' for standard certificates.",
						},
						"cert_type_existing": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this is an existing certificate (imported or generated).",
						},
						"cert_type_csr": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this entry represents a Certificate Signing Request (CSR) rather than a signed certificate.",
						},
						"cert_type_ca": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this certificate is a Certificate Authority (CA) certificate.",
						},
						"key_length": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the cryptographic key in bits. `null` if key information is unavailable.",
						},
						"key_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of cryptographic key algorithm (e.g., 'RSA', 'EC', 'DSA'). `null` if key information is unavail",
						},
						"country": schema.StringAttribute{
							Computed:    true,
							Description: "ISO 3166-1 alpha-2 country code from the certificate subject. `null` if not specified.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "State or province name from the certificate subject. `null` if not specified.",
						},
						"city": schema.StringAttribute{
							Computed:    true,
							Description: "City or locality name from the certificate subject. `null` if not specified.",
						},
						"organization": schema.StringAttribute{
							Computed:    true,
							Description: "Organization name from the certificate subject. `null` if not specified.",
						},
						"organizational_unit": schema.StringAttribute{
							Computed:    true,
							Description: "Organizational unit from the certificate subject. `null` if not specified.",
						},
						"common": schema.StringAttribute{
							Computed:    true,
							Description: "Common name (CN) from the certificate subject. `null` if not specified.",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "Email address from the certificate subject. `null` if not specified.",
						},
						"dn": schema.StringAttribute{
							Computed:    true,
							Description: "Distinguished Name (DN) of the certificate subject in RFC 2253 format. `null` if certificate parsing",
						},
						"subject_name_hash": schema.Int64Attribute{
							Computed:    true,
							Description: "Hash of the certificate subject name. `null` if certificate parsing failed.",
						},
						"extensions": schema.StringAttribute{
							Computed:    true,
							Description: "X.509 certificate extensions parsed into a dictionary structure.",
						},
						"digest_algorithm": schema.StringAttribute{
							Computed:    true,
							Description: "Cryptographic hash algorithm used for certificate signing (e.g., 'SHA256'). `null` if unavailable.",
						},
						"lifetime": schema.Int64Attribute{
							Computed:    true,
							Description: "Certificate validity period in seconds. `null` if certificate parsing failed.",
						},
						"from": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate validity start date in ISO 8601 format. `null` if certificate parsing failed.",
						},
						"until": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate validity end date in ISO 8601 format. `null` if certificate parsing failed.",
						},
						"serial": schema.Int64Attribute{
							Computed:    true,
							Description: "Certificate serial number. `null` if certificate parsing failed.",
						},
						"chain": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this certificate has an associated certificate chain. `null` if unavailable.",
						},
						"fingerprint": schema.StringAttribute{
							Computed:    true,
							Description: "SHA-256 fingerprint of the certificate in hexadecimal format. `null` if certificate parsing failed.",
						},
						"expired": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the certificate has expired. `null` if certificate parsing failed.",
						},
						"parsed": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the certificate data was successfully parsed and validated.",
						},
					},
				},
			},
		},
	}
}

func (d *CertificatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *CertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertificatesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("certificate.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query certificates: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]CertificatesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := CertificatesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Type = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["certificate"]; ok && v != nil {
			itemModel.Certificate = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["privatekey"]; ok && v != nil {
			itemModel.Privatekey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["CSR"]; ok && v != nil {
			itemModel.Csr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["acme_uri"]; ok && v != nil {
			itemModel.AcmeUri = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["domains_authenticators"]; ok && v != nil {
			itemModel.DomainsAuthenticators = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["renew_days"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.RenewDays = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["acme"]; ok && v != nil {
			itemModel.Acme = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["add_to_trusted_store"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.AddToTrustedStore = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["root_path"]; ok && v != nil {
			itemModel.RootPath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["certificate_path"]; ok && v != nil {
			itemModel.CertificatePath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["privatekey_path"]; ok && v != nil {
			itemModel.PrivatekeyPath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["csr_path"]; ok && v != nil {
			itemModel.CsrPath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cert_type"]; ok && v != nil {
			itemModel.CertType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cert_type_existing"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.CertTypeExisting = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["cert_type_CSR"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.CertTypeCsr = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["cert_type_CA"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.CertTypeCa = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["key_length"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.KeyLength = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["key_type"]; ok && v != nil {
			itemModel.KeyType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["country"]; ok && v != nil {
			itemModel.Country = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["city"]; ok && v != nil {
			itemModel.City = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["organization"]; ok && v != nil {
			itemModel.Organization = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["organizational_unit"]; ok && v != nil {
			itemModel.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["common"]; ok && v != nil {
			itemModel.Common = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["email"]; ok && v != nil {
			itemModel.Email = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["DN"]; ok && v != nil {
			itemModel.Dn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["subject_name_hash"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.SubjectNameHash = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["extensions"]; ok && v != nil {
			itemModel.Extensions = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["digest_algorithm"]; ok && v != nil {
			itemModel.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lifetime"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Lifetime = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["from"]; ok && v != nil {
			itemModel.From = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["until"]; ok && v != nil {
			itemModel.Until = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["serial"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Serial = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["chain"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Chain = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["fingerprint"]; ok && v != nil {
			itemModel.Fingerprint = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["expired"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Expired = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["parsed"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Parsed = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"csr":                    types.StringType,
			"dn":                     types.StringType,
			"acme":                   types.StringType,
			"acme_uri":               types.StringType,
			"add_to_trusted_store":   types.BoolType,
			"cert_type":              types.StringType,
			"cert_type_ca":           types.BoolType,
			"cert_type_csr":          types.BoolType,
			"cert_type_existing":     types.BoolType,
			"certificate":            types.StringType,
			"certificate_path":       types.StringType,
			"chain":                  types.BoolType,
			"city":                   types.StringType,
			"common":                 types.StringType,
			"country":                types.StringType,
			"csr_path":               types.StringType,
			"digest_algorithm":       types.StringType,
			"domains_authenticators": types.StringType,
			"email":                  types.StringType,
			"expired":                types.BoolType,
			"extensions":             types.StringType,
			"fingerprint":            types.StringType,
			"from":                   types.StringType,
			"id":                     types.StringType,
			"key_length":             types.Int64Type,
			"key_type":               types.StringType,
			"lifetime":               types.Int64Type,
			"name":                   types.StringType,
			"organization":           types.StringType,
			"organizational_unit":    types.StringType,
			"parsed":                 types.BoolType,
			"privatekey":             types.StringType,
			"privatekey_path":        types.StringType,
			"renew_days":             types.Int64Type,
			"root_path":              types.StringType,
			"serial":                 types.Int64Type,
			"state":                  types.StringType,
			"subject_name_hash":      types.Int64Type,
			"type":                   types.Int64Type,
			"until":                  types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
