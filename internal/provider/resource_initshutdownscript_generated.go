package provider

import (
	"context"
	"fmt"
	"strconv"

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

func (r *InitshutdownscriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS initshutdownscript resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"command": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"script": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"when": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"timeout": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
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
	params["type"] = data.Type.ValueString()
	if !data.Command.IsNull() {
		params["command"] = data.Command.ValueString()
	}
	if !data.Script.IsNull() {
		params["script"] = data.Script.ValueString()
	}
	params["when"] = data.When.ValueString()
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
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

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

	_, err := r.client.Call("initshutdownscript.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InitshutdownscriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InitshutdownscriptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state InitshutdownscriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["type"] = data.Type.ValueString()
	if !data.Command.IsNull() {
		params["command"] = data.Command.ValueString()
	}
	if !data.Script.IsNull() {
		params["script"] = data.Script.ValueString()
	}
	params["when"] = data.When.ValueString()
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Timeout.IsNull() {
		params["timeout"] = data.Timeout.ValueInt64()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("initshutdownscript.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InitshutdownscriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InitshutdownscriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("initshutdownscript.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
