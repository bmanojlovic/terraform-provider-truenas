package provider

import (
	"context"
	"fmt"
	"strconv"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type PoolSnapshottaskResource struct {
	client *client.Client
}

type PoolSnapshottaskResourceModel struct {
	ID types.String `tfsdk:"id"`
	Dataset types.String `tfsdk:"dataset"`
	Recursive types.Bool `tfsdk:"recursive"`
	LifetimeValue types.Int64 `tfsdk:"lifetime_value"`
	LifetimeUnit types.String `tfsdk:"lifetime_unit"`
	Enabled types.Bool `tfsdk:"enabled"`
	Exclude types.List `tfsdk:"exclude"`
	NamingSchema types.String `tfsdk:"naming_schema"`
	AllowEmpty types.Bool `tfsdk:"allow_empty"`
	Schedule types.String `tfsdk:"schedule"`
	FixateRemovalDate types.Bool `tfsdk:"fixate_removal_date"`
}

func NewPoolSnapshottaskResource() resource.Resource {
	return &PoolSnapshottaskResource{}
}

func (r *PoolSnapshottaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_snapshottask"
}

func (r *PoolSnapshottaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PoolSnapshottaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Periodic Snapshot Task",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"dataset": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "The dataset to take snapshots of.",
			},
			"recursive": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to recursively snapshot child datasets.",
			},
			"lifetime_value": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of time units to retain snapshots. `lifetime_unit` gives the time unit.",
			},
			"lifetime_unit": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Unit of time for snapshot retention.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether this periodic snapshot task is enabled.",
			},
			"exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of dataset patterns to exclude from recursive snapshots.",
			},
			"naming_schema": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Naming pattern for generated snapshots using strftime format.",
			},
			"allow_empty": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to take snapshots even if no data has changed.",
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Cron schedule for when snapshots should be taken.",
			},
			"fixate_removal_date": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to fix the removal date of existing snapshots when retention settings change.",
			},
		},
	}
}

func (r *PoolSnapshottaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolSnapshottaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolSnapshottaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Dataset.IsNull() {
		params["dataset"] = data.Dataset.ValueString()
	}
	if !data.Recursive.IsNull() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.LifetimeValue.IsNull() {
		params["lifetime_value"] = data.LifetimeValue.ValueInt64()
	}
	if !data.LifetimeUnit.IsNull() {
		params["lifetime_unit"] = data.LifetimeUnit.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.NamingSchema.IsNull() {
		params["naming_schema"] = data.NamingSchema.ValueString()
	}
	if !data.AllowEmpty.IsNull() {
		params["allow_empty"] = data.AllowEmpty.ValueBool()
	}
	if !data.Schedule.IsNull() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.FixateRemovalDate.IsNull() {
		params["fixate_removal_date"] = data.FixateRemovalDate.ValueBool()
	}

	result, err := r.client.Call("pool.snapshottask.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create pool_snapshottask: %s", err))
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

func (r *PoolSnapshottaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolSnapshottaskResourceModel
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

	result, err := r.client.Call("pool.snapshottask.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool_snapshottask: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["dataset"]; ok && v != nil {
			data.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recursive"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Recursive = types.BoolValue(bv) }
		}
		if v, ok := resultMap["lifetime_value"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.LifetimeValue = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["lifetime_unit"]; ok && v != nil {
			data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["exclude"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Exclude, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["naming_schema"]; ok && v != nil {
			data.NamingSchema = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["allow_empty"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AllowEmpty = types.BoolValue(bv) }
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["fixate_removal_date"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.FixateRemovalDate = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshottaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolSnapshottaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PoolSnapshottaskResourceModel
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
	if !data.Dataset.IsNull() {
		params["dataset"] = data.Dataset.ValueString()
	}
	if !data.Recursive.IsNull() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.LifetimeValue.IsNull() {
		params["lifetime_value"] = data.LifetimeValue.ValueInt64()
	}
	if !data.LifetimeUnit.IsNull() {
		params["lifetime_unit"] = data.LifetimeUnit.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.NamingSchema.IsNull() {
		params["naming_schema"] = data.NamingSchema.ValueString()
	}
	if !data.AllowEmpty.IsNull() {
		params["allow_empty"] = data.AllowEmpty.ValueBool()
	}
	if !data.Schedule.IsNull() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.FixateRemovalDate.IsNull() {
		params["fixate_removal_date"] = data.FixateRemovalDate.ValueBool()
	}

	_, err = r.client.Call("pool.snapshottask.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update pool_snapshottask: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshottaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolSnapshottaskResourceModel
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

	_, err = r.client.Call("pool.snapshottask.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool_snapshottask: %s", err))
		return
	}
}
