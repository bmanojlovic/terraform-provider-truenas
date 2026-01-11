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

func (r *KerberosRealmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *KerberosRealmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new kerberos realm. This will be automatically populated during the",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"realm": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Kerberos realm name. This is external to TrueNAS and is case-sensitive.     The general convention f",
			},
			"primary_kdc": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The master Kerberos domain controller for this realm. TrueNAS uses this as a fallback if it cannot g",
			},
			"kdc": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of kerberos domain controllers. If the list is empty then the kerberos     libraries will use D",
			},
			"admin_server": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of kerberos admin servers. If the list is empty then the kerberos     libraries will use DNS to",
			},
			"kpasswd_server": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of kerberos kpasswd servers. If the list is empty then DNS will be used     to look them up if ",
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
	if !data.Realm.IsNull() {
		params["realm"] = data.Realm.ValueString()
	}
	if !data.PrimaryKdc.IsNull() {
		params["primary_kdc"] = data.PrimaryKdc.ValueString()
	}
	if !data.Kdc.IsNull() {
		var kdcList []string
		data.Kdc.ElementsAs(ctx, &kdcList, false)
		params["kdc"] = kdcList
	}
	if !data.AdminServer.IsNull() {
		var admin_serverList []string
		data.AdminServer.ElementsAs(ctx, &admin_serverList, false)
		params["admin_server"] = admin_serverList
	}
	if !data.KpasswdServer.IsNull() {
		var kpasswd_serverList []string
		data.KpasswdServer.ElementsAs(ctx, &kpasswd_serverList, false)
		params["kpasswd_server"] = kpasswd_serverList
	}

	result, err := r.client.Call("kerberos.realm.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create kerberos_realm: %s", err))
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

func (r *KerberosRealmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	result, err := r.client.Call("kerberos.realm.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read kerberos_realm: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["realm"]; ok && v != nil {
			data.Realm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["primary_kdc"]; ok && v != nil {
			data.PrimaryKdc = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["kdc"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Kdc, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["admin_server"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.AdminServer, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["kpasswd_server"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.KpasswdServer, _ = types.ListValue(types.StringType, strVals)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KerberosRealmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state KerberosRealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	params := map[string]interface{}{}
	if !data.Realm.IsNull() {
		params["realm"] = data.Realm.ValueString()
	}
	if !data.PrimaryKdc.IsNull() {
		params["primary_kdc"] = data.PrimaryKdc.ValueString()
	}
	if !data.Kdc.IsNull() {
		var kdcList []string
		data.Kdc.ElementsAs(ctx, &kdcList, false)
		params["kdc"] = kdcList
	}
	if !data.AdminServer.IsNull() {
		var admin_serverList []string
		data.AdminServer.ElementsAs(ctx, &admin_serverList, false)
		params["admin_server"] = admin_serverList
	}
	if !data.KpasswdServer.IsNull() {
		var kpasswd_serverList []string
		data.KpasswdServer.ElementsAs(ctx, &kpasswd_serverList, false)
		params["kpasswd_server"] = kpasswd_serverList
	}

	_, err = r.client.Call("kerberos.realm.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update kerberos_realm: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KerberosRealmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KerberosRealmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	_, err = r.client.Call("kerberos.realm.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete kerberos_realm: %s", err))
		return
	}
}
