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

type NvmetHostResource struct {
	client *client.Client
}

type NvmetHostResourceModel struct {
	ID types.String `tfsdk:"id"`
	Hostnqn types.String `tfsdk:"hostnqn"`
	DhchapKey types.String `tfsdk:"dhchap_key"`
	DhchapCtrlKey types.String `tfsdk:"dhchap_ctrl_key"`
	DhchapDhgroup types.String `tfsdk:"dhchap_dhgroup"`
	DhchapHash types.String `tfsdk:"dhchap_hash"`
}

func NewNvmetHostResource() resource.Resource {
	return &NvmetHostResource{}
}

func (r *NvmetHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_host"
}

func (r *NvmetHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS nvmet_host resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"hostnqn": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"dhchap_key": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"dhchap_ctrl_key": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"dhchap_dhgroup": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"dhchap_hash": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *NvmetHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NvmetHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["hostnqn"] = data.Hostnqn.ValueString()
	if !data.DhchapKey.IsNull() {
		params["dhchap_key"] = data.DhchapKey.ValueString()
	}
	if !data.DhchapCtrlKey.IsNull() {
		params["dhchap_ctrl_key"] = data.DhchapCtrlKey.ValueString()
	}
	if !data.DhchapDhgroup.IsNull() {
		params["dhchap_dhgroup"] = data.DhchapDhgroup.ValueString()
	}
	if !data.DhchapHash.IsNull() {
		params["dhchap_hash"] = data.DhchapHash.ValueString()
	}

	result, err := r.client.Call("nvmet.host.create", params)
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

func (r *NvmetHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("nvmet.host.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state NvmetHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["hostnqn"] = data.Hostnqn.ValueString()
	if !data.DhchapKey.IsNull() {
		params["dhchap_key"] = data.DhchapKey.ValueString()
	}
	if !data.DhchapCtrlKey.IsNull() {
		params["dhchap_ctrl_key"] = data.DhchapCtrlKey.ValueString()
	}
	if !data.DhchapDhgroup.IsNull() {
		params["dhchap_dhgroup"] = data.DhchapDhgroup.ValueString()
	}
	if !data.DhchapHash.IsNull() {
		params["dhchap_hash"] = data.DhchapHash.ValueString()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("nvmet.host.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("nvmet.host.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
