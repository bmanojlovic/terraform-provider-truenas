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

var _ datasource.DataSource = &PoolSnapshotsDataSource{}

func NewPoolSnapshotsDataSource() datasource.DataSource {
	return &PoolSnapshotsDataSource{}
}

type PoolSnapshotsDataSource struct {
	client *client.Client
}

type PoolSnapshotsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type PoolSnapshotsItemModel struct {
	ID types.String `tfsdk:"id"`
	Properties types.String `tfsdk:"properties"`
	Pool types.String `tfsdk:"pool"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	SnapshotName types.String `tfsdk:"snapshot_name"`
	Dataset types.String `tfsdk:"dataset"`
	Createtxg types.String `tfsdk:"createtxg"`
	Holds types.String `tfsdk:"holds"`
	Retention types.String `tfsdk:"retention"`
}

func (d *PoolSnapshotsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_snapshots"
}

func (d *PoolSnapshotsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query all ZFS Snapshots with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of pool_snapshots resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"properties": schema.StringAttribute{
				Computed: true,
				Description: "Object mapping ZFS property names to their values and metadata.",
			},
			"pool": schema.StringAttribute{
				Computed: true,
				Description: "Name of the ZFS pool containing this snapshot.",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Full name of the snapshot including dataset path.",
			},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "Type identifier indicating this is a ZFS snapshot.",
			},
			"snapshot_name": schema.StringAttribute{
				Computed: true,
				Description: "Just the snapshot name portion without the dataset path.",
			},
			"dataset": schema.StringAttribute{
				Computed: true,
				Description: "Name of the dataset this snapshot was taken from.",
			},
			"createtxg": schema.StringAttribute{
				Computed: true,
				Description: "Transaction group ID when the snapshot was created.",
			},
			"holds": schema.StringAttribute{
				Computed: true,
				Description: "Returned when options.extra.holds is set.",
			},
			"retention": schema.StringAttribute{
				Computed: true,
				Description: "Returned when options.extra.retention is set.",
			},
					},
				},
			},
		},
	}
}

func (d *PoolSnapshotsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolSnapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PoolSnapshotsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("pool.snapshot.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query pool_snapshots: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]PoolSnapshotsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := PoolSnapshotsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["properties"]; ok && v != nil {
			itemModel.Properties = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pool"]; ok && v != nil {
			itemModel.Pool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapshot_name"]; ok && v != nil {
			itemModel.SnapshotName = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dataset"]; ok && v != nil {
			itemModel.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["createtxg"]; ok && v != nil {
			itemModel.Createtxg = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["holds"]; ok && v != nil {
			itemModel.Holds = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["retention"]; ok && v != nil {
			itemModel.Retention = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"createtxg": types.StringType,
			"dataset": types.StringType,
			"holds": types.StringType,
			"id": types.StringType,
			"name": types.StringType,
			"pool": types.StringType,
			"properties": types.StringType,
			"retention": types.StringType,
			"snapshot_name": types.StringType,
			"type": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
