package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type AppResource struct {
	client *client.Client
}

type AppResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	CustomApp                 types.Bool   `tfsdk:"custom_app"`
	Values                    types.String `tfsdk:"values"`
	CustomComposeConfig       types.String `tfsdk:"custom_compose_config"`
	CustomComposeConfigString types.String `tfsdk:"custom_compose_config_string"`
	CatalogApp                types.String `tfsdk:"catalog_app"`
	AppName                   types.String `tfsdk:"app_name"`
	Train                     types.String `tfsdk:"train"`
	Version                   types.String `tfsdk:"version"`
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
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"custom_app": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether to create a custom application (`true`) or install from catalog (`false`).",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown(), boolplanmodifier.RequiresReplace()},
			},
			"values": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Updated configuration values for the application.",
			},
			"custom_compose_config": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Updated Docker Compose configuration as a structured object.",
			},
			"custom_compose_config_string": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Updated Docker Compose configuration as a YAML string.",
			},
			"catalog_app": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Name of the catalog application to install. Required when `custom_app` is `false`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"app_name": schema.StringAttribute{
				Required:      true,
				Optional:      false,
				Description:   "Application name must have the following:  * Lowercase alphanumeric characters can be specified. * N",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"train": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "The catalog train to install from.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"version": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "The version of the application to install.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
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
	if !data.CustomApp.IsNull() && !data.CustomApp.IsUnknown() {
		params["custom_app"] = data.CustomApp.ValueBool()
	}
	if !data.Values.IsNull() && !data.Values.IsUnknown() {
		var valuesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Values.ValueString()), &valuesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse values: %s", err))
			return
		}
		params["values"] = valuesObj
	}
	if !data.CustomComposeConfig.IsNull() && !data.CustomComposeConfig.IsUnknown() {
		var custom_compose_configObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.CustomComposeConfig.ValueString()), &custom_compose_configObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse custom_compose_config: %s", err))
			return
		}
		params["custom_compose_config"] = custom_compose_configObj
	}
	if !data.CustomComposeConfigString.IsNull() && !data.CustomComposeConfigString.IsUnknown() {
		params["custom_compose_config_string"] = data.CustomComposeConfigString.ValueString()
	}
	if !data.CatalogApp.IsNull() && !data.CatalogApp.IsUnknown() {
		params["catalog_app"] = data.CatalogApp.ValueString()
	}
	if !data.AppName.IsNull() && !data.AppName.IsUnknown() {
		params["app_name"] = data.AppName.ValueString()
	}
	if !data.Train.IsNull() && !data.Train.IsUnknown() {
		params["train"] = data.Train.ValueString()
	}
	if !data.Version.IsNull() && !data.Version.IsUnknown() {
		params["version"] = data.Version.ValueString()
	}

	result, err := r.client.CallWithJob("app.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create app: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	// Read back to populate computed fields
	id := data.ID.ValueString()
	result, err = r.client.Call("app.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back app: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["custom_app"]; ok {
		if bv, ok := v.(bool); ok {
			data.CustomApp = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["version"]; ok {
		switch val := v.(type) {
		case string:
			data.Version = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Version = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Version = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.CustomApp.IsUnknown() {
		data.CustomApp = types.BoolNull()
	}
	if data.Values.IsUnknown() {
		data.Values = types.StringNull()
	}
	if data.CustomComposeConfig.IsUnknown() {
		data.CustomComposeConfig = types.StringNull()
	}
	if data.CustomComposeConfigString.IsUnknown() {
		data.CustomComposeConfigString = types.StringNull()
	}
	if data.CatalogApp.IsUnknown() {
		data.CatalogApp = types.StringNull()
	}
	if data.Train.IsUnknown() {
		data.Train = types.StringNull()
	}
	if data.Version.IsUnknown() {
		data.Version = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	result, err := r.client.Call("app.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read app: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["custom_app"]; ok {
		if bv, ok := v.(bool); ok {
			data.CustomApp = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["version"]; ok {
		switch val := v.(type) {
		case string:
			data.Version = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Version = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Version = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.CustomApp.IsUnknown() {
		data.CustomApp = types.BoolNull()
	}
	if data.Values.IsUnknown() {
		data.Values = types.StringNull()
	}
	if data.CustomComposeConfig.IsUnknown() {
		data.CustomComposeConfig = types.StringNull()
	}
	if data.CustomComposeConfigString.IsUnknown() {
		data.CustomComposeConfigString = types.StringNull()
	}
	if data.CatalogApp.IsUnknown() {
		data.CatalogApp = types.StringNull()
	}
	if data.Train.IsUnknown() {
		data.Train = types.StringNull()
	}
	if data.Version.IsUnknown() {
		data.Version = types.StringNull()
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

	var id interface{}
	var err error
	id = state.ID.ValueString()

	params := map[string]interface{}{}
	if !data.Values.IsNull() && !data.Values.IsUnknown() {
		var valuesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Values.ValueString()), &valuesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse values: %s", err))
			return
		}
		params["values"] = valuesObj
	}
	if !data.CustomComposeConfig.IsNull() && !data.CustomComposeConfig.IsUnknown() {
		var custom_compose_configObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.CustomComposeConfig.ValueString()), &custom_compose_configObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse custom_compose_config: %s", err))
			return
		}
		params["custom_compose_config"] = custom_compose_configObj
	}
	if !data.CustomComposeConfigString.IsNull() && !data.CustomComposeConfigString.IsUnknown() {
		params["custom_compose_config_string"] = data.CustomComposeConfigString.ValueString()
	}

	_, err = r.client.CallWithJob("app.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update app: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("app.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back app: %s", readErr))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["custom_app"]; ok {
		if bv, ok := v.(bool); ok {
			data.CustomApp = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["version"]; ok {
		switch val := v.(type) {
		case string:
			data.Version = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Version = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Version = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.CustomApp.IsUnknown() {
		data.CustomApp = types.BoolNull()
	}
	if data.Values.IsUnknown() {
		data.Values = types.StringNull()
	}
	if data.CustomComposeConfig.IsUnknown() {
		data.CustomComposeConfig = types.StringNull()
	}
	if data.CustomComposeConfigString.IsUnknown() {
		data.CustomComposeConfigString = types.StringNull()
	}
	if data.CatalogApp.IsUnknown() {
		data.CatalogApp = types.StringNull()
	}
	if data.Train.IsUnknown() {
		data.Train = types.StringNull()
	}
	if data.Version.IsUnknown() {
		data.Version = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := []interface{}{data.ID.ValueString(), map[string]interface{}{}}

	_, _ = r.client.Call("app.stop", data.ID.ValueString())
	time.Sleep(2 * time.Second)

	_, err := r.client.CallWithJob("app.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete app: %s", err))
		return
	}
}
