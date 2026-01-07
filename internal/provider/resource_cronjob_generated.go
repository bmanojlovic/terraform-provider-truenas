package provider

import (
	"context"
	"fmt"

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
	Schedule types.Object `tfsdk:"schedule"`
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

func (r *CronjobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS cronjob resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"stderr": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"stdout": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"schedule": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"command": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"user": schema.StringAttribute{
				Required: true,
				Optional: false,
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

	params := map[string]interface{}{
		"enabled": data.Enabled.ValueBool(),
		"stderr": data.Stderr.ValueBool(),
		"stdout": data.Stdout.ValueBool(),
		"command": data.Command.ValueString(),
		"description": data.Description.ValueString(),
		"user": data.User.ValueString(),
	}

	result, err := r.client.Call("cronjob.create", params)
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

func (r *CronjobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("cronjob.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronjobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"enabled": data.Enabled.ValueBool(),
		"stderr": data.Stderr.ValueBool(),
		"stdout": data.Stdout.ValueBool(),
		"command": data.Command.ValueString(),
		"description": data.Description.ValueString(),
		"user": data.User.ValueString(),
	}

	_, err := r.client.Call("cronjob.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronjobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CronjobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("cronjob.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
