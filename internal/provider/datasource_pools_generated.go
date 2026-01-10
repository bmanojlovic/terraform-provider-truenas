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

var _ datasource.DataSource = &PoolsDataSource{}

func NewPoolsDataSource() datasource.DataSource {
	return &PoolsDataSource{}
}

type PoolsDataSource struct {
	client *client.Client
}

type PoolsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type PoolsItemModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Guid types.String `tfsdk:"guid"`
	Status types.String `tfsdk:"status"`
	Path types.String `tfsdk:"path"`
	Scan types.String `tfsdk:"scan"`
	Expand types.String `tfsdk:"expand"`
	IsUpgraded types.Bool `tfsdk:"is_upgraded"`
	Healthy types.Bool `tfsdk:"healthy"`
	Warning types.Bool `tfsdk:"warning"`
	StatusCode types.String `tfsdk:"status_code"`
	StatusDetail types.String `tfsdk:"status_detail"`
	Size types.Int64 `tfsdk:"size"`
	Allocated types.Int64 `tfsdk:"allocated"`
	Free types.Int64 `tfsdk:"free"`
	Freeing types.Int64 `tfsdk:"freeing"`
	DedupTableSize types.Int64 `tfsdk:"dedup_table_size"`
	DedupTableQuota types.String `tfsdk:"dedup_table_quota"`
	Fragmentation types.String `tfsdk:"fragmentation"`
	SizeStr types.String `tfsdk:"size_str"`
	AllocatedStr types.String `tfsdk:"allocated_str"`
	FreeStr types.String `tfsdk:"free_str"`
	FreeingStr types.String `tfsdk:"freeing_str"`
	Autotrim types.String `tfsdk:"autotrim"`
	Topology types.String `tfsdk:"topology"`
}

func (d *PoolsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pools"
}

func (d *PoolsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query pools resources",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of pools resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the storage pool.",
			},
			"guid": schema.StringAttribute{
				Computed: true,
				Description: "Globally unique identifier (GUID) for this pool.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Current status of the pool.",
			},
			"path": schema.StringAttribute{
				Computed: true,
				Description: "Filesystem path where the pool is mounted.",
			},
			"scan": schema.StringAttribute{
				Computed: true,
				Description: "Information about any active scrub or resilver operation. `null` if no operation is running.",
			},
			"expand": schema.StringAttribute{
				Computed: true,
				Description: "Information about any active pool expansion operation. `null` if no expansion is running.",
			},
			"is_upgraded": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this pool has been upgraded to the latest feature flags.",
			},
			"healthy": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the pool is in a healthy state with no errors or warnings.",
			},
			"warning": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the pool has warning conditions that require attention.",
			},
			"status_code": schema.StringAttribute{
				Computed: true,
				Description: "Detailed status code for the pool condition. `null` if not applicable.",
			},
			"status_detail": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable description of the pool status. `null` if not available.",
			},
			"size": schema.Int64Attribute{
				Computed: true,
				Description: "Total size of the pool in bytes. `null` if not available.",
			},
			"allocated": schema.Int64Attribute{
				Computed: true,
				Description: "Amount of space currently allocated in the pool in bytes. `null` if not available.",
			},
			"free": schema.Int64Attribute{
				Computed: true,
				Description: "Amount of free space available in the pool in bytes. `null` if not available.",
			},
			"freeing": schema.Int64Attribute{
				Computed: true,
				Description: "Amount of space being freed (in bytes) by ongoing operations. `null` if not available.",
			},
			"dedup_table_size": schema.Int64Attribute{
				Computed: true,
				Description: "Size of the deduplication table in bytes. `null` if deduplication is not enabled.",
			},
			"dedup_table_quota": schema.StringAttribute{
				Computed: true,
				Description: "Quota limit for the deduplication table. `null` if no quota is set.",
			},
			"fragmentation": schema.StringAttribute{
				Computed: true,
				Description: "Percentage of pool fragmentation as a string. `null` if not available.",
			},
			"size_str": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable string representation of the pool size. `null` if not available.",
			},
			"allocated_str": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable string representation of allocated space. `null` if not available.",
			},
			"free_str": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable string representation of free space. `null` if not available.",
			},
			"freeing_str": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable string representation of space being freed. `null` if not available.",
			},
			"autotrim": schema.StringAttribute{
				Computed: true,
				Description: "Auto-trim configuration for the pool indicating whether automatic TRIM operations are enabled.",
			},
			"topology": schema.StringAttribute{
				Computed: true,
				Description: "Physical topology and structure of the pool including vdevs. `null` if not available.",
			},
					},
				},
			},
		},
	}
}

func (d *PoolsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PoolsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("pool.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query pools: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]PoolsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := PoolsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["guid"]; ok && v != nil {
			itemModel.Guid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["status"]; ok && v != nil {
			itemModel.Status = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			itemModel.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["scan"]; ok && v != nil {
			itemModel.Scan = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["expand"]; ok && v != nil {
			itemModel.Expand = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["is_upgraded"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.IsUpgraded = types.BoolValue(bv) }
		}
		if v, ok := resultMap["healthy"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Healthy = types.BoolValue(bv) }
		}
		if v, ok := resultMap["warning"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Warning = types.BoolValue(bv) }
		}
		if v, ok := resultMap["status_code"]; ok && v != nil {
			itemModel.StatusCode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["status_detail"]; ok && v != nil {
			itemModel.StatusDetail = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["size"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Size = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["allocated"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Allocated = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["free"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Free = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["freeing"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Freeing = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["dedup_table_size"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.DedupTableSize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["dedup_table_quota"]; ok && v != nil {
			itemModel.DedupTableQuota = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["fragmentation"]; ok && v != nil {
			itemModel.Fragmentation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["size_str"]; ok && v != nil {
			itemModel.SizeStr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["allocated_str"]; ok && v != nil {
			itemModel.AllocatedStr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["free_str"]; ok && v != nil {
			itemModel.FreeStr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["freeing_str"]; ok && v != nil {
			itemModel.FreeingStr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["autotrim"]; ok && v != nil {
			itemModel.Autotrim = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["topology"]; ok && v != nil {
			itemModel.Topology = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"allocated": types.Int64Type,
			"allocated_str": types.StringType,
			"autotrim": types.StringType,
			"dedup_table_quota": types.StringType,
			"dedup_table_size": types.Int64Type,
			"expand": types.StringType,
			"fragmentation": types.StringType,
			"free": types.Int64Type,
			"free_str": types.StringType,
			"freeing": types.Int64Type,
			"freeing_str": types.StringType,
			"guid": types.StringType,
			"healthy": types.BoolType,
			"id": types.StringType,
			"is_upgraded": types.BoolType,
			"name": types.StringType,
			"path": types.StringType,
			"scan": types.StringType,
			"size": types.Int64Type,
			"size_str": types.StringType,
			"status": types.StringType,
			"status_code": types.StringType,
			"status_detail": types.StringType,
			"topology": types.StringType,
			"warning": types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
