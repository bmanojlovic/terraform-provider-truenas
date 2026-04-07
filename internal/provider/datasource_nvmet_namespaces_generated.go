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

var _ datasource.DataSource = &NvmetNamespacesDataSource{}

func NewNvmetNamespacesDataSource() datasource.DataSource {
	return &NvmetNamespacesDataSource{}
}

type NvmetNamespacesDataSource struct {
	client *client.Client
}

type NvmetNamespacesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type NvmetNamespacesItemModel struct {
	ID types.String `tfsdk:"id"`
	Nsid types.Int64 `tfsdk:"nsid"`
	Subsys types.String `tfsdk:"subsys"`
	DeviceType types.String `tfsdk:"device_type"`
	DevicePath types.String `tfsdk:"device_path"`
	Filesize types.Int64 `tfsdk:"filesize"`
	DeviceUuid types.String `tfsdk:"device_uuid"`
	DeviceNguid types.String `tfsdk:"device_nguid"`
	Enabled types.Bool `tfsdk:"enabled"`
	Locked types.Bool `tfsdk:"locked"`
}

func (d *NvmetNamespacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_namespaces"
}

func (d *NvmetNamespacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query nvmet_namespaces",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of nvmet_namespaces resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"nsid": schema.Int64Attribute{
				Computed: true,
				Description: "Namespace ID (NSID).  Each namespace within a subsystem has an associated NSID, unique within that s",
			},
			"subsys": schema.StringAttribute{
				Computed: true,
				Description: "NVMe-oF subsystem that contains this namespace.",
			},
			"device_type": schema.StringAttribute{
				Computed: true,
				Description: "Type of device (or file) used to implement the namespace. ",
			},
			"device_path": schema.StringAttribute{
				Computed: true,
				Description: "Path to the device or file being used to implement the namespace.  When `device_type` is:  * \"ZVOL\":",
			},
			"filesize": schema.Int64Attribute{
				Computed: true,
				Description: "When `device_type` is \"FILE\" then this will be the size of the file in bytes.",
			},
			"device_uuid": schema.StringAttribute{
				Computed: true,
				Description: "Unique device identifier for the namespace.",
			},
			"device_nguid": schema.StringAttribute{
				Computed: true,
				Description: "Namespace Globally Unique Identifier for the namespace.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "If `enabled` is `False` then the namespace will not be accessible.  Some namespace configuration cha",
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				Description: "Reflect the locked state of the namespace.  The underlying `device_path` could be an encrypted ZVOL,",
			},
					},
				},
			},
		},
	}
}

func (d *NvmetNamespacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NvmetNamespacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NvmetNamespacesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("nvmet.namespace.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query nvmet_namespaces: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]NvmetNamespacesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := NvmetNamespacesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["nsid"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Nsid = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["subsys"]; ok && v != nil {
			itemModel.Subsys = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["device_type"]; ok && v != nil {
			itemModel.DeviceType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["device_path"]; ok && v != nil {
			itemModel.DevicePath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["filesize"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Filesize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["device_uuid"]; ok && v != nil {
			itemModel.DeviceUuid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["device_nguid"]; ok && v != nil {
			itemModel.DeviceNguid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Locked = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"device_nguid": types.StringType,
			"device_path": types.StringType,
			"device_type": types.StringType,
			"device_uuid": types.StringType,
			"enabled": types.BoolType,
			"filesize": types.Int64Type,
			"id": types.StringType,
			"locked": types.BoolType,
			"nsid": types.Int64Type,
			"subsys": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
