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

var _ datasource.DataSource = &PoolScrubsDataSource{}

func NewPoolScrubsDataSource() datasource.DataSource {
	return &PoolScrubsDataSource{}
}

type PoolScrubsDataSource struct {
	client *client.Client
}

type PoolScrubsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type PoolScrubsItemModel struct {
	ID types.String `tfsdk:"id"`
	Pool types.Int64 `tfsdk:"pool"`
	Threshold types.Int64 `tfsdk:"threshold"`
	Description types.String `tfsdk:"description"`
	Schedule types.String `tfsdk:"schedule"`
	Enabled types.Bool `tfsdk:"enabled"`
	PoolName types.String `tfsdk:"pool_name"`
}

func (d *PoolScrubsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_scrubs"
}

func (d *PoolScrubsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query pool_scrubs",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of pool_scrubs resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"pool": schema.Int64Attribute{
				Computed: true,
				Description: "ID of the pool to scrub.",
			},
			"threshold": schema.Int64Attribute{
				Computed: true,
				Description: "Days before a scrub is due when a scrub should automatically start.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Description or notes for this scrub schedule.",
			},
			"schedule": schema.StringAttribute{
				Computed: true,
				Description: "Cron schedule for when scrubs should run.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this scrub schedule is enabled.",
			},
			"pool_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the pool being scrubbed.",
			},
					},
				},
			},
		},
	}
}

func (d *PoolScrubsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolScrubsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PoolScrubsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("pool.scrub.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query pool_scrubs: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]PoolScrubsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := PoolScrubsItemModel{}
		if v, ok := resultMap["pool"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Pool = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["threshold"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Threshold = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			itemModel.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pool_name"]; ok && v != nil {
			itemModel.PoolName = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"enabled": types.BoolType,
			"id": types.StringType,
			"pool": types.Int64Type,
			"pool_name": types.StringType,
			"schedule": types.StringType,
			"threshold": types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
