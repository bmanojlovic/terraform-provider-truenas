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
	ID           types.String `tfsdk:"id"`
	Identifier   types.String `tfsdk:"identifier"`
	Name         types.String `tfsdk:"name"`
	Subsystem    types.String `tfsdk:"subsystem"`
	Number       types.Int64  `tfsdk:"number"`
	Serial       types.String `tfsdk:"serial"`
	Lunid        types.String `tfsdk:"lunid"`
	Size         types.Int64  `tfsdk:"size"`
	Description  types.String `tfsdk:"description"`
	Transfermode types.String `tfsdk:"transfermode"`
	Hddstandby   types.String `tfsdk:"hddstandby"`
	Advpowermgmt types.String `tfsdk:"advpowermgmt"`
	Expiretime   types.String `tfsdk:"expiretime"`
	Model        types.String `tfsdk:"model"`
	Rotationrate types.Int64  `tfsdk:"rotationrate"`
	Type         types.String `tfsdk:"type"`
	ZfsGuid      types.String `tfsdk:"zfs_guid"`
	Bus          types.String `tfsdk:"bus"`
	Devname      types.String `tfsdk:"devname"`
	Enclosure    types.String `tfsdk:"enclosure"`
	Pool         types.String `tfsdk:"pool"`
	Passwd       types.String `tfsdk:"passwd"`
	KmipUid      types.String `tfsdk:"kmip_uid"`
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
				Computed:    true,
				Description: "Unique identifier for the disk device.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "System name of the disk device.",
			},
			"subsystem": schema.StringAttribute{
				Computed:    true,
				Description: "Storage subsystem type.",
			},
			"number": schema.Int64Attribute{
				Computed:    true,
				Description: "Numeric identifier assigned to the disk.",
			},
			"serial": schema.StringAttribute{
				Computed:    true,
				Description: "Manufacturer serial number of the disk.",
			},
			"lunid": schema.StringAttribute{
				Computed:    true,
				Description: "Logical unit number identifier or `null` if not applicable.",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Total size of the disk in bytes. `null` if not available.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description of the disk device.",
			},
			"transfermode": schema.StringAttribute{
				Computed:    true,
				Description: "Data transfer mode and capabilities of the disk.",
			},
			"hddstandby": schema.StringAttribute{
				Computed:    true,
				Description: "Hard disk standby timer in minutes or `ALWAYS ON` to disable standby.",
			},
			"advpowermgmt": schema.StringAttribute{
				Computed:    true,
				Description: "Advanced power management level or `DISABLED` to turn off power management.",
			},
			"expiretime": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration timestamp for disk data or `null` if not applicable.",
			},
			"model": schema.StringAttribute{
				Computed:    true,
				Description: "Manufacturer model name/number of the disk. `null` if not available.",
			},
			"rotationrate": schema.Int64Attribute{
				Computed:    true,
				Description: "Disk rotation speed in RPM or `null` for SSDs and unknown devices.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Disk type classification or `null` if not determined.",
			},
			"zfs_guid": schema.StringAttribute{
				Computed:    true,
				Description: "ZFS globally unique identifier for this disk or `null` if not used in ZFS.",
			},
			"bus": schema.StringAttribute{
				Computed:    true,
				Description: "System bus type the disk is connected to.",
			},
			"devname": schema.StringAttribute{
				Computed:    true,
				Description: "Device name in the operating system.",
			},
			"enclosure": schema.StringAttribute{
				Computed:    true,
				Description: "Physical enclosure information or `null` if not in an enclosure.",
			},
			"pool": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the storage pool this disk belongs to. `null` if not part of any pool.",
			},
			"passwd": schema.StringAttribute{
				Computed:    true,
				Description: "Disk encryption password (masked for security).",
			},
			"kmip_uid": schema.StringAttribute{
				Computed:    true,
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

	if v, ok := resultMap["identifier"]; ok {
		switch val := v.(type) {
		case string:
			data.Identifier = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Identifier = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Identifier = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["subsystem"]; ok {
		switch val := v.(type) {
		case string:
			data.Subsystem = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Subsystem = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Subsystem = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["number"]; ok {
		switch val := v.(type) {
		case float64:
			data.Number = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Number = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["serial"]; ok {
		switch val := v.(type) {
		case string:
			data.Serial = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Serial = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Serial = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["lunid"]; ok {
		if v == nil {
			data.Lunid = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Lunid = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Lunid = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Lunid = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["size"]; ok {
		if v == nil {
			data.Size = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.Size = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.Size = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["description"]; ok {
		switch val := v.(type) {
		case string:
			data.Description = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Description = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transfermode"]; ok {
		switch val := v.(type) {
		case string:
			data.Transfermode = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Transfermode = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Transfermode = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["hddstandby"]; ok {
		switch val := v.(type) {
		case string:
			data.Hddstandby = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Hddstandby = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Hddstandby = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["advpowermgmt"]; ok {
		switch val := v.(type) {
		case string:
			data.Advpowermgmt = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Advpowermgmt = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Advpowermgmt = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["expiretime"]; ok {
		if v == nil {
			data.Expiretime = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Expiretime = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Expiretime = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Expiretime = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["model"]; ok {
		if v == nil {
			data.Model = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Model = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Model = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Model = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["rotationrate"]; ok {
		if v == nil {
			data.Rotationrate = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.Rotationrate = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.Rotationrate = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["type"]; ok {
		if v == nil {
			data.Type = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Type = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Type = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["zfs_guid"]; ok {
		if v == nil {
			data.ZfsGuid = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ZfsGuid = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ZfsGuid = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ZfsGuid = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["bus"]; ok {
		switch val := v.(type) {
		case string:
			data.Bus = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Bus = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Bus = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["devname"]; ok {
		switch val := v.(type) {
		case string:
			data.Devname = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Devname = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Devname = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["enclosure"]; ok {
		if v == nil {
			data.Enclosure = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Enclosure = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Enclosure = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Enclosure = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["pool"]; ok {
		if v == nil {
			data.Pool = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Pool = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Pool = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Pool = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["passwd"]; ok {
		switch val := v.(type) {
		case string:
			data.Passwd = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Passwd = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Passwd = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["kmip_uid"]; ok {
		if v == nil {
			data.KmipUid = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.KmipUid = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.KmipUid = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.KmipUid = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Identifier.IsUnknown() {
		data.Identifier = types.StringNull()
	}
	if data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	if data.Subsystem.IsUnknown() {
		data.Subsystem = types.StringNull()
	}
	if data.Number.IsUnknown() {
		data.Number = types.Int64Null()
	}
	if data.Serial.IsUnknown() {
		data.Serial = types.StringNull()
	}
	if data.Lunid.IsUnknown() {
		data.Lunid = types.StringNull()
	}
	if data.Size.IsUnknown() {
		data.Size = types.Int64Null()
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Transfermode.IsUnknown() {
		data.Transfermode = types.StringNull()
	}
	if data.Hddstandby.IsUnknown() {
		data.Hddstandby = types.StringNull()
	}
	if data.Advpowermgmt.IsUnknown() {
		data.Advpowermgmt = types.StringNull()
	}
	if data.Expiretime.IsUnknown() {
		data.Expiretime = types.StringNull()
	}
	if data.Model.IsUnknown() {
		data.Model = types.StringNull()
	}
	if data.Rotationrate.IsUnknown() {
		data.Rotationrate = types.Int64Null()
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.ZfsGuid.IsUnknown() {
		data.ZfsGuid = types.StringNull()
	}
	if data.Bus.IsUnknown() {
		data.Bus = types.StringNull()
	}
	if data.Devname.IsUnknown() {
		data.Devname = types.StringNull()
	}
	if data.Enclosure.IsUnknown() {
		data.Enclosure = types.StringNull()
	}
	if data.Pool.IsUnknown() {
		data.Pool = types.StringNull()
	}
	if data.Passwd.IsUnknown() {
		data.Passwd = types.StringNull()
	}
	if data.KmipUid.IsUnknown() {
		data.KmipUid = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
