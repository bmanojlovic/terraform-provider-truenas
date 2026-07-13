package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilesystemChownResource struct {
	client *client.Client
}

type FilesystemChownResourceModel struct {
	Path             types.String `tfsdk:"path"`
	Uid              types.Int64  `tfsdk:"uid"`
	User             types.String `tfsdk:"user"`
	Gid              types.Int64  `tfsdk:"gid"`
	Group            types.String `tfsdk:"group"`
	OptionsRecursive types.Bool   `tfsdk:"options_recursive"`
	OptionsTraverse  types.Bool   `tfsdk:"options_traverse"`
	// Computed outputs
	ActionID types.String  `tfsdk:"action_id"`
	JobID    types.Int64   `tfsdk:"job_id"`
	State    types.String  `tfsdk:"state"`
	Progress types.Float64 `tfsdk:"progress"`
	Result   types.String  `tfsdk:"result"`
	Error    types.String  `tfsdk:"error"`
}

func NewFilesystemChownResource() resource.Resource {
	return &FilesystemChownResource{}
}

func (r *FilesystemChownResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_chown"
}

func (r *FilesystemChownResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Change owner or group of file at `path`.  `uid` and `gid` specify new owner of the file. If either key is absent or None, then existing value on the file is not changed.  `user` and `group` alternativ",
		Attributes: map[string]schema.Attribute{
			"path":              schema.StringAttribute{Required: true, MarkdownDescription: "Filesystem path to modify.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"uid":               schema.Int64Attribute{Optional: true, MarkdownDescription: "Numeric user ID to set as owner. `null` to leave unchanged.", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"user":              schema.StringAttribute{Optional: true, MarkdownDescription: "Username to set as owner. `null` to leave unchanged.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"gid":               schema.Int64Attribute{Optional: true, MarkdownDescription: "Numeric group ID to set as group owner. `null` to leave unchanged.", PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"group":             schema.StringAttribute{Optional: true, MarkdownDescription: "Group name to set as group owner. `null` to leave unchanged.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"options_recursive": schema.BoolAttribute{Optional: true, MarkdownDescription: "Whether to apply the operation recursively to subdirectories.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
			"options_traverse":  schema.BoolAttribute{Optional: true, MarkdownDescription: "If set do not limit to single dataset / filesystem.", PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
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

func (r *FilesystemChownResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FilesystemChownResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FilesystemChownResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build parameters
	params := map[string]interface{}{}
	params["path"] = data.Path.ValueString()
	if !data.Uid.IsNull() && !data.Uid.IsUnknown() {
		params["uid"] = data.Uid.ValueInt64()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Gid.IsNull() && !data.Gid.IsUnknown() {
		params["gid"] = data.Gid.ValueInt64()
	}
	if !data.Group.IsNull() && !data.Group.IsUnknown() {
		params["group"] = data.Group.ValueString()
	}
	optionsOpts := map[string]interface{}{}
	if !data.OptionsRecursive.IsNull() {
		optionsOpts["recursive"] = data.OptionsRecursive.ValueBool()
	}
	if !data.OptionsTraverse.IsNull() {
		optionsOpts["traverse"] = data.OptionsTraverse.ValueBool()
	}
	if len(optionsOpts) > 0 {
		params["options"] = optionsOpts
	}
	paramsArr := []interface{}{params}

	// Execute action
	result, err := r.client.Call("filesystem.chown", paramsArr)
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute filesystem.chown: %s", err.Error()))
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
	data.ActionID = types.StringValue(fmt.Sprintf("filesystem.chown-%d", time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemChownResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are immutable - just return current state
	var data FilesystemChownResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *FilesystemChownResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Actions cannot be updated, only recreated")
}

func (r *FilesystemChownResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op - actions cannot be undone
}
