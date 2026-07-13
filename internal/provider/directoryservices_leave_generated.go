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

type DirectoryservicesLeaveResource struct {
	client *client.Client
}

type DirectoryservicesLeaveResourceModel struct {
	CredentialCredentialType types.String `tfsdk:"credential_credential_type"`
	CredentialUsername       types.String `tfsdk:"credential_username"`
	CredentialPassword       types.String `tfsdk:"credential_password"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewDirectoryservicesLeaveResource() resource.Resource {
	return &DirectoryservicesLeaveResource{}
}

func (r *DirectoryservicesLeaveResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directoryservices_leave"
}

func (r *DirectoryservicesLeaveResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Leave an Active Directory or IPA domain. Calling this endpoint when the directory services status is `HEALTHY` will cause TrueNAS to remove its account from the domain and then reset the local directo",
		Attributes: map[string]schema.Attribute{
			"credential_credential_type": schema.StringAttribute{Required: true, MarkdownDescription: "Credential type identifier for Kerberos user authentication.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"credential_username":        schema.StringAttribute{Required: true, MarkdownDescription: "Username of the account to use to create a kerberos ticket for authentication to directory services. This     account must exist on the domain controller. ", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"credential_password":        schema.StringAttribute{Required: true, MarkdownDescription: "The password for the user account that will obtain the kerberos ticket. ", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
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

func (r *DirectoryservicesLeaveResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DirectoryservicesLeaveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DirectoryservicesLeaveResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := map[string]interface{}{}
	credentialOpts := map[string]interface{}{}
	if !data.CredentialCredentialType.IsNull() {
		credentialOpts["credential_type"] = data.CredentialCredentialType.ValueString()
	}
	if !data.CredentialUsername.IsNull() {
		credentialOpts["username"] = data.CredentialUsername.ValueString()
	}
	if !data.CredentialPassword.IsNull() {
		credentialOpts["password"] = data.CredentialPassword.ValueString()
	}
	if len(credentialOpts) > 0 {
		params["credential"] = credentialOpts
	}
	paramsArr := []interface{}{params}

	// Execute action
	result, err := r.client.Call("directoryservices.leave", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute directoryservices.leave: %s", err.Error()))
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
	data.ActionID = types.StringValue(fmt.Sprintf("directoryservices.leave-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DirectoryservicesLeaveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data DirectoryservicesLeaveResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *DirectoryservicesLeaveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *DirectoryservicesLeaveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
