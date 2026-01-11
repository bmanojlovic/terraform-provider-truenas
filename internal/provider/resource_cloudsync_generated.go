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

type CloudsyncResource struct {
	client *client.Client
}

type CloudsyncResourceModel struct {
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
	Bwlimit types.List `tfsdk:"bwlimit"`
	Transfers types.Int64 `tfsdk:"transfers"`
	Direction types.String `tfsdk:"direction"`
	TransferMode types.String `tfsdk:"transfer_mode"`
	Encryption types.Bool `tfsdk:"encryption"`
	FilenameEncryption types.Bool `tfsdk:"filename_encryption"`
	EncryptionPassword types.String `tfsdk:"encryption_password"`
	EncryptionSalt types.String `tfsdk:"encryption_salt"`
	CreateEmptySrcDirs types.Bool `tfsdk:"create_empty_src_dirs"`
	FollowSymlinks types.Bool `tfsdk:"follow_symlinks"`
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
				Description: "ID of the cloud credential.",
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
			"bwlimit": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Schedule of bandwidth limits.",
			},
			"transfers": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum number of parallel file transfers. `null` for default.",
			},
			"direction": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Direction of the cloud sync operation.  * `PUSH`: Upload local files to cloud storage * `PULL`: Down",
			},
			"transfer_mode": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "How files are transferred between local and cloud storage.  * `SYNC`: Synchronize directories (add n",
			},
			"encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to encrypt files before uploading to cloud storage.",
			},
			"filename_encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to encrypt filenames in addition to file contents.",
			},
			"encryption_password": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Password for client-side encryption. Empty string if encryption is disabled.",
			},
			"encryption_salt": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Salt value for encryption key derivation. Empty string if encryption is disabled.",
			},
			"create_empty_src_dirs": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to create empty directories in the destination that exist in the source.",
			},
			"follow_symlinks": schema.BoolAttribute{
				Required: false,
				Optional: true,
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
	if !data.Bwlimit.IsNull() {
		var bwlimitList []string
		data.Bwlimit.ElementsAs(ctx, &bwlimitList, false)
		params["bwlimit"] = bwlimitList
	}
	if !data.Transfers.IsNull() {
		params["transfers"] = data.Transfers.ValueInt64()
	}
	if !data.Direction.IsNull() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.TransferMode.IsNull() {
		params["transfer_mode"] = data.TransferMode.ValueString()
	}
	if !data.Encryption.IsNull() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.FilenameEncryption.IsNull() {
		params["filename_encryption"] = data.FilenameEncryption.ValueBool()
	}
	if !data.EncryptionPassword.IsNull() {
		params["encryption_password"] = data.EncryptionPassword.ValueString()
	}
	if !data.EncryptionSalt.IsNull() {
		params["encryption_salt"] = data.EncryptionSalt.ValueString()
	}
	if !data.CreateEmptySrcDirs.IsNull() {
		params["create_empty_src_dirs"] = data.CreateEmptySrcDirs.ValueBool()
	}
	if !data.FollowSymlinks.IsNull() {
		params["follow_symlinks"] = data.FollowSymlinks.ValueBool()
	}

	result, err := r.client.Call("cloudsync.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create cloudsync: %s", err))
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

func (r *CloudsyncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudsyncResourceModel
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

	result, err := r.client.Call("cloudsync.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read cloudsync: %s", err))
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
		if v, ok := resultMap["bwlimit"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Bwlimit, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["transfers"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Transfers = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["direction"]; ok && v != nil {
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["transfer_mode"]; ok && v != nil {
			data.TransferMode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Encryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["filename_encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.FilenameEncryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_password"]; ok && v != nil {
			data.EncryptionPassword = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_salt"]; ok && v != nil {
			data.EncryptionSalt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["create_empty_src_dirs"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.CreateEmptySrcDirs = types.BoolValue(bv) }
		}
		if v, ok := resultMap["follow_symlinks"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.FollowSymlinks = types.BoolValue(bv) }
		}
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
	if !data.Bwlimit.IsNull() {
		var bwlimitList []string
		data.Bwlimit.ElementsAs(ctx, &bwlimitList, false)
		params["bwlimit"] = bwlimitList
	}
	if !data.Transfers.IsNull() {
		params["transfers"] = data.Transfers.ValueInt64()
	}
	if !data.Direction.IsNull() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.TransferMode.IsNull() {
		params["transfer_mode"] = data.TransferMode.ValueString()
	}
	if !data.Encryption.IsNull() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.FilenameEncryption.IsNull() {
		params["filename_encryption"] = data.FilenameEncryption.ValueBool()
	}
	if !data.EncryptionPassword.IsNull() {
		params["encryption_password"] = data.EncryptionPassword.ValueString()
	}
	if !data.EncryptionSalt.IsNull() {
		params["encryption_salt"] = data.EncryptionSalt.ValueString()
	}
	if !data.CreateEmptySrcDirs.IsNull() {
		params["create_empty_src_dirs"] = data.CreateEmptySrcDirs.ValueBool()
	}
	if !data.FollowSymlinks.IsNull() {
		params["follow_symlinks"] = data.FollowSymlinks.ValueBool()
	}

	_, err = r.client.Call("cloudsync.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update cloudsync: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudsyncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudsyncResourceModel
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

	_, err = r.client.Call("cloudsync.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete cloudsync: %s", err))
		return
	}
}
