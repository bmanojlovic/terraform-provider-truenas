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

var _ datasource.DataSource = &NvmetPortSubsyssDataSource{}

func NewNvmetPortSubsyssDataSource() datasource.DataSource {
	return &NvmetPortSubsyssDataSource{}
}

type NvmetPortSubsyssDataSource struct {
	client *client.Client
}

type NvmetPortSubsyssDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type NvmetPortSubsyssItemModel struct {
	ID types.String `tfsdk:"id"`
	Port types.String `tfsdk:"port"`
	Subsys types.String `tfsdk:"subsys"`
}

func (d *NvmetPortSubsyssDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_port_subsyss"
}

func (d *NvmetPortSubsyssDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query nvmet_port_subsyss",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of nvmet_port_subsyss resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"port": schema.StringAttribute{
				Computed: true,
				Description: "NVMe-oF port that provides access to the subsystem.",
			},
			"subsys": schema.StringAttribute{
				Computed: true,
				Description: "NVMe-oF subsystem that is accessible through the port.",
			},
					},
				},
			},
		},
	}
}

func (d *NvmetPortSubsyssDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NvmetPortSubsyssDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NvmetPortSubsyssDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("nvmet.port_subsys.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query nvmet_port_subsyss: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]NvmetPortSubsyssItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := NvmetPortSubsyssItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["port"]; ok && v != nil {
			itemModel.Port = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["subsys"]; ok && v != nil {
			itemModel.Subsys = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id": types.StringType,
			"port": types.StringType,
			"subsys": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
