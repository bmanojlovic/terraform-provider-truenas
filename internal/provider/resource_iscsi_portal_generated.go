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

type IscsiPortalResource struct {
	client *client.Client
}

type IscsiPortalResourceModel struct {
	ID types.String `tfsdk:"id"`
	Listen types.List `tfsdk:"listen"`
	Comment types.String `tfsdk:"comment"`
}

func NewIscsiPortalResource() resource.Resource {
	return &IscsiPortalResource{}
}

func (r *IscsiPortalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_portal"
}

func (r *IscsiPortalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS iscsi_portal resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"listen": schema.ListAttribute{
				ElementType: types.StringType,
				Required: true,
				Optional: false,
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *IscsiPortalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiPortalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiPortalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	result, err := r.client.Call("iscsi.portal.create", params)
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

func (r *IscsiPortalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiPortalResourceModel
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

	_, err = r.client.Call("iscsi.portal.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiPortalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiPortalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state IscsiPortalResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("iscsi.portal.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiPortalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiPortalResourceModel
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

	_, err = r.client.Call("iscsi.portal.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
