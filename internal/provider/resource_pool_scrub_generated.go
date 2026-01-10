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

type PoolScrubResource struct {
	client *client.Client
}

type PoolScrubResourceModel struct {
	ID types.String `tfsdk:"id"`
	Pool types.Int64 `tfsdk:"pool"`
	Threshold types.Int64 `tfsdk:"threshold"`
	Description types.String `tfsdk:"description"`
	Schedule types.String `tfsdk:"schedule"`
	Enabled types.Bool `tfsdk:"enabled"`
}

func NewPoolScrubResource() resource.Resource {
	return &PoolScrubResource{}
}

func (r *PoolScrubResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_scrub"
}

func (r *PoolScrubResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PoolScrubResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a scrub task for a pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"pool": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the pool to scrub.",
			},
			"threshold": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Days before a scrub is due when a scrub should automatically start.",
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Description or notes for this scrub schedule.",
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Cron schedule for when scrubs should run.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether this scrub schedule is enabled.",
			},
		},
	}
}

func (r *PoolScrubResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolScrubResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolScrubResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Pool.IsNull() {
		params["pool"] = data.Pool.ValueInt64()
	}
	if !data.Threshold.IsNull() {
		params["threshold"] = data.Threshold.ValueInt64()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Schedule.IsNull() {
		params["schedule"] = data.Schedule.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	result, err := r.client.Call("pool.scrub.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create pool_scrub: %s", err))
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

func (r *PoolScrubResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolScrubResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("pool.scrub.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool_scrub: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["pool"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Pool = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["threshold"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Threshold = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolScrubResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolScrubResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PoolScrubResourceModel
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
	if !data.Pool.IsNull() {
		params["pool"] = data.Pool.ValueInt64()
	}
	if !data.Threshold.IsNull() {
		params["threshold"] = data.Threshold.ValueInt64()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Schedule.IsNull() {
		params["schedule"] = data.Schedule.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	_, err = r.client.Call("pool.scrub.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update pool_scrub: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolScrubResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolScrubResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("pool.scrub.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool_scrub: %s", err))
		return
	}
}
