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

func (r *JbofResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS jbof resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"mgmt_ip1": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"mgmt_ip2": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"mgmt_username": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"mgmt_password": schema.StringAttribute{
				Required: true,
				Optional: false,
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
	params["mgmt_ip1"] = data.MgmtIp1.ValueString()
	if !data.MgmtIp2.IsNull() {
		params["mgmt_ip2"] = data.MgmtIp2.ValueString()
	}
	params["mgmt_username"] = data.MgmtUsername.ValueString()
	params["mgmt_password"] = data.MgmtPassword.ValueString()

	result, err := r.client.Call("jbof.create", params)
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

func (r *JbofResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("jbof.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JbofResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state JbofResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	params["mgmt_ip1"] = data.MgmtIp1.ValueString()
	if !data.MgmtIp2.IsNull() {
		params["mgmt_ip2"] = data.MgmtIp2.ValueString()
	}
	params["mgmt_username"] = data.MgmtUsername.ValueString()
	params["mgmt_password"] = data.MgmtPassword.ValueString()

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("jbof.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JbofResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JbofResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("jbof.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
