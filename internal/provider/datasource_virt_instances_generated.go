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

var _ datasource.DataSource = &VirtInstancesDataSource{}

func NewVirtInstancesDataSource() datasource.DataSource {
	return &VirtInstancesDataSource{}
}

type VirtInstancesDataSource struct {
	client *client.Client
}

type VirtInstancesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type VirtInstancesItemModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Status         types.String `tfsdk:"status"`
	Cpu            types.String `tfsdk:"cpu"`
	Memory         types.Int64  `tfsdk:"memory"`
	Autostart      types.Bool   `tfsdk:"autostart"`
	Environment    types.String `tfsdk:"environment"`
	Image          types.String `tfsdk:"image"`
	UsernsIdmap    types.String `tfsdk:"userns_idmap"`
	Raw            types.String `tfsdk:"raw"`
	VncEnabled     types.Bool   `tfsdk:"vnc_enabled"`
	VncPort        types.Int64  `tfsdk:"vnc_port"`
	VncPassword    types.String `tfsdk:"vnc_password"`
	SecureBoot     types.Bool   `tfsdk:"secure_boot"`
	PrivilegedMode types.Bool   `tfsdk:"privileged_mode"`
	RootDiskSize   types.Int64  `tfsdk:"root_disk_size"`
	RootDiskIoBus  types.String `tfsdk:"root_disk_io_bus"`
	StoragePool    types.String `tfsdk:"storage_pool"`
}

func (d *VirtInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virt_instances"
}

func (d *VirtInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query all instances with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of virt_instances resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable name for the virtual instance.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of virtual instance.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current operational status of the virtual instance.",
						},
						"cpu": schema.StringAttribute{
							Computed:    true,
							Description: "CPU configuration string or `null` for default allocation.",
						},
						"memory": schema.Int64Attribute{
							Computed:    true,
							Description: "Memory allocation in bytes or `null` for default allocation.",
						},
						"autostart": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the instance automatically starts when the host boots.",
						},
						"environment": schema.StringAttribute{
							Computed:    true,
							Description: "Environment variables to set inside the instance.",
						},
						"image": schema.StringAttribute{
							Computed:    true,
							Description: "Image information used to create this instance.",
						},
						"userns_idmap": schema.StringAttribute{
							Computed:    true,
							Description: "User namespace ID mapping configuration for privilege isolation.",
						},
						"raw": schema.StringAttribute{
							Computed:    true,
							Description: "Raw low-level configuration options (advanced use only).",
						},
						"vnc_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether VNC remote access is enabled for the instance.",
						},
						"vnc_port": schema.Int64Attribute{
							Computed:    true,
							Description: "TCP port number for VNC connections or `null` if VNC is disabled.",
						},
						"vnc_password": schema.StringAttribute{
							Computed:    true,
							Description: "Password for VNC access or `null` if no password is set.",
						},
						"secure_boot": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether UEFI Secure Boot is enabled (VMs only) or `null` for containers.",
						},
						"privileged_mode": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the container runs in privileged mode or `null` for VMs.",
						},
						"root_disk_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the root disk in GB or `null` for default size.",
						},
						"root_disk_io_bus": schema.StringAttribute{
							Computed:    true,
							Description: "I/O bus type for the root disk or `null` for default.",
						},
						"storage_pool": schema.StringAttribute{
							Computed:    true,
							Description: "Storage pool in which the root of the instance is located.",
						},
					},
				},
			},
		},
	}
}

func (d *VirtInstancesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VirtInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtInstancesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("virt.instance.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query virt_instances: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]VirtInstancesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := VirtInstancesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["status"]; ok && v != nil {
			itemModel.Status = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cpu"]; ok && v != nil {
			itemModel.Cpu = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["memory"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Memory = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["autostart"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Autostart = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["environment"]; ok && v != nil {
			itemModel.Environment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["image"]; ok && v != nil {
			itemModel.Image = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["userns_idmap"]; ok && v != nil {
			itemModel.UsernsIdmap = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["raw"]; ok && v != nil {
			itemModel.Raw = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vnc_enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.VncEnabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["vnc_port"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.VncPort = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["vnc_password"]; ok && v != nil {
			itemModel.VncPassword = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["secure_boot"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.SecureBoot = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["privileged_mode"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.PrivilegedMode = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["root_disk_size"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.RootDiskSize = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["root_disk_io_bus"]; ok && v != nil {
			itemModel.RootDiskIoBus = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["storage_pool"]; ok && v != nil {
			itemModel.StoragePool = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"autostart":        types.BoolType,
			"cpu":              types.StringType,
			"environment":      types.StringType,
			"id":               types.StringType,
			"image":            types.StringType,
			"memory":           types.Int64Type,
			"name":             types.StringType,
			"privileged_mode":  types.BoolType,
			"raw":              types.StringType,
			"root_disk_io_bus": types.StringType,
			"root_disk_size":   types.Int64Type,
			"secure_boot":      types.BoolType,
			"status":           types.StringType,
			"storage_pool":     types.StringType,
			"type":             types.StringType,
			"userns_idmap":     types.StringType,
			"vnc_enabled":      types.BoolType,
			"vnc_password":     types.StringType,
			"vnc_port":         types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
