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

type IscsiTargetextentResource struct {
	client *client.Client
}

type IscsiTargetextentResourceModel struct {
	ID types.String `tfsdk:"id"`
	Target types.Int64 `tfsdk:"target"`
	Lunid types.Int64 `tfsdk:"lunid"`
	Extent types.Int64 `tfsdk:"extent"`
}

func NewIscsiTargetextentResource() resource.Resource {
	return &IscsiTargetextentResource{}
}

func (r *IscsiTargetextentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_targetextent"
}

func (r *IscsiTargetextentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiTargetextentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an Associated Target.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"target": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the iSCSI target to associate with the extent.",
			},
			"lunid": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "LUN ID to assign or `null` to auto-assign the next available LUN.",
			},
			"extent": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the iSCSI extent to associate with the target.",
			},
		},
	}
}

func (r *IscsiTargetextentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiTargetextentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Target.IsNull() {
		params["target"] = data.Target.ValueInt64()
	}
	if !data.Lunid.IsNull() {
		params["lunid"] = data.Lunid.ValueInt64()
	}
	if !data.Extent.IsNull() {
		params["extent"] = data.Extent.ValueInt64()
	}

	result, err := r.client.Call("iscsi.targetextent.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_targetextent: %s", err))
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

func (r *IscsiTargetextentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("iscsi.targetextent.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_targetextent: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["target"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Target = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["lunid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Lunid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["extent"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Extent = types.Int64Value(int64(fv)) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetextentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiTargetextentResourceModel
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
	if !data.Target.IsNull() {
		params["target"] = data.Target.ValueInt64()
	}
	if !data.Lunid.IsNull() {
		params["lunid"] = data.Lunid.ValueInt64()
	}
	if !data.Extent.IsNull() {
		params["extent"] = data.Extent.ValueInt64()
	}

	_, err = r.client.Call("iscsi.targetextent.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_targetextent: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetextentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("iscsi.targetextent.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_targetextent: %s", err))
		return
	}
}
