package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &DiskDataSource{}

func NewDiskDataSource() datasource.DataSource {
	return &DiskDataSource{}
}

type DiskDataSource struct {
	client *client.Client
}

type DiskDataSourceModel struct {
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

func (d *DiskDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disk"
}

func (d *DiskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
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
	}
}

func (d *DiskDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DiskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DiskDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("disk.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read disk: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["identifier"]; ok && v != nil {
			data.Identifier = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["subsystem"]; ok && v != nil {
			data.Subsystem = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["number"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Number = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["serial"]; ok && v != nil {
			data.Serial = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lunid"]; ok && v != nil {
			data.Lunid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["size"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Size = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["transfermode"]; ok && v != nil {
			data.Transfermode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hddstandby"]; ok && v != nil {
			data.Hddstandby = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["advpowermgmt"]; ok && v != nil {
			data.Advpowermgmt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["expiretime"]; ok && v != nil {
			data.Expiretime = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["model"]; ok && v != nil {
			data.Model = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["rotationrate"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Rotationrate = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["zfs_guid"]; ok && v != nil {
			data.ZfsGuid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["bus"]; ok && v != nil {
			data.Bus = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["devname"]; ok && v != nil {
			data.Devname = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enclosure"]; ok && v != nil {
			data.Enclosure = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pool"]; ok && v != nil {
			data.Pool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["passwd"]; ok && v != nil {
			data.Passwd = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["kmip_uid"]; ok && v != nil {
			data.KmipUid = types.StringValue(fmt.Sprintf("%v", v))
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
