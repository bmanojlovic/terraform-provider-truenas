package provider

import (
	"context"
	"fmt"
	"strconv"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type CloudBackupResource struct {
	client *client.Client
}

type CloudBackupResourceModel struct {
	ID types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Path types.String `tfsdk:"path"`
	Credentials types.Int64 `tfsdk:"credentials"`
	Attributes types.String `tfsdk:"attributes"`
	Schedule types.String `tfsdk:"schedule"`
	PreScript types.String `tfsdk:"pre_script"`
	PostScript types.String `tfsdk:"post_script"`
	Snapshot types.Bool `tfsdk:"snapshot"`
	Include types.List `tfsdk:"include"`
	Exclude types.List `tfsdk:"exclude"`
	Args types.String `tfsdk:"args"`
	Enabled types.Bool `tfsdk:"enabled"`
	Password types.String `tfsdk:"password"`
	KeepLast types.Int64 `tfsdk:"keep_last"`
	TransferSetting types.String `tfsdk:"transfer_setting"`
	AbsolutePaths types.Bool `tfsdk:"absolute_paths"`
	CachePath types.String `tfsdk:"cache_path"`
	RateLimit types.Int64 `tfsdk:"rate_limit"`
}

func NewCloudBackupResource() resource.Resource {
	return &CloudBackupResource{}
}

func (r *CloudBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_backup"
}

func (r *CloudBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CloudBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new cloud backup task",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The name of the task to display in the UI.",
			},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "The local path to back up beginning with `/mnt` or `/dev/zvol`.",
			},
			"credentials": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the cloud credential to use for each backup.",
			},
			"attributes": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Additional information for each backup, e.g. bucket name.",
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Cron schedule dictating when the task should run.",
			},
			"pre_script": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "A Bash script to run immediately before every backup.",
			},
			"post_script": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "A Bash script to run immediately after every backup if it succeeds.",
			},
			"snapshot": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to create a temporary snapshot of the dataset before every backup.",
			},
			"include": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Paths to pass to `restic backup --include`.",
			},
			"exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Paths to pass to `restic backup --exclude`.",
			},
			"args": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "(Slated for removal).",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Can enable/disable the task.",
			},
			"password": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Password for the remote repository.",
			},
			"keep_last": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "How many of the most recent backup snapshots to keep after each backup.",
			},
			"transfer_setting": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "* DEFAULT:     * pack size given by `$RESTIC_PACK_SIZE` (default 16 MiB)     * read concurrency give",
			},
			"absolute_paths": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Preserve absolute paths in each backup (cannot be set when `snapshot=True`).",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"cache_path": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Cache path. If not set, performance may degrade.",
			},
			"rate_limit": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum upload/download rate in KiB/s. Passed to `restic --limit-upload` on `cloud_backup.sync` and ",
			},
		},
	}
}

func (r *CloudBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Credentials.IsNull() {
		params["credentials"] = data.Credentials.ValueInt64()
	}
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Schedule.IsNull() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.PreScript.IsNull() {
		params["pre_script"] = data.PreScript.ValueString()
	}
	if !data.PostScript.IsNull() {
		params["post_script"] = data.PostScript.ValueString()
	}
	if !data.Snapshot.IsNull() {
		params["snapshot"] = data.Snapshot.ValueBool()
	}
	if !data.Include.IsNull() {
		var includeList []string
		data.Include.ElementsAs(ctx, &includeList, false)
		params["include"] = includeList
	}
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Args.IsNull() {
		params["args"] = data.Args.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}
	if !data.KeepLast.IsNull() {
		params["keep_last"] = data.KeepLast.ValueInt64()
	}
	if !data.TransferSetting.IsNull() {
		params["transfer_setting"] = data.TransferSetting.ValueString()
	}
	if !data.AbsolutePaths.IsNull() {
		params["absolute_paths"] = data.AbsolutePaths.ValueBool()
	}
	if !data.CachePath.IsNull() {
		params["cache_path"] = data.CachePath.ValueString()
	}
	if !data.RateLimit.IsNull() {
		params["rate_limit"] = data.RateLimit.ValueInt64()
	}

	result, err := r.client.Call("cloud_backup.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create cloud_backup: %s", err))
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

func (r *CloudBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudBackupResourceModel
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

	result, err := r.client.Call("cloud_backup.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read cloud_backup: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["credentials"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Credentials = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			data.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pre_script"]; ok && v != nil {
			data.PreScript = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["post_script"]; ok && v != nil {
			data.PostScript = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapshot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Snapshot = types.BoolValue(bv) }
		}
		if v, ok := resultMap["include"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Include, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["exclude"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Exclude, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["args"]; ok && v != nil {
			data.Args = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["keep_last"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.KeepLast = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["transfer_setting"]; ok && v != nil {
			data.TransferSetting = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["absolute_paths"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AbsolutePaths = types.BoolValue(bv) }
		}
		if v, ok := resultMap["cache_path"]; ok && v != nil {
			data.CachePath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["rate_limit"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.RateLimit = types.Int64Value(int64(fv)) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CloudBackupResourceModel
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
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Credentials.IsNull() {
		params["credentials"] = data.Credentials.ValueInt64()
	}
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Schedule.IsNull() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.PreScript.IsNull() {
		params["pre_script"] = data.PreScript.ValueString()
	}
	if !data.PostScript.IsNull() {
		params["post_script"] = data.PostScript.ValueString()
	}
	if !data.Snapshot.IsNull() {
		params["snapshot"] = data.Snapshot.ValueBool()
	}
	if !data.Include.IsNull() {
		var includeList []string
		data.Include.ElementsAs(ctx, &includeList, false)
		params["include"] = includeList
	}
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Args.IsNull() {
		params["args"] = data.Args.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}
	if !data.KeepLast.IsNull() {
		params["keep_last"] = data.KeepLast.ValueInt64()
	}
	if !data.TransferSetting.IsNull() {
		params["transfer_setting"] = data.TransferSetting.ValueString()
	}
	if !data.CachePath.IsNull() {
		params["cache_path"] = data.CachePath.ValueString()
	}
	if !data.RateLimit.IsNull() {
		params["rate_limit"] = data.RateLimit.ValueInt64()
	}

	_, err = r.client.Call("cloud_backup.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update cloud_backup: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudBackupResourceModel
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

	_, err = r.client.Call("cloud_backup.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete cloud_backup: %s", err))
		return
	}
}
