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

var _ datasource.DataSource = &JbofsDataSource{}

func NewJbofsDataSource() datasource.DataSource {
	return &JbofsDataSource{}
}

type JbofsDataSource struct {
	client *client.Client
}

type JbofsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type JbofsItemModel struct {
	ID           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	MgmtIp1      types.String `tfsdk:"mgmt_ip1"`
	MgmtIp2      types.String `tfsdk:"mgmt_ip2"`
	MgmtUsername types.String `tfsdk:"mgmt_username"`
	MgmtPassword types.String `tfsdk:"mgmt_password"`
	Index        types.Int64  `tfsdk:"index"`
	Uuid         types.String `tfsdk:"uuid"`
}

func (d *JbofsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jbofs"
}

func (d *JbofsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query jbofs",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of jbofs resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Optional description of the JBOF.",
						},
						"mgmt_ip1": schema.StringAttribute{
							Computed:    true,
							Description: "IP of first Redfish management interface.",
						},
						"mgmt_ip2": schema.StringAttribute{
							Computed:    true,
							Description: "Optional IP of second Redfish management interface.",
						},
						"mgmt_username": schema.StringAttribute{
							Computed:    true,
							Description: "Redfish administrative username.",
						},
						"mgmt_password": schema.StringAttribute{
							Computed:    true,
							Description: "Redfish administrative password.",
						},
						"index": schema.Int64Attribute{
							Computed:    true,
							Description: "Index of the JBOF.  Used to determine data plane IP addresses.",
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the JBOF as reported by the enclosure firmware.",
						},
					},
				},
			},
		},
	}
}

func (d *JbofsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *JbofsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JbofsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("jbof.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query jbofs: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]JbofsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := JbofsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_ip1"]; ok && v != nil {
			itemModel.MgmtIp1 = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_ip2"]; ok && v != nil {
			itemModel.MgmtIp2 = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_username"]; ok && v != nil {
			itemModel.MgmtUsername = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mgmt_password"]; ok && v != nil {
			itemModel.MgmtPassword = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["index"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Index = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["uuid"]; ok && v != nil {
			itemModel.Uuid = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description":   types.StringType,
			"id":            types.StringType,
			"index":         types.Int64Type,
			"mgmt_ip1":      types.StringType,
			"mgmt_ip2":      types.StringType,
			"mgmt_password": types.StringType,
			"mgmt_username": types.StringType,
			"uuid":          types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
