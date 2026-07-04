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

var _ datasource.DataSource = &SystemNtpserversDataSource{}

func NewSystemNtpserversDataSource() datasource.DataSource {
	return &SystemNtpserversDataSource{}
}

type SystemNtpserversDataSource struct {
	client *client.Client
}

type SystemNtpserversDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type SystemNtpserversItemModel struct {
	ID      types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
	Burst   types.Bool   `tfsdk:"burst"`
	Iburst  types.Bool   `tfsdk:"iburst"`
	Prefer  types.Bool   `tfsdk:"prefer"`
	Minpoll types.Int64  `tfsdk:"minpoll"`
	Maxpoll types.Int64  `tfsdk:"maxpoll"`
}

func (d *SystemNtpserversDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_ntpservers"
}

func (d *SystemNtpserversDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query system_ntpservers",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of system_ntpservers resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"address": schema.StringAttribute{
							Computed:    true,
							Description: "Hostname or IP address of the NTP server.",
						},
						"burst": schema.BoolAttribute{
							Computed:    true,
							Description: "Send a burst of packets when the server is reachable.",
						},
						"iburst": schema.BoolAttribute{
							Computed:    true,
							Description: "Send a burst of packets when the server is unreachable.",
						},
						"prefer": schema.BoolAttribute{
							Computed:    true,
							Description: "Mark this server as preferred for time synchronization.",
						},
						"minpoll": schema.Int64Attribute{
							Computed:    true,
							Description: "Minimum polling interval (log2 seconds).",
						},
						"maxpoll": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum polling interval (log2 seconds).",
						},
					},
				},
			},
		},
	}
}

func (d *SystemNtpserversDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemNtpserversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemNtpserversDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("system.ntpserver.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query system_ntpservers: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]SystemNtpserversItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := SystemNtpserversItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["address"]; ok && v != nil {
			itemModel.Address = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["burst"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Burst = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["iburst"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Iburst = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["prefer"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Prefer = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["minpoll"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Minpoll = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["maxpoll"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Maxpoll = types.Int64Value(int64(fv))
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"address": types.StringType,
			"burst":   types.BoolType,
			"iburst":  types.BoolType,
			"id":      types.StringType,
			"maxpoll": types.Int64Type,
			"minpoll": types.Int64Type,
			"prefer":  types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
