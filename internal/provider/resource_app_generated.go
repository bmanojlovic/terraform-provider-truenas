package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type AppResource struct {
	client *client.Client
}

type AppResourceModel struct {
	ID types.String `tfsdk:"id"`
	StartOnCreate types.Bool `tfsdk:"start_on_create"`
	CustomApp types.Bool `tfsdk:"custom_app"`
	Values types.Object `tfsdk:"values"`
	CustomComposeConfig types.Object `tfsdk:"custom_compose_config"`
	CustomComposeConfigString types.String `tfsdk:"custom_compose_config_string"`
	CatalogApp types.String `tfsdk:"catalog_app"`
	AppName types.String `tfsdk:"app_name"`
	Train types.String `tfsdk:"train"`
	Version types.String `tfsdk:"version"`
}

func NewAppResource() resource.Resource {
	return &AppResource{}
}

func (r *AppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *AppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS app resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"start_on_create": schema.BoolAttribute{
				Optional: true,
				Description: "Start the resource immediately after creation (default: true if not specified)",
			},
			"custom_app": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"values": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"custom_compose_config": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"custom_compose_config_string": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"catalog_app": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"app_name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"train": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"version": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *AppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.CustomApp.IsNull() {
		params["custom_app"] = data.CustomApp.ValueBool()
	}
	if !data.CustomComposeConfigString.IsNull() {
		params["custom_compose_config_string"] = data.CustomComposeConfigString.ValueString()
	}
	if !data.CatalogApp.IsNull() {
		params["catalog_app"] = data.CatalogApp.ValueString()
	}
	params["app_name"] = data.AppName.ValueString()
	if !data.Train.IsNull() {
		params["train"] = data.Train.ValueString()
	}
	if !data.Version.IsNull() {
		params["version"] = data.Version.ValueString()
	}

	result, err := r.client.Call("app.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Handle lifecycle action - start on create if requested
	startOnCreate := true  // default when not specified
	if !data.StartOnCreate.IsNull() {
		startOnCreate = data.StartOnCreate.ValueBool()
	}
	if startOnCreate {
		_, err = r.client.Call("app.start", data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("app.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.CustomApp.IsNull() {
		params["custom_app"] = data.CustomApp.ValueBool()
	}
	if !data.CustomComposeConfigString.IsNull() {
		params["custom_compose_config_string"] = data.CustomComposeConfigString.ValueString()
	}
	if !data.CatalogApp.IsNull() {
		params["catalog_app"] = data.CatalogApp.ValueString()
	}
	params["app_name"] = data.AppName.ValueString()
	if !data.Train.IsNull() {
		params["train"] = data.Train.ValueString()
	}
	if !data.Version.IsNull() {
		params["version"] = data.Version.ValueString()
	}

	_, err := r.client.Call("app.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("app.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
