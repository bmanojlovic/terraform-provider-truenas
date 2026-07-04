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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type CloudsyncResource struct {
	client *client.Client
}

type CloudsyncResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Description        types.String `tfsdk:"description"`
	Path               types.String `tfsdk:"path"`
	Credentials        types.Int64  `tfsdk:"credentials"`
	Attributes         types.String `tfsdk:"attributes"`
	Schedule           types.String `tfsdk:"schedule"`
	PreScript          types.String `tfsdk:"pre_script"`
	PostScript         types.String `tfsdk:"post_script"`
	Snapshot           types.Bool   `tfsdk:"snapshot"`
	Include            types.List   `tfsdk:"include"`
	Exclude            types.List   `tfsdk:"exclude"`
	Args               types.String `tfsdk:"args"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	Bwlimit            types.List   `tfsdk:"bwlimit"`
	Transfers          types.Int64  `tfsdk:"transfers"`
	Direction          types.String `tfsdk:"direction"`
	TransferMode       types.String `tfsdk:"transfer_mode"`
	Encryption         types.Bool   `tfsdk:"encryption"`
	FilenameEncryption types.Bool   `tfsdk:"filename_encryption"`
	EncryptionPassword types.String `tfsdk:"encryption_password"`
	EncryptionSalt     types.String `tfsdk:"encryption_salt"`
	CreateEmptySrcDirs types.Bool   `tfsdk:"create_empty_src_dirs"`
	FollowSymlinks     types.Bool   `tfsdk:"follow_symlinks"`
}

func NewCloudsyncResource() resource.Resource {
	return &CloudsyncResource{}
}

func (r *CloudsyncResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudsync"
}

func (r *CloudsyncResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CloudsyncResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a new cloud_sync entry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the task to display in the UI.",
			},
			"path": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "The local path to back up beginning with `/mnt` or `/dev/zvol`.",
			},
			"credentials": schema.Int64Attribute{
				Required:    true,
				Optional:    false,
				Description: "ID of the cloud credential.",
			},
			"attributes": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Additional information for each backup, e.g. bucket name.",
			},
			"schedule": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cron schedule dictating when the task should run.",
			},
			"pre_script": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A Bash script to run immediately before every backup.",
			},
			"post_script": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A Bash script to run immediately after every backup if it succeeds.",
			},
			"snapshot": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to create a temporary snapshot of the dataset before every backup.",
			},
			"include": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Paths to pass to `restic backup --include`.",
			},
			"exclude": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Paths to pass to `restic backup --exclude`.",
			},
			"args": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "(Slated for removal).",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Can enable/disable the task.",
			},
			"bwlimit": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Schedule of bandwidth limits.",
			},
			"transfers": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of parallel file transfers. `null` for default.",
			},
			"direction": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Direction of the cloud sync operation.  * `PUSH`: Upload local files to cloud storage * `PULL`: Down",
			},
			"transfer_mode": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "How files are transferred between local and cloud storage.  * `SYNC`: Synchronize directories (add n",
			},
			"encryption": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to encrypt files before uploading to cloud storage.",
			},
			"filename_encryption": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to encrypt filenames in addition to file contents.",
			},
			"encryption_password": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Password for client-side encryption. Empty string if encryption is disabled.",
			},
			"encryption_salt": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Salt value for encryption key derivation. Empty string if encryption is disabled.",
			},
			"create_empty_src_dirs": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to create empty directories in the destination that exist in the source.",
			},
			"follow_symlinks": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to follow symbolic links and sync the files they point to.",
			},
		},
	}
}

func (r *CloudsyncResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudsyncResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Credentials.IsNull() && !data.Credentials.IsUnknown() {
		params["credentials"] = data.Credentials.ValueInt64()
	}
	if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.PreScript.IsNull() && !data.PreScript.IsUnknown() {
		params["pre_script"] = data.PreScript.ValueString()
	}
	if !data.PostScript.IsNull() && !data.PostScript.IsUnknown() {
		params["post_script"] = data.PostScript.ValueString()
	}
	if !data.Snapshot.IsNull() && !data.Snapshot.IsUnknown() {
		params["snapshot"] = data.Snapshot.ValueBool()
	}
	if !data.Include.IsNull() && !data.Include.IsUnknown() {
		var includeList []string
		data.Include.ElementsAs(ctx, &includeList, false)
		params["include"] = includeList
	}
	if !data.Exclude.IsNull() && !data.Exclude.IsUnknown() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Args.IsNull() && !data.Args.IsUnknown() {
		params["args"] = data.Args.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Bwlimit.IsNull() && !data.Bwlimit.IsUnknown() {
		var bwlimitList []string
		data.Bwlimit.ElementsAs(ctx, &bwlimitList, false)
		var bwlimitObjs []map[string]interface{}
		for _, jsonStr := range bwlimitList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse bwlimit item: %s", err))
				return
			}
			bwlimitObjs = append(bwlimitObjs, obj)
		}
		params["bwlimit"] = bwlimitObjs
	}
	if !data.Transfers.IsNull() && !data.Transfers.IsUnknown() {
		params["transfers"] = data.Transfers.ValueInt64()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.TransferMode.IsNull() && !data.TransferMode.IsUnknown() {
		params["transfer_mode"] = data.TransferMode.ValueString()
	}
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.FilenameEncryption.IsNull() && !data.FilenameEncryption.IsUnknown() {
		params["filename_encryption"] = data.FilenameEncryption.ValueBool()
	}
	if !data.EncryptionPassword.IsNull() && !data.EncryptionPassword.IsUnknown() {
		params["encryption_password"] = data.EncryptionPassword.ValueString()
	}
	if !data.EncryptionSalt.IsNull() && !data.EncryptionSalt.IsUnknown() {
		params["encryption_salt"] = data.EncryptionSalt.ValueString()
	}
	if !data.CreateEmptySrcDirs.IsNull() && !data.CreateEmptySrcDirs.IsUnknown() {
		params["create_empty_src_dirs"] = data.CreateEmptySrcDirs.ValueBool()
	}
	if !data.FollowSymlinks.IsNull() && !data.FollowSymlinks.IsUnknown() {
		params["follow_symlinks"] = data.FollowSymlinks.ValueBool()
	}

	result, err := r.client.Call("cloudsync.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create cloudsync: %s", err))
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
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("cloudsync.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back cloudsync: %s", err))
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
	if v, ok := resultMap["path"]; ok {
		switch val := v.(type) {
		case string:
			data.Path = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["credentials"]; ok {
		switch val := v.(type) {
		case float64:
			data.Credentials = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Credentials = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["attributes"]; ok {
		switch val := v.(type) {
		case string:
			data.Attributes = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Attributes = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfers"]; ok {
		if v == nil {
			data.Transfers = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.Transfers = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.Transfers = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["direction"]; ok {
		switch val := v.(type) {
		case string:
			data.Direction = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Direction = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfer_mode"]; ok {
		switch val := v.(type) {
		case string:
			data.TransferMode = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.TransferMode = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.TransferMode = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Schedule.IsUnknown() {
		data.Schedule = types.StringNull()
	}
	if data.PreScript.IsUnknown() {
		data.PreScript = types.StringNull()
	}
	if data.PostScript.IsUnknown() {
		data.PostScript = types.StringNull()
	}
	if data.Snapshot.IsUnknown() {
		data.Snapshot = types.BoolNull()
	}
	if data.Include.IsUnknown() {
		data.Include, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Args.IsUnknown() {
		data.Args = types.StringNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}
	if data.Bwlimit.IsUnknown() {
		data.Bwlimit, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Transfers.IsUnknown() {
		data.Transfers = types.Int64Null()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.FilenameEncryption.IsUnknown() {
		data.FilenameEncryption = types.BoolNull()
	}
	if data.EncryptionPassword.IsUnknown() {
		data.EncryptionPassword = types.StringNull()
	}
	if data.EncryptionSalt.IsUnknown() {
		data.EncryptionSalt = types.StringNull()
	}
	if data.CreateEmptySrcDirs.IsUnknown() {
		data.CreateEmptySrcDirs = types.BoolNull()
	}
	if data.FollowSymlinks.IsUnknown() {
		data.FollowSymlinks = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudsyncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("cloudsync.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read cloudsync: %s", err))
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
	if v, ok := resultMap["path"]; ok {
		switch val := v.(type) {
		case string:
			data.Path = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["credentials"]; ok {
		switch val := v.(type) {
		case float64:
			data.Credentials = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Credentials = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["attributes"]; ok {
		switch val := v.(type) {
		case string:
			data.Attributes = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Attributes = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfers"]; ok {
		if v == nil {
			data.Transfers = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.Transfers = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.Transfers = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["direction"]; ok {
		switch val := v.(type) {
		case string:
			data.Direction = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Direction = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfer_mode"]; ok {
		switch val := v.(type) {
		case string:
			data.TransferMode = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.TransferMode = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.TransferMode = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Schedule.IsUnknown() {
		data.Schedule = types.StringNull()
	}
	if data.PreScript.IsUnknown() {
		data.PreScript = types.StringNull()
	}
	if data.PostScript.IsUnknown() {
		data.PostScript = types.StringNull()
	}
	if data.Snapshot.IsUnknown() {
		data.Snapshot = types.BoolNull()
	}
	if data.Include.IsUnknown() {
		data.Include, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Args.IsUnknown() {
		data.Args = types.StringNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}
	if data.Bwlimit.IsUnknown() {
		data.Bwlimit, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Transfers.IsUnknown() {
		data.Transfers = types.Int64Null()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.FilenameEncryption.IsUnknown() {
		data.FilenameEncryption = types.BoolNull()
	}
	if data.EncryptionPassword.IsUnknown() {
		data.EncryptionPassword = types.StringNull()
	}
	if data.EncryptionSalt.IsUnknown() {
		data.EncryptionSalt = types.StringNull()
	}
	if data.CreateEmptySrcDirs.IsUnknown() {
		data.CreateEmptySrcDirs = types.BoolNull()
	}
	if data.FollowSymlinks.IsUnknown() {
		data.FollowSymlinks = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudsyncResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CloudsyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Credentials.IsNull() && !data.Credentials.IsUnknown() {
		params["credentials"] = data.Credentials.ValueInt64()
	}
	if !data.Attributes.IsNull() && !data.Attributes.IsUnknown() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.PreScript.IsNull() && !data.PreScript.IsUnknown() {
		params["pre_script"] = data.PreScript.ValueString()
	}
	if !data.PostScript.IsNull() && !data.PostScript.IsUnknown() {
		params["post_script"] = data.PostScript.ValueString()
	}
	if !data.Snapshot.IsNull() && !data.Snapshot.IsUnknown() {
		params["snapshot"] = data.Snapshot.ValueBool()
	}
	if !data.Include.IsNull() && !data.Include.IsUnknown() {
		var includeList []string
		data.Include.ElementsAs(ctx, &includeList, false)
		params["include"] = includeList
	}
	if !data.Exclude.IsNull() && !data.Exclude.IsUnknown() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Args.IsNull() && !data.Args.IsUnknown() {
		params["args"] = data.Args.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Bwlimit.IsNull() && !data.Bwlimit.IsUnknown() {
		var bwlimitList []string
		data.Bwlimit.ElementsAs(ctx, &bwlimitList, false)
		var bwlimitObjs []map[string]interface{}
		for _, jsonStr := range bwlimitList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse bwlimit item: %s", err))
				return
			}
			bwlimitObjs = append(bwlimitObjs, obj)
		}
		params["bwlimit"] = bwlimitObjs
	}
	if !data.Transfers.IsNull() && !data.Transfers.IsUnknown() {
		params["transfers"] = data.Transfers.ValueInt64()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.TransferMode.IsNull() && !data.TransferMode.IsUnknown() {
		params["transfer_mode"] = data.TransferMode.ValueString()
	}
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.FilenameEncryption.IsNull() && !data.FilenameEncryption.IsUnknown() {
		params["filename_encryption"] = data.FilenameEncryption.ValueBool()
	}
	if !data.EncryptionPassword.IsNull() && !data.EncryptionPassword.IsUnknown() {
		params["encryption_password"] = data.EncryptionPassword.ValueString()
	}
	if !data.EncryptionSalt.IsNull() && !data.EncryptionSalt.IsUnknown() {
		params["encryption_salt"] = data.EncryptionSalt.ValueString()
	}
	if !data.CreateEmptySrcDirs.IsNull() && !data.CreateEmptySrcDirs.IsUnknown() {
		params["create_empty_src_dirs"] = data.CreateEmptySrcDirs.ValueBool()
	}
	if !data.FollowSymlinks.IsNull() && !data.FollowSymlinks.IsUnknown() {
		params["follow_symlinks"] = data.FollowSymlinks.ValueBool()
	}

	_, err = r.client.Call("cloudsync.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update cloudsync: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("cloudsync.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back cloudsync: %s", readErr))
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
	if v, ok := resultMap["path"]; ok {
		switch val := v.(type) {
		case string:
			data.Path = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["credentials"]; ok {
		switch val := v.(type) {
		case float64:
			data.Credentials = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Credentials = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["attributes"]; ok {
		switch val := v.(type) {
		case string:
			data.Attributes = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Attributes = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfers"]; ok {
		if v == nil {
			data.Transfers = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.Transfers = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.Transfers = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["direction"]; ok {
		switch val := v.(type) {
		case string:
			data.Direction = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Direction = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfer_mode"]; ok {
		switch val := v.(type) {
		case string:
			data.TransferMode = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.TransferMode = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.TransferMode = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Schedule.IsUnknown() {
		data.Schedule = types.StringNull()
	}
	if data.PreScript.IsUnknown() {
		data.PreScript = types.StringNull()
	}
	if data.PostScript.IsUnknown() {
		data.PostScript = types.StringNull()
	}
	if data.Snapshot.IsUnknown() {
		data.Snapshot = types.BoolNull()
	}
	if data.Include.IsUnknown() {
		data.Include, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Args.IsUnknown() {
		data.Args = types.StringNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}
	if data.Bwlimit.IsUnknown() {
		data.Bwlimit, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Transfers.IsUnknown() {
		data.Transfers = types.Int64Null()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.FilenameEncryption.IsUnknown() {
		data.FilenameEncryption = types.BoolNull()
	}
	if data.EncryptionPassword.IsUnknown() {
		data.EncryptionPassword = types.StringNull()
	}
	if data.EncryptionSalt.IsUnknown() {
		data.EncryptionSalt = types.StringNull()
	}
	if data.CreateEmptySrcDirs.IsUnknown() {
		data.CreateEmptySrcDirs = types.BoolNull()
	}
	if data.FollowSymlinks.IsUnknown() {
		data.FollowSymlinks = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudsyncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("cloudsync.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete cloudsync: %s", err))
		return
	}
}
