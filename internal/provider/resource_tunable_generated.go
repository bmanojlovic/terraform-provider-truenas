package provider

import (
	"context"
	"fmt"


	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type TunableResource struct {
	client *client.Client
}

type TunableResourceModel struct {
	ID types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Var types.String `tfsdk:"var"`
	Value types.String `tfsdk:"value"`
	Comment types.String `tfsdk:"comment"`
	Enabled types.Bool `tfsdk:"enabled"`
	UpdateInitramfs types.Bool `tfsdk:"update_initramfs"`
}

func NewTunableResource() resource.Resource {
	return &TunableResource{}
}

func (r *TunableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunable"
}

func (r *TunableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS tunable resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"var": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"value": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"update_initramfs": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *TunableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TunableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TunableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	params["var"] = data.Var.ValueString()
	params["value"] = data.Value.ValueString()
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.UpdateInitramfs.IsNull() {
		params["update_initramfs"] = data.UpdateInitramfs.ValueBool()
	}

	result, err := r.client.Call("tunable.create", params)
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

func (r *TunableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TunableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("tunable.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TunableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	params["var"] = data.Var.ValueString()
	params["value"] = data.Value.ValueString()
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.UpdateInitramfs.IsNull() {
		params["update_initramfs"] = data.UpdateInitramfs.ValueBool()
	}

	_, err := r.client.Call("tunable.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TunableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("tunable.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
