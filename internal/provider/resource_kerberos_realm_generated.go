package provider

import (
	"context"
	"fmt"


	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type KerberosRealmResource struct {
	client *client.Client
}

type KerberosRealmResourceModel struct {
	ID types.String `tfsdk:"id"`
	Realm types.String `tfsdk:"realm"`
	PrimaryKdc types.String `tfsdk:"primary_kdc"`
	Kdc types.List `tfsdk:"kdc"`
	AdminServer types.List `tfsdk:"admin_server"`
	KpasswdServer types.List `tfsdk:"kpasswd_server"`
}

func NewKerberosRealmResource() resource.Resource {
	return &KerberosRealmResource{}
}

func (r *KerberosRealmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_realm"
}

func (r *KerberosRealmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS kerberos_realm resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"realm": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"primary_kdc": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"kdc": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"admin_server": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"kpasswd_server": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *KerberosRealmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KerberosRealmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["realm"] = data.Realm.ValueString()
	if !data.PrimaryKdc.IsNull() {
		params["primary_kdc"] = data.PrimaryKdc.ValueString()
	}

	result, err := r.client.Call("kerberos/realm.create", params)
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

func (r *KerberosRealmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("kerberos/realm.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KerberosRealmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["realm"] = data.Realm.ValueString()
	if !data.PrimaryKdc.IsNull() {
		params["primary_kdc"] = data.PrimaryKdc.ValueString()
	}

	_, err := r.client.Call("kerberos/realm.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KerberosRealmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("kerberos/realm.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
