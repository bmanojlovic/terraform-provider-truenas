package provider

import (
	"context"
	"fmt"

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

	params := map[string]interface{}{
		"hostnqn": data.Hostnqn.ValueString(),
		"dhchap_key": data.DhchapKey.ValueString(),
		"dhchap_ctrl_key": data.DhchapCtrlKey.ValueString(),
		"dhchap_dhgroup": data.DhchapDhgroup.ValueString(),
		"dhchap_hash": data.DhchapHash.ValueString(),
	}

	result, err := r.client.Call("nvmet/host.create", params)
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

	_, err := r.client.Call("nvmet/host.get_instance", data.ID.ValueString())
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

	params := map[string]interface{}{
		"hostnqn": data.Hostnqn.ValueString(),
		"dhchap_key": data.DhchapKey.ValueString(),
		"dhchap_ctrl_key": data.DhchapCtrlKey.ValueString(),
		"dhchap_dhgroup": data.DhchapDhgroup.ValueString(),
		"dhchap_hash": data.DhchapHash.ValueString(),
	}

	_, err := r.client.Call("nvmet/host.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("nvmet/host.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
