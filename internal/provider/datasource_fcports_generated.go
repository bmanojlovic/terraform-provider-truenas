package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &FcportsDataSource{}

func NewFcportsDataSource() datasource.DataSource {
	return &FcportsDataSource{}
}

type FcportsDataSource struct {
	client *client.Client
}

type FcportsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type FcportsItemModel struct {
	ID types.String `tfsdk:"id"`
	Port types.String `tfsdk:"port"`
	Wwpn types.String `tfsdk:"wwpn"`
	WwpnB types.String `tfsdk:"wwpn_b"`
	Target types.String `tfsdk:"target"`
}

func (d *FcportsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fcports"
}

func (d *FcportsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query fcports",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of fcports resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"port": schema.StringAttribute{
				Computed: true,
				Description: "Alias name for the Fibre Channel port.",
			},
			"wwpn": schema.StringAttribute{
				Computed: true,
				Description: "World Wide Port Name for port A or `null` if not configured.",
			},
			"wwpn_b": schema.StringAttribute{
				Computed: true,
				Description: "World Wide Port Name for port B or `null` if not configured.",
			},
			"target": schema.StringAttribute{
				Computed: true,
				Description: "Target configuration object or `null` if not configured.",
			},
					},
				},
			},
		},
	}
}

func (d *FcportsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FcportsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FcportsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("fcport.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query fcports: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]FcportsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := FcportsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["port"]; ok && v != nil {
			itemModel.Port = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["wwpn"]; ok && v != nil {
			itemModel.Wwpn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["wwpn_b"]; ok && v != nil {
			itemModel.WwpnB = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["target"]; ok && v != nil {
			itemModel.Target = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id": types.StringType,
			"port": types.StringType,
			"target": types.StringType,
			"wwpn": types.StringType,
			"wwpn_b": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
