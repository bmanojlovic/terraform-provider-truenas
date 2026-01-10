package provider

import (
	"context"
	"fmt"
	"strconv"

	"time"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	Values types.String `tfsdk:"values"`
	CustomComposeConfig types.String `tfsdk:"custom_compose_config"`
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

func (r *AppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an app with `app_name` using `catalog_app` with `train` and `version`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"start_on_create": schema.BoolAttribute{Optional: true, Description: "Start the resource immediately after creation (default: true)"},
			"custom_app": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to create a custom application (`true`) or install from catalog (`false`).",
			},
			"values": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Configuration values for the application installation.",
			},
			"custom_compose_config": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Docker Compose configuration as a structured object for custom applications.",
			},
			"custom_compose_config_string": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Docker Compose configuration as a YAML string for custom applications.",
			},
			"catalog_app": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Name of the catalog application to install. Required when `custom_app` is `false`.",
			},
			"app_name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Application name must have the following:  * Lowercase alphanumeric characters can be specified. * N",
			},
			"train": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The catalog train to install from.",
			},
			"version": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The version of the application to install.",
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
	if !data.Values.IsNull() {
		params["values"] = data.Values.ValueString()
	}
	if !data.CustomComposeConfig.IsNull() {
		params["custom_compose_config"] = data.CustomComposeConfig.ValueString()
	}
	if !data.CustomComposeConfigString.IsNull() {
		params["custom_compose_config_string"] = data.CustomComposeConfigString.ValueString()
	}
	if !data.CatalogApp.IsNull() {
		params["catalog_app"] = data.CatalogApp.ValueString()
	}
	if !data.AppName.IsNull() {
		params["app_name"] = data.AppName.ValueString()
	}
	if !data.Train.IsNull() {
		params["train"] = data.Train.ValueString()
	}
	if !data.Version.IsNull() {
		params["version"] = data.Version.ValueString()
	}

	result, err := r.client.CallWithJob("app.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create app: %s", err))
		return
	}

	// Extract ID from result
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
		vmID, err := strconv.Atoi(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
			return
		}
		_, err = r.client.Call("app.start", vmID)
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

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("app.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read app: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["custom_app"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.CustomApp = types.BoolValue(bv) }
		}
		if v, ok := resultMap["values"]; ok && v != nil {
			data.Values = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["custom_compose_config"]; ok && v != nil {
			data.CustomComposeConfig = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["custom_compose_config_string"]; ok && v != nil {
			data.CustomComposeConfigString = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["catalog_app"]; ok && v != nil {
			data.CatalogApp = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["app_name"]; ok && v != nil {
			data.AppName = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["train"]; ok && v != nil {
			data.Train = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["version"]; ok && v != nil {
			data.Version = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AppResourceModel
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
	if !data.CustomApp.IsNull() {
		params["custom_app"] = data.CustomApp.ValueBool()
	}
	if !data.Values.IsNull() {
		params["values"] = data.Values.ValueString()
	}
	if !data.CustomComposeConfig.IsNull() {
		params["custom_compose_config"] = data.CustomComposeConfig.ValueString()
	}
	if !data.CustomComposeConfigString.IsNull() {
		params["custom_compose_config_string"] = data.CustomComposeConfigString.ValueString()
	}
	if !data.CatalogApp.IsNull() {
		params["catalog_app"] = data.CatalogApp.ValueString()
	}
	if !data.AppName.IsNull() {
		params["app_name"] = data.AppName.ValueString()
	}
	if !data.Train.IsNull() {
		params["train"] = data.Train.ValueString()
	}
	if !data.Version.IsNull() {
		params["version"] = data.Version.ValueString()
	}

	_, err = r.client.CallWithJob("app.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update app: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	// Stop VM before deletion if running
	vmID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}
	_, _ = r.client.Call("app.stop", vmID)  // Ignore errors - VM might already be stopped
	time.Sleep(2 * time.Second)  // Wait for VM to stop

	_, err = r.client.CallWithJob("app.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete app: %s", err))
		return
	}
}
