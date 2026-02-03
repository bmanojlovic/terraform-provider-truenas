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

type BootDetachResource struct {
	client *client.Client
}

type BootDetachResourceModel struct {
	Dev types.String `tfsdk:"dev"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewBootDetachResource() resource.Resource {
	return &BootDetachResource{}
}

func (r *BootDetachResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_boot_detach"
}

func (r *BootDetachResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Detach given `dev` from boot pool.",
		Attributes: map[string]schema.Attribute{
			"dev": schema.StringAttribute{Required: true, MarkdownDescription: "Device name or path to detach from the boot pool."},
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

func (r *BootDetachResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BootDetachResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BootDetachResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	paramsArr := []interface{}{}
	paramsArr = append(paramsArr, data.Dev.ValueString())

	// Execute action
	result, err := r.client.Call("boot.detach", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute boot.detach: %s", err.Error()))
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
	data.ActionID = types.StringValue(fmt.Sprintf("boot.detach-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BootDetachResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data BootDetachResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *BootDetachResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *BootDetachResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
