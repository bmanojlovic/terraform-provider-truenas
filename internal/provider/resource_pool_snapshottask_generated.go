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
	Schedule types.Object `tfsdk:"schedule"`
}

func NewPoolSnapshottaskResource() resource.Resource {
	return &PoolSnapshottaskResource{}
}

func (r *PoolSnapshottaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_snapshottask"
}

func (r *PoolSnapshottaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS pool_snapshottask resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"dataset": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"recursive": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"lifetime_value": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"lifetime_unit": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"exclude": schema.ListAttribute{
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"naming_schema": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"allow_empty": schema.BoolAttribute{
				Required: false,
				Optional: true,
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
	params["dataset"] = data.Dataset.ValueString()
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
	if !data.NamingSchema.IsNull() {
		params["naming_schema"] = data.NamingSchema.ValueString()
	}
	if !data.AllowEmpty.IsNull() {
		params["allow_empty"] = data.AllowEmpty.ValueBool()
	}

	result, err := r.client.Call("pool.snapshottask.create", params)
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

func (r *PoolSnapshottaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolSnapshottaskResourceModel
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

	_, err = r.client.Call("pool.snapshottask.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshottaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolSnapshottaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state PoolSnapshottaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["dataset"] = data.Dataset.ValueString()
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
	if !data.NamingSchema.IsNull() {
		params["naming_schema"] = data.NamingSchema.ValueString()
	}
	if !data.AllowEmpty.IsNull() {
		params["allow_empty"] = data.AllowEmpty.ValueBool()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("pool.snapshottask.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshottaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolSnapshottaskResourceModel
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

	_, err = r.client.Call("pool.snapshottask.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
