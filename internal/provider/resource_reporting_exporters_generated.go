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

type ReportingExportersResource struct {
	client *client.Client
}

type ReportingExportersResourceModel struct {
	ID types.String `tfsdk:"id"`
	Enabled types.Bool `tfsdk:"enabled"`
	Attributes types.String `tfsdk:"attributes"`
	Name types.String `tfsdk:"name"`
}

func NewReportingExportersResource() resource.Resource {
	return &ReportingExportersResource{}
}

func (r *ReportingExportersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reporting_exporters"
}

func (r *ReportingExportersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ReportingExportersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a specific reporting exporter configuration containing required details for exporting reporting metrics.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"enabled": schema.BoolAttribute{
				Required: true,
				Optional: false,
				Description: "Whether this exporter is enabled and active.",
			},
			"attributes": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Specific attributes for the exporter.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "User defined name of exporter configuration.",
			},
		},
	}
}

func (r *ReportingExportersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReportingExportersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}

	result, err := r.client.Call("reporting.exporters.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create reporting_exporters: %s", err))
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

func (r *ReportingExportersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("reporting.exporters.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read reporting_exporters: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			data.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReportingExportersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ReportingExportersResourceModel
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
	if !data.Attributes.IsNull() {
		var attributesObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Attributes.ValueString()), &attributesObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse attributes: %s", err))
			return
		}
		params["attributes"] = attributesObj
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}

	_, err = r.client.Call("reporting.exporters.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update reporting_exporters: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReportingExportersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("reporting.exporters.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete reporting_exporters: %s", err))
		return
	}
}
