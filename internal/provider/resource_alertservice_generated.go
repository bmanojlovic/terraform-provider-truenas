package provider

import (
	"context"
	"fmt"
	"strconv"

	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type AlertserviceResource struct {
	client *client.Client
}

type AlertserviceResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Attributes types.String `tfsdk:"attributes"`
	Level types.String `tfsdk:"level"`
	Enabled types.Bool `tfsdk:"enabled"`
}

func NewAlertserviceResource() resource.Resource {
	return &AlertserviceResource{}
}

func (r *AlertserviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alertservice"
}

func (r *AlertserviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AlertserviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an Alert Service of specified `type`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Human-readable name for the alert service.",
			},
			"attributes": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Service-specific configuration attributes (credentials, endpoints, etc.).",
			},
			"level": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Minimum alert severity level that triggers notifications through this service.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the alert service is active and will send notifications.",
			},
		},
	}
}

func (r *AlertserviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertserviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AlertserviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Level.IsNull() {
		params["level"] = data.Level.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	result, err := r.client.Call("alertservice.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create alertservice: %s", err))
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

func (r *AlertserviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AlertserviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("alertservice.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read alertservice: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			data.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["level"]; ok && v != nil {
			data.Level = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertserviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AlertserviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AlertserviceResourceModel
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Level.IsNull() {
		params["level"] = data.Level.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	_, err = r.client.Call("alertservice.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update alertservice: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertserviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AlertserviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("alertservice.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete alertservice: %s", err))
		return
	}
}
