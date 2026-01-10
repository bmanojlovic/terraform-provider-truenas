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

var _ datasource.DataSource = &DisksDataSource{}

func NewDisksDataSource() datasource.DataSource {
	return &DisksDataSource{}
}

type DisksDataSource struct {
	client *client.Client
}

type DisksDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type DisksItemModel struct {
	ID types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Name types.String `tfsdk:"name"`
	Subsystem types.String `tfsdk:"subsystem"`
	Number types.Int64 `tfsdk:"number"`
	Serial types.String `tfsdk:"serial"`
	Lunid types.String `tfsdk:"lunid"`
	Size types.Int64 `tfsdk:"size"`
	Description types.String `tfsdk:"description"`
	Transfermode types.String `tfsdk:"transfermode"`
	Hddstandby types.String `tfsdk:"hddstandby"`
	Advpowermgmt types.String `tfsdk:"advpowermgmt"`
	Expiretime types.String `tfsdk:"expiretime"`
	Model types.String `tfsdk:"model"`
	Rotationrate types.Int64 `tfsdk:"rotationrate"`
	Type types.String `tfsdk:"type"`
	ZfsGuid types.String `tfsdk:"zfs_guid"`
	Bus types.String `tfsdk:"bus"`
	Devname types.String `tfsdk:"devname"`
	Enclosure types.String `tfsdk:"enclosure"`
	Pool types.String `tfsdk:"pool"`
	Passwd types.String `tfsdk:"passwd"`
	KmipUid types.String `tfsdk:"kmip_uid"`
}

func (d *DisksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disks"
}

func (d *DisksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query disks.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of disks resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"identifier": schema.StringAttribute{
				Computed: true,
				Description: "Unique identifier for the disk device.",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "System name of the disk device.",
			},
			"subsystem": schema.StringAttribute{
				Computed: true,
				Description: "Storage subsystem type.",
			},
			"number": schema.Int64Attribute{
				Computed: true,
				Description: "Numeric identifier assigned to the disk.",
			},
			"serial": schema.StringAttribute{
				Computed: true,
				Description: "Manufacturer serial number of the disk.",
			},
			"lunid": schema.StringAttribute{
				Computed: true,
				Description: "Logical unit number identifier or `null` if not applicable.",
			},
			"size": schema.Int64Attribute{
				Computed: true,
				Description: "Total size of the disk in bytes. `null` if not available.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable description of the disk device.",
			},
			"transfermode": schema.StringAttribute{
				Computed: true,
				Description: "Data transfer mode and capabilities of the disk.",
			},
			"hddstandby": schema.StringAttribute{
				Computed: true,
				Description: "Hard disk standby timer in minutes or `ALWAYS ON` to disable standby.",
			},
			"advpowermgmt": schema.StringAttribute{
				Computed: true,
				Description: "Advanced power management level or `DISABLED` to turn off power management.",
			},
			"expiretime": schema.StringAttribute{
				Computed: true,
				Description: "Expiration timestamp for disk data or `null` if not applicable.",
			},
			"model": schema.StringAttribute{
				Computed: true,
				Description: "Manufacturer model name/number of the disk. `null` if not available.",
			},
			"rotationrate": schema.Int64Attribute{
				Computed: true,
				Description: "Disk rotation speed in RPM or `null` for SSDs and unknown devices.",
			},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "Disk type classification or `null` if not determined.",
			},
			"zfs_guid": schema.StringAttribute{
				Computed: true,
				Description: "ZFS globally unique identifier for this disk or `null` if not used in ZFS.",
			},
			"bus": schema.StringAttribute{
				Computed: true,
				Description: "System bus type the disk is connected to.",
			},
			"devname": schema.StringAttribute{
				Computed: true,
				Description: "Device name in the operating system.",
			},
			"enclosure": schema.StringAttribute{
				Computed: true,
				Description: "Physical enclosure information or `null` if not in an enclosure.",
			},
			"pool": schema.StringAttribute{
				Computed: true,
				Description: "Name of the storage pool this disk belongs to. `null` if not part of any pool.",
			},
			"passwd": schema.StringAttribute{
				Computed: true,
				Description: "Disk encryption password (masked for security).",
			},
			"kmip_uid": schema.StringAttribute{
				Computed: true,
				Description: "KMIP (Key Management Interoperability Protocol) unique identifier or `null`.",
			},
					},
				},
			},
		},
	}
}

func (d *DisksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DisksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DisksDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("disk.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query disks: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]DisksItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := DisksItemModel{}
		if v, ok := resultMap["identifier"]; ok && v != nil {
			itemModel.Identifier = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["subsystem"]; ok && v != nil {
			itemModel.Subsystem = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["number"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Number = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["serial"]; ok && v != nil {
			itemModel.Serial = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lunid"]; ok && v != nil {
			itemModel.Lunid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["size"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Size = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["transfermode"]; ok && v != nil {
			itemModel.Transfermode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hddstandby"]; ok && v != nil {
			itemModel.Hddstandby = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["advpowermgmt"]; ok && v != nil {
			itemModel.Advpowermgmt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["expiretime"]; ok && v != nil {
			itemModel.Expiretime = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["model"]; ok && v != nil {
			itemModel.Model = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["rotationrate"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Rotationrate = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["zfs_guid"]; ok && v != nil {
			itemModel.ZfsGuid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["bus"]; ok && v != nil {
			itemModel.Bus = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["devname"]; ok && v != nil {
			itemModel.Devname = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enclosure"]; ok && v != nil {
			itemModel.Enclosure = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pool"]; ok && v != nil {
			itemModel.Pool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["passwd"]; ok && v != nil {
			itemModel.Passwd = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["kmip_uid"]; ok && v != nil {
			itemModel.KmipUid = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"advpowermgmt": types.StringType,
			"bus": types.StringType,
			"description": types.StringType,
			"devname": types.StringType,
			"enclosure": types.StringType,
			"expiretime": types.StringType,
			"hddstandby": types.StringType,
			"identifier": types.StringType,
			"kmip_uid": types.StringType,
			"lunid": types.StringType,
			"model": types.StringType,
			"name": types.StringType,
			"number": types.Int64Type,
			"passwd": types.StringType,
			"pool": types.StringType,
			"rotationrate": types.Int64Type,
			"serial": types.StringType,
			"size": types.Int64Type,
			"subsystem": types.StringType,
			"transfermode": types.StringType,
			"type": types.StringType,
			"zfs_guid": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
