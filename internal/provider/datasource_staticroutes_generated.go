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

var _ datasource.DataSource = &StaticroutesDataSource{}

func NewStaticroutesDataSource() datasource.DataSource {
	return &StaticroutesDataSource{}
}

type StaticroutesDataSource struct {
	client *client.Client
}

type StaticroutesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type StaticroutesItemModel struct {
	ID types.String `tfsdk:"id"`
	Destination types.String `tfsdk:"destination"`
	Gateway types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
}

func (d *StaticroutesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staticroutes"
}

func (d *StaticroutesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query staticroutes",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of staticroutes resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"destination": schema.StringAttribute{
				Computed: true,
				Description: "Destination network or host for this static route.",
			},
			"gateway": schema.StringAttribute{
				Computed: true,
				Description: "Gateway IP address for this static route.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Optional description for this static route.",
			},
					},
				},
			},
		},
	}
}

func (d *StaticroutesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StaticroutesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StaticroutesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("staticroute.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query staticroutes: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]StaticroutesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := StaticroutesItemModel{}
		if v, ok := resultMap["destination"]; ok && v != nil {
			itemModel.Destination = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["gateway"]; ok && v != nil {
			itemModel.Gateway = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"destination": types.StringType,
			"gateway": types.StringType,
			"id": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
