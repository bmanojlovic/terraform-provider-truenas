package provider

import (
	"context"
	"fmt"
	"time"

	"encoding/json"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ActionCloud_BackupRestoreResource struct {
	client *client.Client
}

type ActionCloud_BackupRestoreResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	SnapshotId      types.String `tfsdk:"snapshot_id"`
	Subfolder       types.String `tfsdk:"subfolder"`
	DestinationPath types.String `tfsdk:"destination_path"`
	Options         types.String `tfsdk:"options"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewActionCloud_BackupRestoreResource() resource.Resource {
	return &ActionCloud_BackupRestoreResource{}
}

func (r *ActionCloud_BackupRestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_action_cloud_backup_restore"
}

func (r *ActionCloud_BackupRestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Restore files to the directory `destination_path` from the `snapshot_id` subfolder `subfolder` created by the cloud backup job `id`.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.Int64Attribute{Required: true, MarkdownDescription: "ID of the cloud backup task.", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"snapshot_id":      schema.StringAttribute{Required: true, MarkdownDescription: "ID of the snapshot to restore.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"subfolder":        schema.StringAttribute{Required: true, MarkdownDescription: "Path within the snapshot to restore.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"destination_path": schema.StringAttribute{Required: true, MarkdownDescription: "Local path to restore to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"options":          schema.StringAttribute{Optional: true, MarkdownDescription: "Additional restore options.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
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

func (r *ActionCloud_BackupRestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ActionCloud_BackupRestoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ActionCloud_BackupRestoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := []interface{}{}
	params = append(params, data.ID.ValueInt64())
	params = append(params, data.SnapshotId.ValueString())
	params = append(params, data.Subfolder.ValueString())
	params = append(params, data.DestinationPath.ValueString())
	if !data.Options.IsNull() {
		var optionsVal interface{}
		if err := json.Unmarshal([]byte(data.Options.ValueString()), &optionsVal); err == nil {
			params = append(params, optionsVal)
		}
	}

	// Execute action
	result, err := r.client.Call("cloud_backup.restore", params)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute cloud_backup.restore: %s", err.Error()))
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
	data.ActionID = types.StringValue(fmt.Sprintf("cloud_backup.restore-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ActionCloud_BackupRestoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data ActionCloud_BackupRestoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *ActionCloud_BackupRestoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *ActionCloud_BackupRestoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
