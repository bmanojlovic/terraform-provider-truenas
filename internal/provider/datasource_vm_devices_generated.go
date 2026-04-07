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

var _ datasource.DataSource = &VmDevicesDataSource{}

func NewVmDevicesDataSource() datasource.DataSource {
	return &VmDevicesDataSource{}
}

type VmDevicesDataSource struct {
	client *client.Client
}

type VmDevicesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type VmDevicesItemModel struct {
	ID types.String `tfsdk:"id"`
	Attributes types.String `tfsdk:"attributes"`
	Vm types.Int64 `tfsdk:"vm"`
	Order types.Int64 `tfsdk:"order"`
}

func (d *VmDevicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_devices"
}

func (d *VmDevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query vm_devices",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of vm_devices resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"attributes": schema.StringAttribute{
				Computed: true,
				Description: "Device-specific configuration attributes.",
			},
			"vm": schema.Int64Attribute{
				Computed: true,
				Description: "ID of the virtual machine this device belongs to.",
			},
			"order": schema.Int64Attribute{
				Computed: true,
				Description: "Boot order priority for this device (lower numbers boot first).",
			},
					},
				},
			},
		},
	}
}

func (d *VmDevicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VmDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VmDevicesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("vm.device.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query vm_devices: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]VmDevicesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := VmDevicesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			itemModel.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vm"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Vm = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["order"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Order = types.Int64Value(int64(fv)) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attributes": types.StringType,
			"id": types.StringType,
			"order": types.Int64Type,
			"vm": types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
