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

type IscsiAuthResource struct {
	client *client.Client
}

type IscsiAuthResourceModel struct {
	ID types.String `tfsdk:"id"`
	Tag types.Int64 `tfsdk:"tag"`
	User types.String `tfsdk:"user"`
	Secret types.String `tfsdk:"secret"`
	Peeruser types.String `tfsdk:"peeruser"`
	Peersecret types.String `tfsdk:"peersecret"`
	DiscoveryAuth types.String `tfsdk:"discovery_auth"`
}

func NewIscsiAuthResource() resource.Resource {
	return &IscsiAuthResource{}
}

func (r *IscsiAuthResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_auth"
}

func (r *IscsiAuthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS iscsi_auth resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"tag": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"user": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"secret": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"peeruser": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"peersecret": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"discovery_auth": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *IscsiAuthResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["tag"] = data.Tag.ValueInt64()
	params["user"] = data.User.ValueString()
	params["secret"] = data.Secret.ValueString()
	if !data.Peeruser.IsNull() {
		params["peeruser"] = data.Peeruser.ValueString()
	}
	if !data.Peersecret.IsNull() {
		params["peersecret"] = data.Peersecret.ValueString()
	}
	if !data.DiscoveryAuth.IsNull() {
		params["discovery_auth"] = data.DiscoveryAuth.ValueString()
	}

	result, err := r.client.Call("iscsi/auth.create", params)
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

func (r *IscsiAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("iscsi/auth.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state IscsiAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["tag"] = data.Tag.ValueInt64()
	params["user"] = data.User.ValueString()
	params["secret"] = data.Secret.ValueString()
	if !data.Peeruser.IsNull() {
		params["peeruser"] = data.Peeruser.ValueString()
	}
	if !data.Peersecret.IsNull() {
		params["peersecret"] = data.Peersecret.ValueString()
	}
	if !data.DiscoveryAuth.IsNull() {
		params["discovery_auth"] = data.DiscoveryAuth.ValueString()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("iscsi/auth.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("iscsi/auth.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
