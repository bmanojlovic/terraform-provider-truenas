package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type PoolSnapshotResource struct {
	client *client.Client
}

type PoolSnapshotResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Dataset              types.String `tfsdk:"dataset"`
	Recursive            types.Bool   `tfsdk:"recursive"`
	Exclude              types.List   `tfsdk:"exclude"`
	VmwareSync           types.Bool   `tfsdk:"vmware_sync"`
	Properties           types.String `tfsdk:"properties"`
	Name                 types.String `tfsdk:"name"`
	NamingSchema         types.String `tfsdk:"naming_schema"`
	UserPropertiesUpdate types.List   `tfsdk:"user_properties_update"`
	UserPropertiesRemove types.List   `tfsdk:"user_properties_remove"`
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
				Required:      true,
				Optional:      false,
				Description:   "Name of the dataset to create a snapshot of.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"recursive": schema.BoolAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Whether to recursively snapshot child datasets.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"exclude": schema.ListAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "Array of dataset patterns to exclude from recursive snapshots.",
			},
			"vmware_sync": schema.BoolAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Whether to sync VMware VMs before taking the snapshot.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"properties": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Object mapping ZFS property names to values to set on the snapshot.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Explicit name for the snapshot.",
			},
			"naming_schema": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Naming schema pattern to generate the snapshot name automatically.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_properties_update": schema.ListAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "Properties to update.",
			},
			"user_properties_remove": schema.ListAttribute{
				Required:    false,
				Optional:    true,
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
	if !data.Dataset.IsNull() && !data.Dataset.IsUnknown() {
		params["dataset"] = data.Dataset.ValueString()
	}
	if !data.Recursive.IsNull() && !data.Recursive.IsUnknown() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Exclude.IsNull() && !data.Exclude.IsUnknown() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.VmwareSync.IsNull() && !data.VmwareSync.IsUnknown() {
		params["vmware_sync"] = data.VmwareSync.ValueBool()
	}
	if !data.Properties.IsNull() && !data.Properties.IsUnknown() {
		var propertiesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Properties.ValueString()), &propertiesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse properties: %s", err))
			return
		}
		params["properties"] = propertiesObj
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.NamingSchema.IsNull() && !data.NamingSchema.IsUnknown() {
		params["naming_schema"] = data.NamingSchema.ValueString()
	}
	if !data.UserPropertiesUpdate.IsNull() && !data.UserPropertiesUpdate.IsUnknown() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		var user_properties_updateObjs []map[string]interface{}
		for _, jsonStr := range user_properties_updateList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse user_properties_update item: %s", err))
				return
			}
			user_properties_updateObjs = append(user_properties_updateObjs, obj)
		}
		params["user_properties_update"] = user_properties_updateObjs
	}
	if !data.UserPropertiesRemove.IsNull() && !data.UserPropertiesRemove.IsUnknown() {
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
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	// Read back to populate computed fields
	id := data.ID.ValueString()
	result, err = r.client.Call("pool.snapshot.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back pool_snapshot: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["dataset"]; ok {
		switch val := v.(type) {
		case string:
			data.Dataset = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Dataset = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["properties"]; ok {
		switch val := v.(type) {
		case string:
			data.Properties = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Properties = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Properties = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["snapshot_name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Recursive.IsUnknown() {
		data.Recursive = types.BoolNull()
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.VmwareSync.IsUnknown() {
		data.VmwareSync = types.BoolNull()
	}
	if data.Properties.IsUnknown() {
		data.Properties = types.StringNull()
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.NamingSchema.IsUnknown() {
		data.NamingSchema = types.StringNull()
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
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool_snapshot: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["dataset"]; ok {
		switch val := v.(type) {
		case string:
			data.Dataset = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Dataset = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["properties"]; ok {
		switch val := v.(type) {
		case string:
			data.Properties = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Properties = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Properties = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["snapshot_name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Recursive.IsUnknown() {
		data.Recursive = types.BoolNull()
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.VmwareSync.IsUnknown() {
		data.VmwareSync = types.BoolNull()
	}
	if data.Properties.IsUnknown() {
		data.Properties = types.StringNull()
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.NamingSchema.IsUnknown() {
		data.NamingSchema = types.StringNull()
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
	if !data.UserPropertiesUpdate.IsNull() && !data.UserPropertiesUpdate.IsUnknown() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		var user_properties_updateObjs []map[string]interface{}
		for _, jsonStr := range user_properties_updateList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse user_properties_update item: %s", err))
				return
			}
			user_properties_updateObjs = append(user_properties_updateObjs, obj)
		}
		params["user_properties_update"] = user_properties_updateObjs
	}
	if !data.UserPropertiesRemove.IsNull() && !data.UserPropertiesRemove.IsUnknown() {
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

	// Read back to populate computed fields
	result, readErr := r.client.Call("pool.snapshot.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back pool_snapshot: %s", readErr))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["dataset"]; ok {
		switch val := v.(type) {
		case string:
			data.Dataset = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Dataset = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["properties"]; ok {
		switch val := v.(type) {
		case string:
			data.Properties = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Properties = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Properties = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["snapshot_name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Recursive.IsUnknown() {
		data.Recursive = types.BoolNull()
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.VmwareSync.IsUnknown() {
		data.VmwareSync = types.BoolNull()
	}
	if data.Properties.IsUnknown() {
		data.Properties = types.StringNull()
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.NamingSchema.IsUnknown() {
		data.NamingSchema = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := []interface{}{data.ID.ValueString(), map[string]interface{}{}}

	_, err := r.client.Call("pool.snapshot.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool_snapshot: %s", err))
		return
	}
}
