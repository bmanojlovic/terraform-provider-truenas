package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &ReportingExporterssDataSource{}

func NewReportingExporterssDataSource() datasource.DataSource {
	return &ReportingExporterssDataSource{}
}

type ReportingExporterssDataSource struct {
	client *client.Client
}

type ReportingExporterssDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type ReportingExporterssItemModel struct {
	ID         types.String `tfsdk:"id"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	Attributes types.String `tfsdk:"attributes"`
	Name       types.String `tfsdk:"name"`
}

func (d *ReportingExporterssDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reporting_exporterss"
}

func (d *ReportingExporterssDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query reporting_exporterss",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of reporting_exporterss resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this exporter is enabled and active.",
						},
						"attributes": schema.StringAttribute{
							Computed:    true,
							Description: "Specific attributes for the exporter.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "User defined name of exporter configuration.",
						},
					},
				},
			},
		},
	}
}

func (d *ReportingExporterssDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *ReportingExporterssDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ReportingExporterssDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("reporting.exporters.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query reporting_exporterss: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]ReportingExporterssItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := ReportingExporterssItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			itemModel.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attributes": types.StringType,
			"enabled":    types.BoolType,
			"id":         types.StringType,
			"name":       types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
