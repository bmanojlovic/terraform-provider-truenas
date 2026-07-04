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

var _ datasource.DataSource = &AppIxVolumesDataSource{}

func NewAppIxVolumesDataSource() datasource.DataSource {
	return &AppIxVolumesDataSource{}
}

type AppIxVolumesDataSource struct {
	client *client.Client
}

type AppIxVolumesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type AppIxVolumesItemModel struct {
	ID      types.String `tfsdk:"id"`
	AppName types.String `tfsdk:"app_name"`
	Name    types.String `tfsdk:"name"`
}

func (d *AppIxVolumesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_ix_volumes"
}

func (d *AppIxVolumesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ix-volumes with `filters` and `options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of app_ix_volumes resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"app_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the application that owns this iX volume.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the iX volume used for persistent storage.",
						},
					},
				},
			},
		},
	}
}

func (d *AppIxVolumesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppIxVolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppIxVolumesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("app.ix_volume.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query app_ix_volumes: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]AppIxVolumesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := AppIxVolumesItemModel{}
		if v, ok := resultMap["app_name"]; ok && v != nil {
			itemModel.AppName = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"app_name": types.StringType,
			"name":     types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
