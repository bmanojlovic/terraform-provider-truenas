package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type CronjobResource struct {
	client *client.Client
}

type CronjobResourceModel struct {
	ID types.String `tfsdk:"id"`
	Enabled types.Bool `tfsdk:"enabled"`
	Stderr types.Bool `tfsdk:"stderr"`
	Stdout types.Bool `tfsdk:"stdout"`
	Schedule types.String `tfsdk:"schedule"`
	Command types.String `tfsdk:"command"`
	Description types.String `tfsdk:"description"`
	User types.String `tfsdk:"user"`
}

func NewCronjobResource() resource.Resource {
	return &CronjobResource{}
}

func (r *CronjobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cronjob"
}

func (r *CronjobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CronjobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new cron job.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the cron job is active and will be executed.",
			},
			"stderr": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to IGNORE standard error (if `false`, it will be added to email).",
			},
			"stdout": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to IGNORE standard output (if `false`, it will be added to email).",
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Cron schedule configuration for when the job runs.",
			},
			"command": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Shell command or script to execute.",
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Human-readable description of what this cron job does.",
			},
			"user": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "System user account to run the command as.",
			},
		},
	}
}

func (r *CronjobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CronjobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Stderr.IsNull() {
		params["stderr"] = data.Stderr.ValueBool()
	}
	if !data.Stdout.IsNull() {
		params["stdout"] = data.Stdout.ValueBool()
	}
	if !data.Schedule.IsNull() {
		params["schedule"] = data.Schedule.ValueString()
	}
	if !data.Command.IsNull() {
		params["command"] = data.Command.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.User.IsNull() {
		params["user"] = data.User.ValueString()
	}

	result, err := r.client.Call("cronjob.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create cronjob: %s", err))
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

func (r *CronjobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("cronjob.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read cronjob: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["stderr"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Stderr = types.BoolValue(bv) }
		}
		if v, ok := resultMap["stdout"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Stdout = types.BoolValue(bv) }
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["command"]; ok && v != nil {
			data.Command = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user"]; ok && v != nil {
			data.User = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronjobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CronjobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Stderr.IsNull() {
		params["stderr"] = data.Stderr.ValueBool()
	}
	if !data.Stdout.IsNull() {
		params["stdout"] = data.Stdout.ValueBool()
	}
	if !data.Schedule.IsNull() {
		params["schedule"] = data.Schedule.ValueString()
	}
	if !data.Command.IsNull() {
		params["command"] = data.Command.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.User.IsNull() {
		params["user"] = data.User.ValueString()
	}

	_, err = r.client.Call("cronjob.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update cronjob: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronjobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("cronjob.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete cronjob: %s", err))
		return
	}
}
