package provider

import (
	"context"
	"fmt"

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

func (r *FcFcHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS fc_fc_host resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"alias": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"wwpn": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"wwpn_b": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"npiv": schema.Int64Attribute{
				Required: false,
				Optional: true,
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

	params := map[string]interface{}{
		"alias": data.Alias.ValueString(),
		"wwpn": data.Wwpn.ValueString(),
		"wwpn_b": data.WwpnB.ValueString(),
		"npiv": data.Npiv.ValueInt64(),
	}

	result, err := r.client.Call("fc/fc_host.create", params)
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

func (r *FcFcHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("fc/fc_host.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcFcHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"alias": data.Alias.ValueString(),
		"wwpn": data.Wwpn.ValueString(),
		"wwpn_b": data.WwpnB.ValueString(),
		"npiv": data.Npiv.ValueInt64(),
	}

	_, err := r.client.Call("fc/fc_host.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcFcHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FcFcHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("fc/fc_host.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
