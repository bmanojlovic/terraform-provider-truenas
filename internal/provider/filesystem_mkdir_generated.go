package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilesystemMkdirResource struct {
	client *client.Client
}

type FilesystemMkdirResourceModel struct {
	Path types.String `tfsdk:"path"`
	OptionsMode types.String `tfsdk:"options_mode"`
	OptionsRaiseChmodError types.Bool `tfsdk:"options_raise_chmod_error"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewFilesystemMkdirResource() resource.Resource {
	return &FilesystemMkdirResource{}
}

func (r *FilesystemMkdirResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_mkdir"
}

func (r *FilesystemMkdirResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a directory at the specified path.  The following options are supported:  `mode` - specify the permissions to set on the new directory (0o755 is default). `raise_chmod_error` - choose whether t",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{Required: true, MarkdownDescription: "Path where the new directory should be created.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"options_mode": schema.StringAttribute{Optional: true, MarkdownDescription: "Unix permissions for the new directory."},
			"options_raise_chmod_error": schema.BoolAttribute{Optional: true, MarkdownDescription: "Whether to raise an error if chmod fails."},
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

func (r *FilesystemMkdirResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FilesystemMkdirResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FilesystemMkdirResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := map[string]interface{}{}
	params["path"] = data.Path.ValueString()
	optionsOpts := map[string]interface{}{}
	if !data.OptionsMode.IsNull() { optionsOpts["mode"] = data.OptionsMode.ValueString() }
	if !data.OptionsRaiseChmodError.IsNull() { optionsOpts["raise_chmod_error"] = data.OptionsRaiseChmodError.ValueBool() }
	if len(optionsOpts) > 0 { params["options"] = optionsOpts }
	paramsArr := []interface{}{params}

	// Execute action
	result, err := r.client.Call("filesystem.mkdir", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute filesystem.mkdir: %s", err.Error()))
		return
	}

	// Check if result is a job ID
	if jobID, ok := result.(float64); ok && false {
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
	data.ActionID = types.StringValue(fmt.Sprintf("filesystem.mkdir-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemMkdirResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data FilesystemMkdirResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *FilesystemMkdirResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FilesystemMkdirResourceModel
	var oldData FilesystemMkdirResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &oldData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update parameters
	updateParams := map[string]interface{}{}
	if !data.Path.IsNull() { updateParams["path"] = data.Path.ValueString() }
	if !data.OptionsMode.IsNull() { updateParams["mode"] = data.OptionsMode.ValueString() }
	updateParamsArr := []interface{}{updateParams}

	// Execute update
	_, err := r.client.Call("filesystem.setperm", updateParamsArr)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", fmt.Sprintf("Failed to execute filesystem.setperm: %s", err.Error()))
		return
	}

	// Preserve computed fields from old state
	data.ActionID = oldData.ActionID
	data.JobID = oldData.JobID
	data.State = oldData.State
	data.Progress = oldData.Progress
	data.Result = oldData.Result
	data.Error = oldData.Error

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemMkdirResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
