package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type ApiKeyResource struct {
	client *client.Client
}

type ApiKeyResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Username  types.String `tfsdk:"username"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	Reset     types.Bool   `tfsdk:"reset"`
}

func NewApiKeyResource() resource.Resource {
	return &ApiKeyResource{}
}

func (r *ApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *ApiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates API Key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Human-readable name for the API key.",
			},
			"username": schema.StringAttribute{
				Required:      true,
				Optional:      false,
				Description:   "",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"expires_at": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Expiration timestamp for the API key or `null` for no expiration.",
			},
			"reset": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to regenerate a new API key value for this entry.",
			},
		},
	}
}

func (r *ApiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApiKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		params["username"] = data.Username.ValueString()
	}
	if !data.ExpiresAt.IsNull() && !data.ExpiresAt.IsUnknown() {
		params["expires_at"] = data.ExpiresAt.ValueString()
	}
	if !data.Reset.IsNull() && !data.Reset.IsUnknown() {
		params["reset"] = data.Reset.ValueBool()
	}

	result, err := r.client.Call("api_key.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create api_key: %s", err))
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
	result, err = r.client.Call("api_key.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back api_key: %s", err))
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
	if v, ok := resultMap["expires_at"]; ok {
		if v == nil {
			data.ExpiresAt = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ExpiresAt = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.ExpiresAt.IsUnknown() {
		data.ExpiresAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApiKeyResourceModel
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

	result, err := r.client.Call("api_key.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read api_key: %s", err))
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
	if v, ok := resultMap["expires_at"]; ok {
		if v == nil {
			data.ExpiresAt = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ExpiresAt = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.ExpiresAt.IsUnknown() {
		data.ExpiresAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApiKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ApiKeyResourceModel
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
	if !data.ExpiresAt.IsNull() && !data.ExpiresAt.IsUnknown() {
		params["expires_at"] = data.ExpiresAt.ValueString()
	}
	if !data.Reset.IsNull() && !data.Reset.IsUnknown() {
		params["reset"] = data.Reset.ValueBool()
	}

	_, err = r.client.Call("api_key.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update api_key: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("api_key.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back api_key: %s", readErr))
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
	if v, ok := resultMap["expires_at"]; ok {
		if v == nil {
			data.ExpiresAt = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ExpiresAt = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.ExpiresAt.IsUnknown() {
		data.ExpiresAt = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApiKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("api_key.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete api_key: %s", err))
		return
	}
}
