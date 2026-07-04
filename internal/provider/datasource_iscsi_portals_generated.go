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

var _ datasource.DataSource = &IscsiPortalsDataSource{}

func NewIscsiPortalsDataSource() datasource.DataSource {
	return &IscsiPortalsDataSource{}
}

type IscsiPortalsDataSource struct {
	client *client.Client
}

type IscsiPortalsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type IscsiPortalsItemModel struct {
	ID      types.String `tfsdk:"id"`
	Tag     types.Int64  `tfsdk:"tag"`
	Comment types.String `tfsdk:"comment"`
}

func (d *IscsiPortalsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_portals"
}

func (d *IscsiPortalsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query iscsi_portals",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of iscsi_portals resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"tag": schema.Int64Attribute{
							Computed:    true,
							Description: "Numeric tag used to associate this portal with iSCSI targets.",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Optional comment describing the portal.",
						},
					},
				},
			},
		},
	}
}

func (d *IscsiPortalsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IscsiPortalsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IscsiPortalsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("iscsi.portal.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query iscsi_portals: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]IscsiPortalsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := IscsiPortalsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["tag"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Tag = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"comment": types.StringType,
			"id":      types.StringType,
			"tag":     types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
