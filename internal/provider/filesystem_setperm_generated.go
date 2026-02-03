package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilesystemSetpermResource struct {
	client *client.Client
}

type FilesystemSetpermResourceModel struct {
	Path types.String `tfsdk:"path"`
	Uid types.Int64 `tfsdk:"uid"`
	User types.String `tfsdk:"user"`
	Gid types.Int64 `tfsdk:"gid"`
	Group types.String `tfsdk:"group"`
	Mode types.String `tfsdk:"mode"`
	OptionsRecursive types.Bool `tfsdk:"options_recursive"`
	OptionsTraverse types.Bool `tfsdk:"options_traverse"`
	OptionsStripacl types.Bool `tfsdk:"options_stripacl"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewFilesystemSetpermResource() resource.Resource {
	return &FilesystemSetpermResource{}
}

func (r *FilesystemSetpermResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_setperm"
}

func (r *FilesystemSetpermResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Set unix permissions on given `path`.  If `mode` is specified then the mode will be applied to the path and files and subdirectories depending on which `options` are selected. Mode should be formatted",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{Required: true, MarkdownDescription: "Filesystem path to modify."},
			"uid": schema.Int64Attribute{Optional: true, MarkdownDescription: "Numeric user ID to set as owner. `null` to leave unchanged."},
			"user": schema.StringAttribute{Optional: true, MarkdownDescription: "Username to set as owner. `null` to leave unchanged."},
			"gid": schema.Int64Attribute{Optional: true, MarkdownDescription: "Numeric group ID to set as group owner. `null` to leave unchanged."},
			"group": schema.StringAttribute{Optional: true, MarkdownDescription: "Group name to set as group owner. `null` to leave unchanged."},
			"mode": schema.StringAttribute{Optional: true, MarkdownDescription: "Unix permissions to set (octal format). `null` to leave unchanged."},
			"options_recursive": schema.BoolAttribute{Optional: true, MarkdownDescription: "Whether to apply the operation recursively to subdirectories."},
			"options_traverse": schema.BoolAttribute{Optional: true, MarkdownDescription: "If set do not limit to single dataset / filesystem."},
			"options_stripacl": schema.BoolAttribute{Optional: true, MarkdownDescription: "Whether to remove existing Access Control Lists when setting permissions."},
			"action_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Action execution identifier",
			},
			"job_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Background job ID (if applicable)",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Job state: SUCCESS, FAILED, or RUNNING",
			},
			"progress": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "Job progress percentage (0-100)",
			},
			"result": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Action result data",
			},
			"error": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Error message if action failed",
			},
		},
	}
}

func (r *FilesystemSetpermResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FilesystemSetpermResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FilesystemSetpermResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := map[string]interface{}{}
	params["path"] = data.Path.ValueString()
	if !data.Uid.IsNull() { params["uid"] = data.Uid.ValueInt64() }
	if !data.User.IsNull() { params["user"] = data.User.ValueString() }
	if !data.Gid.IsNull() { params["gid"] = data.Gid.ValueInt64() }
	if !data.Group.IsNull() { params["group"] = data.Group.ValueString() }
	if !data.Mode.IsNull() { params["mode"] = data.Mode.ValueString() }
	optionsOpts := map[string]interface{}{}
	if !data.OptionsRecursive.IsNull() { optionsOpts["recursive"] = data.OptionsRecursive.ValueBool() }
	if !data.OptionsTraverse.IsNull() { optionsOpts["traverse"] = data.OptionsTraverse.ValueBool() }
	if !data.OptionsStripacl.IsNull() { optionsOpts["stripacl"] = data.OptionsStripacl.ValueBool() }
	if len(optionsOpts) > 0 { params["options"] = optionsOpts }
	paramsArr := []interface{}{params}

	// Execute action
	result, err := r.client.Call("filesystem.setperm", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute filesystem.setperm: %s", err.Error()))
		return
	}

	// Check if result is a job ID
	if jobID, ok := result.(float64); ok && true {
		// Background job - wait for completion
		data.JobID = types.Int64Value(int64(jobID))
		
		jobResult, err := r.client.WaitForJob(int(jobID), 30*time.Minute)
		if err != nil {
			data.State = types.StringValue("FAILED")
			data.Error = types.StringValue(err.Error())
			resp.Diagnostics.AddError("Job Failed", err.Error())
		} else {
			data.State = types.StringValue(jobResult.State)
			data.Progress = types.Float64Value(jobResult.Progress)
			data.Result = types.StringValue(fmt.Sprintf("%v", jobResult.Result))
			if jobResult.Error != "" {
				data.Error = types.StringValue(jobResult.Error)
			} else {
				data.Error = types.StringValue("")
			}
		}
	} else {
		// Immediate result
		data.JobID = types.Int64Value(0)
		data.State = types.StringValue("SUCCESS")
		data.Progress = types.Float64Value(100.0)
		data.Result = types.StringValue(fmt.Sprintf("%v", result))
		data.Error = types.StringValue("")
	}

	// Generate ID from timestamp
	data.ActionID = types.StringValue(fmt.Sprintf("filesystem.setperm-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemSetpermResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data FilesystemSetpermResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *FilesystemSetpermResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *FilesystemSetpermResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
