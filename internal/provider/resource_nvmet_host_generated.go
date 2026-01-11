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

func (r *NvmetHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NvmetHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an NVMe target `host`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"hostnqn": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "NQN of the host that will connect to this TrueNAS. ",
			},
			"dhchap_key": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "If set, the secret that the host must present when connecting.  A suitable secret can be generated u",
			},
			"dhchap_ctrl_key": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "If set, the secret that this TrueNAS will present to the host when the host is connecting (Bi-Direct",
			},
			"dhchap_dhgroup": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "If selected, the DH (Diffie-Hellman) key exchange built on top of CHAP to be used for authentication",
			},
			"dhchap_hash": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "HMAC (Hashed Message Authentication Code) to be used in conjunction if a `dhchap_dhgroup` is selecte",
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
	if !data.Hostnqn.IsNull() {
		params["hostnqn"] = data.Hostnqn.ValueString()
	}
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
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create nvmet_host: %s", err))
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

func (r *NvmetHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetHostResourceModel
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

	result, err := r.client.Call("nvmet.host.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read nvmet_host: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["hostnqn"]; ok && v != nil {
			data.Hostnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_key"]; ok && v != nil {
			data.DhchapKey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_ctrl_key"]; ok && v != nil {
			data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_dhgroup"]; ok && v != nil {
			data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_hash"]; ok && v != nil {
			data.DhchapHash = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NvmetHostResourceModel
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
	if !data.Hostnqn.IsNull() {
		params["hostnqn"] = data.Hostnqn.ValueString()
	}
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

	_, err = r.client.Call("nvmet.host.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update nvmet_host: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetHostResourceModel
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

	_, err = r.client.Call("nvmet.host.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete nvmet_host: %s", err))
		return
	}
}
