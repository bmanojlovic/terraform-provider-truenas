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

type StaticrouteResource struct {
	client *client.Client
}

type StaticrouteResourceModel struct {
	ID types.String `tfsdk:"id"`
	Destination types.String `tfsdk:"destination"`
	Gateway types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
}

func NewStaticrouteResource() resource.Resource {
	return &StaticrouteResource{}
}

func (r *StaticrouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staticroute"
}

func (r *StaticrouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *StaticrouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Static Route.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"destination": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Destination network or host for this static route.",
			},
			"gateway": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Gateway IP address for this static route.",
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional description for this static route.",
			},
		},
	}
}

func (r *StaticrouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StaticrouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StaticrouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Destination.IsNull() {
		params["destination"] = data.Destination.ValueString()
	}
	if !data.Gateway.IsNull() {
		params["gateway"] = data.Gateway.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}

	result, err := r.client.Call("staticroute.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create staticroute: %s", err))
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

func (r *StaticrouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StaticrouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("staticroute.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read staticroute: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["destination"]; ok && v != nil {
			data.Destination = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["gateway"]; ok && v != nil {
			data.Gateway = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticrouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StaticrouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StaticrouteResourceModel
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
	if !data.Destination.IsNull() {
		params["destination"] = data.Destination.ValueString()
	}
	if !data.Gateway.IsNull() {
		params["gateway"] = data.Gateway.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}

	_, err = r.client.Call("staticroute.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update staticroute: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticrouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StaticrouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("staticroute.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete staticroute: %s", err))
		return
	}
}
