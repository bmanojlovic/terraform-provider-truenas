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

var _ datasource.DataSource = &ServicesDataSource{}

func NewServicesDataSource() datasource.DataSource {
	return &ServicesDataSource{}
}

type ServicesDataSource struct {
	client *client.Client
}

type ServicesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type ServicesItemModel struct {
	ID types.String `tfsdk:"id"`
	Service types.String `tfsdk:"service"`
	Enable types.Bool `tfsdk:"enable"`
	State types.String `tfsdk:"state"`
}

func (d *ServicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *ServicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query all system services with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of services resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"service": schema.StringAttribute{
				Computed: true,
				Description: "Name of the system service.",
			},
			"enable": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the service is enabled to start on boot.",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Current state of the service (e.g., 'RUNNING', 'STOPPED').",
			},
					},
				},
			},
		},
	}
}

func (d *ServicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServicesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("service.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query services: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]ServicesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := ServicesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["service"]; ok && v != nil {
			itemModel.Service = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enable": types.BoolType,
			"id": types.StringType,
			"service": types.StringType,
			"state": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
