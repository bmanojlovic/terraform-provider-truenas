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

var _ datasource.DataSource = &IscsiExtentsDataSource{}

func NewIscsiExtentsDataSource() datasource.DataSource {
	return &IscsiExtentsDataSource{}
}

type IscsiExtentsDataSource struct {
	client *client.Client
}

type IscsiExtentsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type IscsiExtentsItemModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Disk           types.String `tfsdk:"disk"`
	Serial         types.String `tfsdk:"serial"`
	Path           types.String `tfsdk:"path"`
	Filesize       types.Int64  `tfsdk:"filesize"`
	Blocksize      types.Int64  `tfsdk:"blocksize"`
	Pblocksize     types.Bool   `tfsdk:"pblocksize"`
	AvailThreshold types.Int64  `tfsdk:"avail_threshold"`
	Comment        types.String `tfsdk:"comment"`
	Naa            types.String `tfsdk:"naa"`
	InsecureTpc    types.Bool   `tfsdk:"insecure_tpc"`
	Xen            types.Bool   `tfsdk:"xen"`
	Rpm            types.String `tfsdk:"rpm"`
	Ro             types.Bool   `tfsdk:"ro"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Vendor         types.String `tfsdk:"vendor"`
	ProductId      types.String `tfsdk:"product_id"`
	Locked         types.Bool   `tfsdk:"locked"`
}

func (d *IscsiExtentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_extents"
}

func (d *IscsiExtentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query iscsi_extents",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of iscsi_extents resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the iSCSI extent.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the extent storage backend.",
						},
						"disk": schema.StringAttribute{
							Computed:    true,
							Description: "Disk device to use for the extent or `null` if using a file.",
						},
						"serial": schema.StringAttribute{
							Computed:    true,
							Description: "Serial number for the extent or `null` to auto-generate.",
						},
						"path": schema.StringAttribute{
							Computed:    true,
							Description: "File path for file-based extents or `null` if using a disk.",
						},
						"filesize": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the file-based extent in bytes.",
						},
						"blocksize": schema.Int64Attribute{
							Computed:    true,
							Description: "Block size for the extent in bytes.",
						},
						"pblocksize": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to use physical block size reporting.",
						},
						"avail_threshold": schema.Int64Attribute{
							Computed:    true,
							Description: "Available space threshold percentage or `null` to disable.",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Optional comment describing the extent.",
						},
						"naa": schema.StringAttribute{
							Computed:    true,
							Description: "Network Address Authority (NAA) identifier for the extent.",
						},
						"insecure_tpc": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to enable insecure Third Party Copy (TPC) operations.",
						},
						"xen": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to enable Xen compatibility mode.",
						},
						"rpm": schema.StringAttribute{
							Computed:    true,
							Description: "Reported RPM type for the extent.",
						},
						"ro": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the extent is read-only.",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the extent is enabled and available for use.",
						},
						"vendor": schema.StringAttribute{
							Computed:    true,
							Description: "Vendor string reported by the extent.",
						},
						"product_id": schema.StringAttribute{
							Computed:    true,
							Description: "Product ID string for the extent or `null` for default.",
						},
						"locked": schema.BoolAttribute{
							Computed:    true,
							Description: "Read-only value indicating whether the iscsi extent is located on a locked dataset.  - `true`: The e",
						},
					},
				},
			},
		},
	}
}

func (d *IscsiExtentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IscsiExtentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IscsiExtentsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("iscsi.extent.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query iscsi_extents: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]IscsiExtentsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := IscsiExtentsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["disk"]; ok && v != nil {
			itemModel.Disk = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["serial"]; ok && v != nil {
			itemModel.Serial = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			itemModel.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["filesize"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Filesize = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["blocksize"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Blocksize = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["pblocksize"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Pblocksize = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["avail_threshold"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.AvailThreshold = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["naa"]; ok && v != nil {
			itemModel.Naa = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["insecure_tpc"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.InsecureTpc = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["xen"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Xen = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["rpm"]; ok && v != nil {
			itemModel.Rpm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ro"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Ro = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["vendor"]; ok && v != nil {
			itemModel.Vendor = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["product_id"]; ok && v != nil {
			itemModel.ProductId = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Locked = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"avail_threshold": types.Int64Type,
			"blocksize":       types.Int64Type,
			"comment":         types.StringType,
			"disk":            types.StringType,
			"enabled":         types.BoolType,
			"filesize":        types.Int64Type,
			"id":              types.StringType,
			"insecure_tpc":    types.BoolType,
			"locked":          types.BoolType,
			"naa":             types.StringType,
			"name":            types.StringType,
			"path":            types.StringType,
			"pblocksize":      types.BoolType,
			"product_id":      types.StringType,
			"ro":              types.BoolType,
			"rpm":             types.StringType,
			"serial":          types.StringType,
			"type":            types.StringType,
			"vendor":          types.StringType,
			"xen":             types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
