package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpdateRunResource struct {
	client *client.Client
}

type UpdateRunResourceModel struct {
	DatasetName types.String `tfsdk:"dataset_name"`
	Resume      types.Bool   `tfsdk:"resume"`
	Train       types.String `tfsdk:"train"`
	Version     types.String `tfsdk:"version"`
	Reboot      types.Bool   `tfsdk:"reboot"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewUpdateRunResource() resource.Resource {
	return &UpdateRunResource{}
}

func (r *UpdateRunResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_update_run"
}

func (r *UpdateRunResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Downloads (if not already in cache) and apply an update.",
		Attributes: map[string]schema.Attribute{
			"dataset_name": schema.StringAttribute{Optional: true, MarkdownDescription: "Name of the ZFS dataset to use for the new boot environment. `null` for automatic naming.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"resume":       schema.BoolAttribute{Optional: true, MarkdownDescription: "Should be set to `true` if a previous call to this method returned a `CallError` with `errno=EAGAIN` meaning     that an upgrade can be performed with a warning and that warning is accepted. In that c", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
			"train":        schema.StringAttribute{Optional: true, MarkdownDescription: "Specifies the train from which to download the update. If both `train` and `version` are `null``, the most     recent version that matches the currently selected update profile is used.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"version":      schema.StringAttribute{Optional: true, MarkdownDescription: "Specific version to update to. `null` to use the latest version from the specified train.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"reboot":       schema.BoolAttribute{Optional: true, MarkdownDescription: "Whether to automatically reboot the system after applying the update.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
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

func (r *UpdateRunResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UpdateRunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UpdateRunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := map[string]interface{}{}
	if !data.DatasetName.IsNull() && !data.DatasetName.IsUnknown() {
		params["dataset_name"] = data.DatasetName.ValueString()
	}
	if !data.Resume.IsNull() && !data.Resume.IsUnknown() {
		params["resume"] = data.Resume.ValueBool()
	}
	if !data.Train.IsNull() && !data.Train.IsUnknown() {
		params["train"] = data.Train.ValueString()
	}
	if !data.Version.IsNull() && !data.Version.IsUnknown() {
		params["version"] = data.Version.ValueString()
	}
	if !data.Reboot.IsNull() && !data.Reboot.IsUnknown() {
		params["reboot"] = data.Reboot.ValueBool()
	}
	paramsArr := []interface{}{params}

	// Execute action
	result, err := r.client.Call("update.run", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute update.run: %s", err.Error()))
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
	data.ActionID = types.StringValue(fmt.Sprintf("update.run-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UpdateRunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data UpdateRunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *UpdateRunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *UpdateRunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
