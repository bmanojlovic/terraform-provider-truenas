package provider

import (
	"context"
	"fmt"
	"strconv"

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
	Csr types.String `tfsdk:"CSR"`
	KeyLength types.String `tfsdk:"key_length"`
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
	CertExtensions types.Object `tfsdk:"cert_extensions"`
	AcmeDirectoryUri types.String `tfsdk:"acme_directory_uri"`
	CsrId types.String `tfsdk:"csr_id"`
	Tos types.String `tfsdk:"tos"`
	DnsMapping types.Object `tfsdk:"dns_mapping"`
	RenewDays types.Int64 `tfsdk:"renew_days"`
}

func NewCertificateResource() resource.Resource {
	return &CertificateResource{}
}

func (r *CertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *CertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS certificate resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"create_type": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"add_to_trusted_store": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"certificate": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"privatekey": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"CSR": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"key_length": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"key_type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"ec_curve": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"passphrase": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"city": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"common": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"country": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"email": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"organization": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"organizational_unit": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"state": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"digest_algorithm": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"san": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"cert_extensions": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"acme_directory_uri": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"csr_id": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"tos": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"dns_mapping": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"renew_days": schema.Int64Attribute{
				Required: false,
				Optional: true,
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
	params["name"] = data.Name.ValueString()
	params["create_type"] = data.CreateType.ValueString()
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
		params["key_length"] = data.KeyLength.ValueString()
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
	if !data.AcmeDirectoryUri.IsNull() {
		params["acme_directory_uri"] = data.AcmeDirectoryUri.ValueString()
	}
	if !data.CsrId.IsNull() {
		params["csr_id"] = data.CsrId.ValueString()
	}
	if !data.Tos.IsNull() {
		params["tos"] = data.Tos.ValueString()
	}
	if !data.RenewDays.IsNull() {
		params["renew_days"] = data.RenewDays.ValueInt64()
	}

	result, err := r.client.Call("certificate.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

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

	_, err := r.client.Call("certificate.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
	params["create_type"] = data.CreateType.ValueString()
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
		params["key_length"] = data.KeyLength.ValueString()
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
	if !data.AcmeDirectoryUri.IsNull() {
		params["acme_directory_uri"] = data.AcmeDirectoryUri.ValueString()
	}
	if !data.CsrId.IsNull() {
		params["csr_id"] = data.CsrId.ValueString()
	}
	if !data.Tos.IsNull() {
		params["tos"] = data.Tos.ValueString()
	}
	if !data.RenewDays.IsNull() {
		params["renew_days"] = data.RenewDays.ValueInt64()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("certificate.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("certificate.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
