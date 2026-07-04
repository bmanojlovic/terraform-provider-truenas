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

var _ datasource.DataSource = &IscsiTargetsDataSource{}

func NewIscsiTargetsDataSource() datasource.DataSource {
	return &IscsiTargetsDataSource{}
}

type IscsiTargetsDataSource struct {
	client *client.Client
}

type IscsiTargetsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type IscsiTargetsItemModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Alias           types.String `tfsdk:"alias"`
	Mode            types.String `tfsdk:"mode"`
	RelTgtId        types.Int64  `tfsdk:"rel_tgt_id"`
	IscsiParameters types.String `tfsdk:"iscsi_parameters"`
}

func (d *IscsiTargetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_targets"
}

func (d *IscsiTargetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query iscsi_targets",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of iscsi_targets resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the iSCSI target (maximum 120 characters).",
						},
						"alias": schema.StringAttribute{
							Computed:    true,
							Description: "Optional alias name for the iSCSI target.",
						},
						"mode": schema.StringAttribute{
							Computed:    true,
							Description: "Protocol mode for the target.  * `ISCSI`: iSCSI protocol only * `FC`: Fibre Channel protocol only * ",
						},
						"rel_tgt_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Relative target ID number assigned by the system.",
						},
						"iscsi_parameters": schema.StringAttribute{
							Computed:    true,
							Description: "Optional iSCSI-specific parameters for this target.",
						},
					},
				},
			},
		},
	}
}

func (d *IscsiTargetsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IscsiTargetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IscsiTargetsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("iscsi.target.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query iscsi_targets: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]IscsiTargetsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := IscsiTargetsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["alias"]; ok && v != nil {
			itemModel.Alias = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mode"]; ok && v != nil {
			itemModel.Mode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["rel_tgt_id"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.RelTgtId = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["iscsi_parameters"]; ok && v != nil {
			itemModel.IscsiParameters = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"alias":            types.StringType,
			"id":               types.StringType,
			"iscsi_parameters": types.StringType,
			"mode":             types.StringType,
			"name":             types.StringType,
			"rel_tgt_id":       types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
