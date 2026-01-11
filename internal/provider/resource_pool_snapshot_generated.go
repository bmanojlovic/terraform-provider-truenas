package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type PoolSnapshotResource struct {
	client *client.Client
}

type PoolSnapshotResourceModel struct {
	ID types.String `tfsdk:"id"`
	Dataset types.String `tfsdk:"dataset"`
	Recursive types.Bool `tfsdk:"recursive"`
	Exclude types.List `tfsdk:"exclude"`
	VmwareSync types.Bool `tfsdk:"vmware_sync"`
	Properties types.String `tfsdk:"properties"`
	Name types.String `tfsdk:"name"`
	NamingSchema types.String `tfsdk:"naming_schema"`
	UserPropertiesUpdate types.List `tfsdk:"user_properties_update"`
	UserPropertiesRemove types.List `tfsdk:"user_properties_remove"`
}

func NewPoolSnapshotResource() resource.Resource {
	return &PoolSnapshotResource{}
}

func (r *PoolSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_snapshot"
}

func (r *PoolSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PoolSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Take a snapshot from a given dataset.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"dataset": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name of the dataset to create a snapshot of.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"recursive": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to recursively snapshot child datasets.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of dataset patterns to exclude from recursive snapshots.",
			},
			"vmware_sync": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to sync VMware VMs before taking the snapshot.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"properties": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Object mapping ZFS property names to values to set on the snapshot.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Explicit name for the snapshot.",
			},
			"naming_schema": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Naming schema pattern to generate the snapshot name automatically.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_properties_update": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Properties to update.",
			},
			"user_properties_remove": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Properties to remove.",
			},
		},
	}
}

func (r *PoolSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolSnapshotResourceModel
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
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.VmwareSync.IsNull() {
		params["vmware_sync"] = data.VmwareSync.ValueBool()
	}
	if !data.Properties.IsNull() {
		var propertiesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Properties.ValueString()), &propertiesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse properties: %s", err))
			return
		}
		params["properties"] = propertiesObj
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.NamingSchema.IsNull() {
		params["naming_schema"] = data.NamingSchema.ValueString()
	}
	if !data.UserPropertiesUpdate.IsNull() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		params["user_properties_update"] = user_properties_updateList
	}
	if !data.UserPropertiesRemove.IsNull() {
		var user_properties_removeList []string
		data.UserPropertiesRemove.ElementsAs(ctx, &user_properties_removeList, false)
		params["user_properties_remove"] = user_properties_removeList
	}

	result, err := r.client.Call("pool.snapshot.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create pool_snapshot: %s", err))
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

func (r *PoolSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	result, err := r.client.Call("pool.snapshot.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool_snapshot: %s", err))
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
		if v, ok := resultMap["exclude"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Exclude, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["vmware_sync"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.VmwareSync = types.BoolValue(bv) }
		}
		if v, ok := resultMap["properties"]; ok && v != nil {
			data.Properties = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["naming_schema"]; ok && v != nil {
			data.NamingSchema = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user_properties_update"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.UserPropertiesUpdate, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["user_properties_remove"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.UserPropertiesRemove, _ = types.ListValue(types.StringType, strVals)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolSnapshotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PoolSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = state.ID.ValueString()

	params := map[string]interface{}{}
	if !data.UserPropertiesUpdate.IsNull() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		params["user_properties_update"] = user_properties_updateList
	}
	if !data.UserPropertiesRemove.IsNull() {
		var user_properties_removeList []string
		data.UserPropertiesRemove.ElementsAs(ctx, &user_properties_removeList, false)
		params["user_properties_remove"] = user_properties_removeList
	}

	_, err = r.client.Call("pool.snapshot.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update pool_snapshot: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	_, err = r.client.Call("pool.snapshot.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool_snapshot: %s", err))
		return
	}
}
