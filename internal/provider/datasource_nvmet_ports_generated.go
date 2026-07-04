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

var _ datasource.DataSource = &NvmetPortsDataSource{}

func NewNvmetPortsDataSource() datasource.DataSource {
	return &NvmetPortsDataSource{}
}

type NvmetPortsDataSource struct {
	client *client.Client
}

type NvmetPortsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type NvmetPortsItemModel struct {
	ID             types.String `tfsdk:"id"`
	Index          types.Int64  `tfsdk:"index"`
	AddrTrtype     types.String `tfsdk:"addr_trtype"`
	AddrTrsvcid    types.Int64  `tfsdk:"addr_trsvcid"`
	AddrTraddr     types.String `tfsdk:"addr_traddr"`
	AddrAdrfam     types.String `tfsdk:"addr_adrfam"`
	InlineDataSize types.Int64  `tfsdk:"inline_data_size"`
	MaxQueueSize   types.Int64  `tfsdk:"max_queue_size"`
	PiEnable       types.Bool   `tfsdk:"pi_enable"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (d *NvmetPortsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_ports"
}

func (d *NvmetPortsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query nvmet_ports",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of nvmet_ports resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"index": schema.Int64Attribute{
							Computed:    true,
							Description: "Index of the port, for internal use. ",
						},
						"addr_trtype": schema.StringAttribute{
							Computed:    true,
							Description: "Fabric transport technology name. ",
						},
						"addr_trsvcid": schema.Int64Attribute{
							Computed:    true,
							Description: "Transport-specific TRSVCID field.  When configured for TCP/IP or RDMA this will be the port number. ",
						},
						"addr_traddr": schema.StringAttribute{
							Computed:    true,
							Description: "A transport-specific field identifying the NVMe host port to use for the connection to the controlle",
						},
						"addr_adrfam": schema.StringAttribute{
							Computed:    true,
							Description: "Address family.",
						},
						"inline_data_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum size for inline data transfers or `null` for default.",
						},
						"max_queue_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of queue entries or `null` for default.",
						},
						"pi_enable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether Protection Information (PI) is enabled or `null` for default.",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Port enabled.  When NVMe target is running, cannot make changes to an enabled port. ",
						},
					},
				},
			},
		},
	}
}

func (d *NvmetPortsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NvmetPortsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NvmetPortsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("nvmet.port.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query nvmet_ports: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]NvmetPortsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := NvmetPortsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["index"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Index = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["addr_trtype"]; ok && v != nil {
			itemModel.AddrTrtype = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["addr_trsvcid"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.AddrTrsvcid = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["addr_traddr"]; ok && v != nil {
			itemModel.AddrTraddr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["addr_adrfam"]; ok && v != nil {
			itemModel.AddrAdrfam = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["inline_data_size"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.InlineDataSize = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["max_queue_size"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.MaxQueueSize = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["pi_enable"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.PiEnable = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"addr_adrfam":      types.StringType,
			"addr_traddr":      types.StringType,
			"addr_trsvcid":     types.Int64Type,
			"addr_trtype":      types.StringType,
			"enabled":          types.BoolType,
			"id":               types.StringType,
			"index":            types.Int64Type,
			"inline_data_size": types.Int64Type,
			"max_queue_size":   types.Int64Type,
			"pi_enable":        types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
