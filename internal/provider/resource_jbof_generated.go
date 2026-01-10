package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type JbofResource struct {
	client *client.Client
}

type JbofResourceModel struct {
	ID types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	MgmtIp1 types.String `tfsdk:"mgmt_ip1"`
	MgmtIp2 types.String `tfsdk:"mgmt_ip2"`
	MgmtUsername types.String `tfsdk:"mgmt_username"`
	MgmtPassword types.String `tfsdk:"mgmt_password"`
}

func NewJbofResource() resource.Resource {
	return &JbofResource{}
}

func (r *JbofResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jbof"
}

func (r *JbofResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *JbofResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new JBOF.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional description of the JBOF.",
			},
			"mgmt_ip1": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "IP of first Redfish management interface.",
			},
			"mgmt_ip2": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional IP of second Redfish management interface.",
			},
			"mgmt_username": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Redfish administrative username.",
			},
			"mgmt_password": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Redfish administrative password.",
			},
		},
	}
}

func (r *JbofResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JbofResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.MgmtIp1.IsNull() {
		params["mgmt_ip1"] = data.MgmtIp1.ValueString()
	}
	if !data.MgmtIp2.IsNull() {
		params["mgmt_ip2"] = data.MgmtIp2.ValueString()
	}
	if !data.MgmtUsername.IsNull() {
		params["mgmt_username"] = data.MgmtUsername.ValueString()
	}
	if !data.MgmtPassword.IsNull() {
		params["mgmt_password"] = data.MgmtPassword.ValueString()
	}

	result, err := r.client.Call("jbof.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create jbof: %s", err))
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

func (r *JbofResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("jbof.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read jbof: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_ip1"]; ok && v != nil {
			data.MgmtIp1 = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_ip2"]; ok && v != nil {
			data.MgmtIp2 = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_username"]; ok && v != nil {
			data.MgmtUsername = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_password"]; ok && v != nil {
			data.MgmtPassword = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JbofResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state JbofResourceModel
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
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.MgmtIp1.IsNull() {
		params["mgmt_ip1"] = data.MgmtIp1.ValueString()
	}
	if !data.MgmtIp2.IsNull() {
		params["mgmt_ip2"] = data.MgmtIp2.ValueString()
	}
	if !data.MgmtUsername.IsNull() {
		params["mgmt_username"] = data.MgmtUsername.ValueString()
	}
	if !data.MgmtPassword.IsNull() {
		params["mgmt_password"] = data.MgmtPassword.ValueString()
	}

	_, err = r.client.Call("jbof.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update jbof: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JbofResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("jbof.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete jbof: %s", err))
		return
	}
}
