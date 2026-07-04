package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type CertificateResource struct {
	client *client.Client
}

type CertificateResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	CreateType         types.String `tfsdk:"create_type"`
	AddToTrustedStore  types.Bool   `tfsdk:"add_to_trusted_store"`
	Certificate        types.String `tfsdk:"certificate"`
	Privatekey         types.String `tfsdk:"privatekey"`
	Csr                types.String `tfsdk:"csr"`
	KeyLength          types.Int64  `tfsdk:"key_length"`
	KeyType            types.String `tfsdk:"key_type"`
	EcCurve            types.String `tfsdk:"ec_curve"`
	Passphrase         types.String `tfsdk:"passphrase"`
	City               types.String `tfsdk:"city"`
	Common             types.String `tfsdk:"common"`
	Country            types.String `tfsdk:"country"`
	Email              types.String `tfsdk:"email"`
	Organization       types.String `tfsdk:"organization"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
	State              types.String `tfsdk:"state"`
	DigestAlgorithm    types.String `tfsdk:"digest_algorithm"`
	San                types.List   `tfsdk:"san"`
	CertExtensions     types.String `tfsdk:"cert_extensions"`
	AcmeDirectoryUri   types.String `tfsdk:"acme_directory_uri"`
	CsrId              types.Int64  `tfsdk:"csr_id"`
	Tos                types.Bool   `tfsdk:"tos"`
	DnsMapping         types.String `tfsdk:"dns_mapping"`
	RenewDays          types.Int64  `tfsdk:"renew_days"`
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
				Required:    true,
				Optional:    false,
				Description: "Certificate name.",
			},
			"create_type": schema.StringAttribute{
				Required:      true,
				Optional:      false,
				Description:   "Type of certificate creation operation.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"add_to_trusted_store": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to add this certificate to the trusted certificate store.",
			},
			"certificate": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "PEM-encoded certificate to import or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"privatekey": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "PEM-encoded private key to import or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"csr": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "PEM-encoded certificate signing request to import or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_length": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Description:   "RSA key length in bits or `null`.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"key_type": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Type of cryptographic key to generate.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ec_curve": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Elliptic curve to use for EC keys.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"passphrase": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Passphrase to protect the private key or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"city": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "City or locality name for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"common": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Common name for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"country": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Country name for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"email": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Email address for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"organization": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Organization name for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"organizational_unit": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Organizational unit for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"state": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "State or province name for certificate subject or `null`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"digest_algorithm": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Hash algorithm for certificate signing.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"san": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Subject alternative names for the certificate.",
			},
			"cert_extensions": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Certificate extensions configuration.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"acme_directory_uri": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "ACME directory URI to be used for ACME certificate creation.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"csr_id": schema.Int64Attribute{
				Required:      false,
				Optional:      true,
				Description:   "CSR to be used for ACME certificate creation.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"tos": schema.BoolAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Set this when creating an ACME certificate to accept terms of service of the ACME service.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"dns_mapping": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "A mapping of domain to ACME DNS Authenticator ID for each domain listed in SAN or common name of the",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"renew_days": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Days before expiration to attempt renewal.",
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
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.CreateType.IsNull() && !data.CreateType.IsUnknown() {
		params["create_type"] = data.CreateType.ValueString()
	}
	if !data.AddToTrustedStore.IsNull() && !data.AddToTrustedStore.IsUnknown() {
		params["add_to_trusted_store"] = data.AddToTrustedStore.ValueBool()
	}
	if !data.Certificate.IsNull() && !data.Certificate.IsUnknown() {
		params["certificate"] = data.Certificate.ValueString()
	}
	if !data.Privatekey.IsNull() && !data.Privatekey.IsUnknown() {
		params["privatekey"] = data.Privatekey.ValueString()
	}
	if !data.Csr.IsNull() && !data.Csr.IsUnknown() {
		params["CSR"] = data.Csr.ValueString()
	}
	if !data.KeyLength.IsNull() && !data.KeyLength.IsUnknown() {
		params["key_length"] = data.KeyLength.ValueInt64()
	}
	if !data.KeyType.IsNull() && !data.KeyType.IsUnknown() {
		params["key_type"] = data.KeyType.ValueString()
	}
	if !data.EcCurve.IsNull() && !data.EcCurve.IsUnknown() {
		params["ec_curve"] = data.EcCurve.ValueString()
	}
	if !data.Passphrase.IsNull() && !data.Passphrase.IsUnknown() {
		params["passphrase"] = data.Passphrase.ValueString()
	}
	if !data.City.IsNull() && !data.City.IsUnknown() {
		params["city"] = data.City.ValueString()
	}
	if !data.Common.IsNull() && !data.Common.IsUnknown() {
		params["common"] = data.Common.ValueString()
	}
	if !data.Country.IsNull() && !data.Country.IsUnknown() {
		params["country"] = data.Country.ValueString()
	}
	if !data.Email.IsNull() && !data.Email.IsUnknown() {
		params["email"] = data.Email.ValueString()
	}
	if !data.Organization.IsNull() && !data.Organization.IsUnknown() {
		params["organization"] = data.Organization.ValueString()
	}
	if !data.OrganizationalUnit.IsNull() && !data.OrganizationalUnit.IsUnknown() {
		params["organizational_unit"] = data.OrganizationalUnit.ValueString()
	}
	if !data.State.IsNull() && !data.State.IsUnknown() {
		params["state"] = data.State.ValueString()
	}
	if !data.DigestAlgorithm.IsNull() && !data.DigestAlgorithm.IsUnknown() {
		params["digest_algorithm"] = data.DigestAlgorithm.ValueString()
	}
	if !data.San.IsNull() && !data.San.IsUnknown() {
		var sanList []string
		data.San.ElementsAs(ctx, &sanList, false)
		params["san"] = sanList
	}
	if !data.CertExtensions.IsNull() && !data.CertExtensions.IsUnknown() {
		var cert_extensionsObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.CertExtensions.ValueString()), &cert_extensionsObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse cert_extensions: %s", err))
			return
		}
		params["cert_extensions"] = cert_extensionsObj
	}
	if !data.AcmeDirectoryUri.IsNull() && !data.AcmeDirectoryUri.IsUnknown() {
		params["acme_directory_uri"] = data.AcmeDirectoryUri.ValueString()
	}
	if !data.CsrId.IsNull() && !data.CsrId.IsUnknown() {
		params["csr_id"] = data.CsrId.ValueInt64()
	}
	if !data.Tos.IsNull() && !data.Tos.IsUnknown() {
		params["tos"] = data.Tos.ValueBool()
	}
	if !data.DnsMapping.IsNull() && !data.DnsMapping.IsUnknown() {
		var dns_mappingObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.DnsMapping.ValueString()), &dns_mappingObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse dns_mapping: %s", err))
			return
		}
		params["dns_mapping"] = dns_mappingObj
	}
	if !data.RenewDays.IsNull() && !data.RenewDays.IsUnknown() {
		params["renew_days"] = data.RenewDays.ValueInt64()
	}

	result, err := r.client.CallWithJob("certificate.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create certificate: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	// Read back to populate computed fields
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("certificate.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back certificate: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["certificate"]; ok {
		if v == nil {
			data.Certificate = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Certificate = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Certificate = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Certificate = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["privatekey"]; ok {
		if v == nil {
			data.Privatekey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Privatekey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Privatekey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Privatekey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["CSR"]; ok {
		if v == nil {
			data.Csr = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Csr = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Csr = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Csr = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["key_length"]; ok {
		if v == nil {
			data.KeyLength = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.KeyLength = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.KeyLength = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["key_type"]; ok {
		switch val := v.(type) {
		case string:
			data.KeyType = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.KeyType = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.KeyType = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["city"]; ok {
		if v == nil {
			data.City = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.City = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.City = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.City = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["common"]; ok {
		if v == nil {
			data.Common = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Common = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Common = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Common = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["country"]; ok {
		if v == nil {
			data.Country = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Country = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Country = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Country = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["email"]; ok {
		if v == nil {
			data.Email = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Email = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Email = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Email = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["organization"]; ok {
		if v == nil {
			data.Organization = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Organization = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Organization = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Organization = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["organizational_unit"]; ok {
		if v == nil {
			data.OrganizationalUnit = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.OrganizationalUnit = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["state"]; ok {
		if v == nil {
			data.State = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.State = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.State = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.State = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["digest_algorithm"]; ok {
		switch val := v.(type) {
		case string:
			data.DigestAlgorithm = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["san"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.San, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if data.AddToTrustedStore.IsUnknown() {
		data.AddToTrustedStore = types.BoolNull()
	}
	if data.Certificate.IsUnknown() {
		data.Certificate = types.StringNull()
	}
	if data.Privatekey.IsUnknown() {
		data.Privatekey = types.StringNull()
	}
	if data.Csr.IsUnknown() {
		data.Csr = types.StringNull()
	}
	if data.KeyLength.IsUnknown() {
		data.KeyLength = types.Int64Null()
	}
	if data.KeyType.IsUnknown() {
		data.KeyType = types.StringNull()
	}
	if data.EcCurve.IsUnknown() {
		data.EcCurve = types.StringNull()
	}
	if data.Passphrase.IsUnknown() {
		data.Passphrase = types.StringNull()
	}
	if data.City.IsUnknown() {
		data.City = types.StringNull()
	}
	if data.Common.IsUnknown() {
		data.Common = types.StringNull()
	}
	if data.Country.IsUnknown() {
		data.Country = types.StringNull()
	}
	if data.Email.IsUnknown() {
		data.Email = types.StringNull()
	}
	if data.Organization.IsUnknown() {
		data.Organization = types.StringNull()
	}
	if data.OrganizationalUnit.IsUnknown() {
		data.OrganizationalUnit = types.StringNull()
	}
	if data.State.IsUnknown() {
		data.State = types.StringNull()
	}
	if data.DigestAlgorithm.IsUnknown() {
		data.DigestAlgorithm = types.StringNull()
	}
	if data.San.IsUnknown() {
		data.San, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.CertExtensions.IsUnknown() {
		data.CertExtensions = types.StringNull()
	}
	if data.AcmeDirectoryUri.IsUnknown() {
		data.AcmeDirectoryUri = types.StringNull()
	}
	if data.CsrId.IsUnknown() {
		data.CsrId = types.Int64Null()
	}
	if data.Tos.IsUnknown() {
		data.Tos = types.BoolNull()
	}
	if data.DnsMapping.IsUnknown() {
		data.DnsMapping = types.StringNull()
	}
	if data.RenewDays.IsUnknown() {
		data.RenewDays = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("certificate.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read certificate: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["certificate"]; ok {
		if v == nil {
			data.Certificate = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Certificate = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Certificate = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Certificate = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["privatekey"]; ok {
		if v == nil {
			data.Privatekey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Privatekey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Privatekey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Privatekey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["CSR"]; ok {
		if v == nil {
			data.Csr = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Csr = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Csr = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Csr = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["key_length"]; ok {
		if v == nil {
			data.KeyLength = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.KeyLength = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.KeyLength = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["key_type"]; ok {
		switch val := v.(type) {
		case string:
			data.KeyType = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.KeyType = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.KeyType = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["city"]; ok {
		if v == nil {
			data.City = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.City = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.City = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.City = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["common"]; ok {
		if v == nil {
			data.Common = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Common = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Common = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Common = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["country"]; ok {
		if v == nil {
			data.Country = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Country = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Country = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Country = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["email"]; ok {
		if v == nil {
			data.Email = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Email = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Email = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Email = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["organization"]; ok {
		if v == nil {
			data.Organization = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Organization = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Organization = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Organization = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["organizational_unit"]; ok {
		if v == nil {
			data.OrganizationalUnit = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.OrganizationalUnit = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["state"]; ok {
		if v == nil {
			data.State = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.State = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.State = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.State = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["digest_algorithm"]; ok {
		switch val := v.(type) {
		case string:
			data.DigestAlgorithm = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["san"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.San, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if data.AddToTrustedStore.IsUnknown() {
		data.AddToTrustedStore = types.BoolNull()
	}
	if data.Certificate.IsUnknown() {
		data.Certificate = types.StringNull()
	}
	if data.Privatekey.IsUnknown() {
		data.Privatekey = types.StringNull()
	}
	if data.Csr.IsUnknown() {
		data.Csr = types.StringNull()
	}
	if data.KeyLength.IsUnknown() {
		data.KeyLength = types.Int64Null()
	}
	if data.KeyType.IsUnknown() {
		data.KeyType = types.StringNull()
	}
	if data.EcCurve.IsUnknown() {
		data.EcCurve = types.StringNull()
	}
	if data.Passphrase.IsUnknown() {
		data.Passphrase = types.StringNull()
	}
	if data.City.IsUnknown() {
		data.City = types.StringNull()
	}
	if data.Common.IsUnknown() {
		data.Common = types.StringNull()
	}
	if data.Country.IsUnknown() {
		data.Country = types.StringNull()
	}
	if data.Email.IsUnknown() {
		data.Email = types.StringNull()
	}
	if data.Organization.IsUnknown() {
		data.Organization = types.StringNull()
	}
	if data.OrganizationalUnit.IsUnknown() {
		data.OrganizationalUnit = types.StringNull()
	}
	if data.State.IsUnknown() {
		data.State = types.StringNull()
	}
	if data.DigestAlgorithm.IsUnknown() {
		data.DigestAlgorithm = types.StringNull()
	}
	if data.San.IsUnknown() {
		data.San, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.CertExtensions.IsUnknown() {
		data.CertExtensions = types.StringNull()
	}
	if data.AcmeDirectoryUri.IsUnknown() {
		data.AcmeDirectoryUri = types.StringNull()
	}
	if data.CsrId.IsUnknown() {
		data.CsrId = types.Int64Null()
	}
	if data.Tos.IsUnknown() {
		data.Tos = types.BoolNull()
	}
	if data.DnsMapping.IsUnknown() {
		data.DnsMapping = types.StringNull()
	}
	if data.RenewDays.IsUnknown() {
		data.RenewDays = types.Int64Null()
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.RenewDays.IsNull() && !data.RenewDays.IsUnknown() {
		params["renew_days"] = data.RenewDays.ValueInt64()
	}
	if !data.AddToTrustedStore.IsNull() && !data.AddToTrustedStore.IsUnknown() {
		params["add_to_trusted_store"] = data.AddToTrustedStore.ValueBool()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}

	_, err = r.client.CallWithJob("certificate.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update certificate: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("certificate.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back certificate: %s", readErr))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["certificate"]; ok {
		if v == nil {
			data.Certificate = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Certificate = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Certificate = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Certificate = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["privatekey"]; ok {
		if v == nil {
			data.Privatekey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Privatekey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Privatekey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Privatekey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["CSR"]; ok {
		if v == nil {
			data.Csr = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Csr = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Csr = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Csr = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["key_length"]; ok {
		if v == nil {
			data.KeyLength = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.KeyLength = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.KeyLength = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["key_type"]; ok {
		switch val := v.(type) {
		case string:
			data.KeyType = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.KeyType = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.KeyType = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["city"]; ok {
		if v == nil {
			data.City = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.City = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.City = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.City = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["common"]; ok {
		if v == nil {
			data.Common = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Common = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Common = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Common = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["country"]; ok {
		if v == nil {
			data.Country = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Country = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Country = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Country = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["email"]; ok {
		if v == nil {
			data.Email = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Email = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Email = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Email = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["organization"]; ok {
		if v == nil {
			data.Organization = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Organization = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Organization = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Organization = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["organizational_unit"]; ok {
		if v == nil {
			data.OrganizationalUnit = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.OrganizationalUnit = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.OrganizationalUnit = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["state"]; ok {
		if v == nil {
			data.State = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.State = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.State = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.State = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["digest_algorithm"]; ok {
		switch val := v.(type) {
		case string:
			data.DigestAlgorithm = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.DigestAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["san"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.San, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if data.AddToTrustedStore.IsUnknown() {
		data.AddToTrustedStore = types.BoolNull()
	}
	if data.Certificate.IsUnknown() {
		data.Certificate = types.StringNull()
	}
	if data.Privatekey.IsUnknown() {
		data.Privatekey = types.StringNull()
	}
	if data.Csr.IsUnknown() {
		data.Csr = types.StringNull()
	}
	if data.KeyLength.IsUnknown() {
		data.KeyLength = types.Int64Null()
	}
	if data.KeyType.IsUnknown() {
		data.KeyType = types.StringNull()
	}
	if data.EcCurve.IsUnknown() {
		data.EcCurve = types.StringNull()
	}
	if data.Passphrase.IsUnknown() {
		data.Passphrase = types.StringNull()
	}
	if data.City.IsUnknown() {
		data.City = types.StringNull()
	}
	if data.Common.IsUnknown() {
		data.Common = types.StringNull()
	}
	if data.Country.IsUnknown() {
		data.Country = types.StringNull()
	}
	if data.Email.IsUnknown() {
		data.Email = types.StringNull()
	}
	if data.Organization.IsUnknown() {
		data.Organization = types.StringNull()
	}
	if data.OrganizationalUnit.IsUnknown() {
		data.OrganizationalUnit = types.StringNull()
	}
	if data.State.IsUnknown() {
		data.State = types.StringNull()
	}
	if data.DigestAlgorithm.IsUnknown() {
		data.DigestAlgorithm = types.StringNull()
	}
	if data.San.IsUnknown() {
		data.San, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.CertExtensions.IsUnknown() {
		data.CertExtensions = types.StringNull()
	}
	if data.AcmeDirectoryUri.IsUnknown() {
		data.AcmeDirectoryUri = types.StringNull()
	}
	if data.CsrId.IsUnknown() {
		data.CsrId = types.Int64Null()
	}
	if data.Tos.IsUnknown() {
		data.Tos = types.BoolNull()
	}
	if data.DnsMapping.IsUnknown() {
		data.DnsMapping = types.StringNull()
	}
	if data.RenewDays.IsUnknown() {
		data.RenewDays = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	id = []interface{}{id, map[string]interface{}{}}

	_, err = r.client.CallWithJob("certificate.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete certificate: %s", err))
		return
	}
}
