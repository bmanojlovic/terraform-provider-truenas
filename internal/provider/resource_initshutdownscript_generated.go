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

type InitshutdownscriptResource struct {
	client *client.Client
}

type InitshutdownscriptResourceModel struct {
	ID types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Command types.String `tfsdk:"command"`
	Script types.String `tfsdk:"script"`
	When types.String `tfsdk:"when"`
	Enabled types.Bool `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
	Comment types.String `tfsdk:"comment"`
}

func NewInitshutdownscriptResource() resource.Resource {
	return &InitshutdownscriptResource{}
}

func (r *InitshutdownscriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_initshutdownscript"
}

func (r *InitshutdownscriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *InitshutdownscriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an initshutdown script task.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"type": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Type of init/shutdown script to execute.  * `COMMAND`: Execute a single command * `SCRIPT`: Execute ",
			},
			"command": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Must be given if `type=\"COMMAND\"`.",
			},
			"script": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Must be given if `type=\"SCRIPT\"`.",
			},
			"when": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "* \"PREINIT\": Early in the boot process before all services have started. * \"POSTINIT\": Late in the b",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the init/shutdown script is enabled to execute.",
			},
			"timeout": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "An integer time in seconds that the system should wait for the execution of the script/command.  A h",
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional comment describing the purpose of this script.",
			},
		},
	}
}

func (r *InitshutdownscriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InitshutdownscriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InitshutdownscriptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Command.IsNull() {
		params["command"] = data.Command.ValueString()
	}
	if !data.Script.IsNull() {
		params["script"] = data.Script.ValueString()
	}
	if !data.When.IsNull() {
		params["when"] = data.When.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Timeout.IsNull() {
		params["timeout"] = data.Timeout.ValueInt64()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	result, err := r.client.Call("initshutdownscript.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create initshutdownscript: %s", err))
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

func (r *InitshutdownscriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InitshutdownscriptResourceModel
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

	result, err := r.client.Call("initshutdownscript.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read initshutdownscript: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["type"]; ok && v != nil {
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["command"]; ok && v != nil {
			data.Command = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["script"]; ok && v != nil {
			data.Script = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["when"]; ok && v != nil {
			data.When = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["timeout"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Timeout = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			data.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InitshutdownscriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InitshutdownscriptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state InitshutdownscriptResourceModel
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
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Command.IsNull() {
		params["command"] = data.Command.ValueString()
	}
	if !data.Script.IsNull() {
		params["script"] = data.Script.ValueString()
	}
	if !data.When.IsNull() {
		params["when"] = data.When.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Timeout.IsNull() {
		params["timeout"] = data.Timeout.ValueInt64()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	_, err = r.client.Call("initshutdownscript.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update initshutdownscript: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InitshutdownscriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InitshutdownscriptResourceModel
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

	_, err = r.client.Call("initshutdownscript.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete initshutdownscript: %s", err))
		return
	}
}
