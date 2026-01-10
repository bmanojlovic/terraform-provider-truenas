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

type FcFcHostResource struct {
	client *client.Client
}

type FcFcHostResourceModel struct {
	ID types.String `tfsdk:"id"`
	Alias types.String `tfsdk:"alias"`
	Wwpn types.String `tfsdk:"wwpn"`
	WwpnB types.String `tfsdk:"wwpn_b"`
	Npiv types.Int64 `tfsdk:"npiv"`
}

func NewFcFcHostResource() resource.Resource {
	return &FcFcHostResource{}
}

func (r *FcFcHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fc_fc_host"
}

func (r *FcFcHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FcFcHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates FC host (pairing).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"alias": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Human-readable alias for the Fibre Channel host.",
			},
			"wwpn": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "World Wide Port Name for port A or `null` if not configured.",
			},
			"wwpn_b": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "World Wide Port Name for port B or `null` if not configured.",
			},
			"npiv": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of N_Port ID Virtualization (NPIV) virtual ports to create.",
			},
		},
	}
}

func (r *FcFcHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FcFcHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Alias.IsNull() {
		params["alias"] = data.Alias.ValueString()
	}
	if !data.Wwpn.IsNull() {
		params["wwpn"] = data.Wwpn.ValueString()
	}
	if !data.WwpnB.IsNull() {
		params["wwpn_b"] = data.WwpnB.ValueString()
	}
	if !data.Npiv.IsNull() {
		params["npiv"] = data.Npiv.ValueInt64()
	}

	result, err := r.client.Call("fc.fc_host.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create fc_fc_host: %s", err))
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

func (r *FcFcHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("fc.fc_host.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read fc_fc_host: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["alias"]; ok && v != nil {
			data.Alias = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["wwpn"]; ok && v != nil {
			data.Wwpn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["wwpn_b"]; ok && v != nil {
			data.WwpnB = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["npiv"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Npiv = types.Int64Value(int64(fv)) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcFcHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FcFcHostResourceModel
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
	if !data.Alias.IsNull() {
		params["alias"] = data.Alias.ValueString()
	}
	if !data.Wwpn.IsNull() {
		params["wwpn"] = data.Wwpn.ValueString()
	}
	if !data.WwpnB.IsNull() {
		params["wwpn_b"] = data.WwpnB.ValueString()
	}
	if !data.Npiv.IsNull() {
		params["npiv"] = data.Npiv.ValueInt64()
	}

	_, err = r.client.Call("fc.fc_host.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update fc_fc_host: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcFcHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("fc.fc_host.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete fc_fc_host: %s", err))
		return
	}
}
