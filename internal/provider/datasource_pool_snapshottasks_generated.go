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

var _ datasource.DataSource = &PoolSnapshottasksDataSource{}

func NewPoolSnapshottasksDataSource() datasource.DataSource {
	return &PoolSnapshottasksDataSource{}
}

type PoolSnapshottasksDataSource struct {
	client *client.Client
}

type PoolSnapshottasksDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type PoolSnapshottasksItemModel struct {
	ID types.String `tfsdk:"id"`
	Dataset types.String `tfsdk:"dataset"`
	Recursive types.Bool `tfsdk:"recursive"`
	LifetimeValue types.Int64 `tfsdk:"lifetime_value"`
	LifetimeUnit types.String `tfsdk:"lifetime_unit"`
	Enabled types.Bool `tfsdk:"enabled"`
	NamingSchema types.String `tfsdk:"naming_schema"`
	AllowEmpty types.Bool `tfsdk:"allow_empty"`
	Schedule types.String `tfsdk:"schedule"`
	State types.String `tfsdk:"state"`
	VmwareSync types.Bool `tfsdk:"vmware_sync"`
}

func (d *PoolSnapshottasksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_snapshottasks"
}

func (d *PoolSnapshottasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query pool_snapshottasks",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of pool_snapshottasks resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"dataset": schema.StringAttribute{
				Computed: true,
				Description: "The dataset to take snapshots of.",
			},
			"recursive": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to recursively snapshot child datasets.",
			},
			"lifetime_value": schema.Int64Attribute{
				Computed: true,
				Description: "Number of time units to retain snapshots. `lifetime_unit` gives the time unit.",
			},
			"lifetime_unit": schema.StringAttribute{
				Computed: true,
				Description: "Unit of time for snapshot retention.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this periodic snapshot task is enabled.",
			},
			"naming_schema": schema.StringAttribute{
				Computed: true,
				Description: "Naming pattern for generated snapshots using strftime format.",
			},
			"allow_empty": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to take snapshots even if no data has changed.",
			},
			"schedule": schema.StringAttribute{
				Computed: true,
				Description: "Cron schedule for when snapshots should be taken.",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Detailed state information for the task.",
			},
			"vmware_sync": schema.BoolAttribute{
				Computed: true,
				Description: "Whether VMware VMs are synced before taking snapshots.",
			},
					},
				},
			},
		},
	}
}

func (d *PoolSnapshottasksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolSnapshottasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PoolSnapshottasksDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("pool.snapshottask.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query pool_snapshottasks: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]PoolSnapshottasksItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := PoolSnapshottasksItemModel{}
		if v, ok := resultMap["dataset"]; ok && v != nil {
			itemModel.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recursive"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Recursive = types.BoolValue(bv) }
		}
		if v, ok := resultMap["lifetime_value"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.LifetimeValue = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["lifetime_unit"]; ok && v != nil {
			itemModel.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["naming_schema"]; ok && v != nil {
			itemModel.NamingSchema = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["allow_empty"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.AllowEmpty = types.BoolValue(bv) }
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			itemModel.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vmware_sync"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.VmwareSync = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"allow_empty": types.BoolType,
			"dataset": types.StringType,
			"enabled": types.BoolType,
			"id": types.StringType,
			"lifetime_unit": types.StringType,
			"lifetime_value": types.Int64Type,
			"naming_schema": types.StringType,
			"recursive": types.BoolType,
			"schedule": types.StringType,
			"state": types.StringType,
			"vmware_sync": types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
