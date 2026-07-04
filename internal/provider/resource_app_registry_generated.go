package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type AppRegistryResource struct {
	client *client.Client
}

type AppRegistryResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Uri         types.String `tfsdk:"uri"`
}

func NewAppRegistryResource() resource.Resource {
	return &AppRegistryResource{}
}

func (r *AppRegistryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_registry"
}

func (r *AppRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AppRegistryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an app registry entry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Human-readable name for the container registry.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional description of the container registry or `null`.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Username for registry authentication (masked for security).",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Password or access token for registry authentication (masked for security).",
			},
			"uri": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Container registry URI endpoint (defaults to Docker Hub).",
			},
		},
	}
}

func (r *AppRegistryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		params["password"] = data.Password.ValueString()
	}
	if !data.Uri.IsNull() && !data.Uri.IsUnknown() {
		params["uri"] = data.Uri.ValueString()
	}

	result, err := r.client.Call("app.registry.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create app_registry: %s", err))
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
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("app.registry.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back app_registry: %s", err))
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
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["description"]; ok {
		if v == nil {
			data.Description = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Description = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Description = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Description = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["username"]; ok {
		switch val := v.(type) {
		case string:
			data.Username = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Username = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["password"]; ok {
		switch val := v.(type) {
		case string:
			data.Password = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Password = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Uri.IsUnknown() {
		data.Uri = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("app.registry.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read app_registry: %s", err))
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
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["description"]; ok {
		if v == nil {
			data.Description = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Description = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Description = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Description = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["username"]; ok {
		switch val := v.(type) {
		case string:
			data.Username = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Username = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["password"]; ok {
		switch val := v.(type) {
		case string:
			data.Password = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Password = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Uri.IsUnknown() {
		data.Uri = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		params["password"] = data.Password.ValueString()
	}
	if !data.Uri.IsNull() && !data.Uri.IsUnknown() {
		params["uri"] = data.Uri.ValueString()
	}

	_, err = r.client.Call("app.registry.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update app_registry: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("app.registry.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back app_registry: %s", readErr))
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
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["description"]; ok {
		if v == nil {
			data.Description = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Description = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Description = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Description = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["username"]; ok {
		switch val := v.(type) {
		case string:
			data.Username = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Username = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["password"]; ok {
		switch val := v.(type) {
		case string:
			data.Password = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Password = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Uri.IsUnknown() {
		data.Uri = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("app.registry.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete app_registry: %s", err))
		return
	}
}
