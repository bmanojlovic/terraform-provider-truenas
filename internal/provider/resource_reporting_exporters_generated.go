package provider

import (
	"context"
	"fmt"

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

func (r *ReportingExportersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS reporting_exporters resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Required: true,
				Optional: false,
			},
			"attributes": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
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

	params := map[string]interface{}{
		"enabled": data.Enabled.ValueBool(),
		"attributes": data.Attributes.ValueString(),
		"name": data.Name.ValueString(),
	}

	result, err := r.client.Call("reporting/exporters.create", params)
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

func (r *ReportingExportersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("reporting/exporters.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReportingExportersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"enabled": data.Enabled.ValueBool(),
		"attributes": data.Attributes.ValueString(),
		"name": data.Name.ValueString(),
	}

	_, err := r.client.Call("reporting/exporters.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReportingExportersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReportingExportersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("reporting/exporters.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
