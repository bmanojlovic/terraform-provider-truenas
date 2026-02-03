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

type FailoverRebootOther_NodeResource struct {
	client *client.Client
}

type FailoverRebootOther_NodeResourceModel struct {
	Reason types.String `tfsdk:"reason"`
	Graceful types.Bool `tfsdk:"graceful"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewFailoverRebootOther_NodeResource() resource.Resource {
	return &FailoverRebootOther_NodeResource{}
}

func (r *FailoverRebootOther_NodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_failover_reboot_other_node"
}

func (r *FailoverRebootOther_NodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reboot the other node and wait for it to come back online.  NOTE: This makes very few checks on HA systems. You need to     know what you're doing before calling this.",
		Attributes: map[string]schema.Attribute{
			"reason": schema.StringAttribute{Optional: true, MarkdownDescription: "Reason for the system reboot."},
			"graceful": schema.BoolAttribute{Optional: true, MarkdownDescription: "If set, call `system.reboot` to gracefully reboot the other node. By default, `failover.become_passive` will be     called on the other node to forcefully reboot and simulate a failover event unless t"},
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

func (r *FailoverRebootOther_NodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FailoverRebootOther_NodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FailoverRebootOther_NodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := map[string]interface{}{}
	if !data.Reason.IsNull() { params["reason"] = data.Reason.ValueString() }
	if !data.Graceful.IsNull() { params["graceful"] = data.Graceful.ValueBool() }
	paramsArr := []interface{}{params}

	// Execute action
	result, err := r.client.Call("failover.reboot.other_node", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute failover.reboot.other_node: %s", err.Error()))
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
	data.ActionID = types.StringValue(fmt.Sprintf("failover.reboot.other_node-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FailoverRebootOther_NodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data FailoverRebootOther_NodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *FailoverRebootOther_NodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *FailoverRebootOther_NodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
