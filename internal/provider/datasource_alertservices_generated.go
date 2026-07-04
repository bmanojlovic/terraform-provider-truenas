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

var _ datasource.DataSource = &AlertservicesDataSource{}

func NewAlertservicesDataSource() datasource.DataSource {
	return &AlertservicesDataSource{}
}

type AlertservicesDataSource struct {
	client *client.Client
}

type AlertservicesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type AlertservicesItemModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Attributes types.String `tfsdk:"attributes"`
	Level      types.String `tfsdk:"level"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	TypeTitle  types.String `tfsdk:"type__title"`
}

func (d *AlertservicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alertservices"
}

func (d *AlertservicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query alertservices",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of alertservices resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable name for the alert service.",
						},
						"attributes": schema.StringAttribute{
							Computed:    true,
							Description: "Service-specific configuration attributes (credentials, endpoints, etc.).",
						},
						"level": schema.StringAttribute{
							Computed:    true,
							Description: "Minimum alert severity level that triggers notifications through this service.",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the alert service is active and will send notifications.",
						},
						"type__title": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable title for the alert service type.",
						},
					},
				},
			},
		},
	}
}

func (d *AlertservicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertservicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AlertservicesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("alertservice.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query alertservices: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]AlertservicesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := AlertservicesItemModel{}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			itemModel.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["level"]; ok && v != nil {
			itemModel.Level = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type__title"]; ok && v != nil {
			itemModel.TypeTitle = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attributes":  types.StringType,
			"enabled":     types.BoolType,
			"id":          types.StringType,
			"level":       types.StringType,
			"name":        types.StringType,
			"type__title": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
