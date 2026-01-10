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

var _ datasource.DataSource = &VmsDataSource{}

func NewVmsDataSource() datasource.DataSource {
	return &VmsDataSource{}
}

type VmsDataSource struct {
	client *client.Client
}

type VmsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type VmsItemModel struct {
	ID types.String `tfsdk:"id"`
	CommandLineArgs types.String `tfsdk:"command_line_args"`
	CpuMode types.String `tfsdk:"cpu_mode"`
	CpuModel types.String `tfsdk:"cpu_model"`
	Name types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Vcpus types.Int64 `tfsdk:"vcpus"`
	Cores types.Int64 `tfsdk:"cores"`
	Threads types.Int64 `tfsdk:"threads"`
	Cpuset types.String `tfsdk:"cpuset"`
	Nodeset types.String `tfsdk:"nodeset"`
	EnableCpuTopologyExtension types.Bool `tfsdk:"enable_cpu_topology_extension"`
	PinVcpus types.Bool `tfsdk:"pin_vcpus"`
	SuspendOnSnapshot types.Bool `tfsdk:"suspend_on_snapshot"`
	TrustedPlatformModule types.Bool `tfsdk:"trusted_platform_module"`
	Memory types.Int64 `tfsdk:"memory"`
	MinMemory types.Int64 `tfsdk:"min_memory"`
	HypervEnlightenments types.Bool `tfsdk:"hyperv_enlightenments"`
	Bootloader types.String `tfsdk:"bootloader"`
	BootloaderOvmf types.String `tfsdk:"bootloader_ovmf"`
	Autostart types.Bool `tfsdk:"autostart"`
	HideFromMsr types.Bool `tfsdk:"hide_from_msr"`
	EnsureDisplayDevice types.Bool `tfsdk:"ensure_display_device"`
	Time types.String `tfsdk:"time"`
	ShutdownTimeout types.Int64 `tfsdk:"shutdown_timeout"`
	ArchType types.String `tfsdk:"arch_type"`
	MachineType types.String `tfsdk:"machine_type"`
	Uuid types.String `tfsdk:"uuid"`
	DisplayAvailable types.Bool `tfsdk:"display_available"`
	Status types.String `tfsdk:"status"`
	EnableSecureBoot types.Bool `tfsdk:"enable_secure_boot"`
}

func (d *VmsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vms"
}

func (d *VmsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query vms resources",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of vms resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"command_line_args": schema.StringAttribute{
				Computed: true,
				Description: "Additional command line arguments passed to the VM hypervisor.",
			},
			"cpu_mode": schema.StringAttribute{
				Computed: true,
				Description: "CPU virtualization mode.  * `CUSTOM`: Use specified model. * `HOST-MODEL`: Mirror host CPU. * `HOST-",
			},
			"cpu_model": schema.StringAttribute{
				Computed: true,
				Description: "Specific CPU model to emulate. `null` to use hypervisor default.",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Display name of the virtual machine.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Optional description or notes about the virtual machine.",
			},
			"vcpus": schema.Int64Attribute{
				Computed: true,
				Description: "Number of virtual CPUs allocated to the VM.",
			},
			"cores": schema.Int64Attribute{
				Computed: true,
				Description: "Number of CPU cores per socket.",
			},
			"threads": schema.Int64Attribute{
				Computed: true,
				Description: "Number of threads per CPU core.",
			},
			"cpuset": schema.StringAttribute{
				Computed: true,
				Description: "Set of host CPU cores to pin VM CPUs to. `null` for no pinning.",
			},
			"nodeset": schema.StringAttribute{
				Computed: true,
				Description: "Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.",
			},
			"enable_cpu_topology_extension": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to expose detailed CPU topology information to the guest OS.",
			},
			"pin_vcpus": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexib",
			},
			"suspend_on_snapshot": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to suspend the VM when taking snapshots.",
			},
			"trusted_platform_module": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to enable virtual Trusted Platform Module (TPM) for the VM.",
			},
			"memory": schema.Int64Attribute{
				Computed: true,
				Description: "Amount of memory allocated to the VM in megabytes.",
			},
			"min_memory": schema.Int64Attribute{
				Computed: true,
				Description: "Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink    ",
			},
			"hyperv_enlightenments": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to enable Hyper-V enlightenments for improved Windows guest performance.",
			},
			"bootloader": schema.StringAttribute{
				Computed: true,
				Description: "Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility.",
			},
			"bootloader_ovmf": schema.StringAttribute{
				Computed: true,
				Description: "OVMF firmware file to use for UEFI boot.",
			},
			"autostart": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to automatically start the VM when the host system boots.",
			},
			"hide_from_msr": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to hide hypervisor signatures from guest OS MSR access.",
			},
			"ensure_display_device": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to ensure at least one display device is configured for the VM.",
			},
			"time": schema.StringAttribute{
				Computed: true,
				Description: "Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time.",
			},
			"shutdown_timeout": schema.Int64Attribute{
				Computed: true,
				Description: "Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances",
			},
			"arch_type": schema.StringAttribute{
				Computed: true,
				Description: "Guest architecture type. `null` to use hypervisor default.",
			},
			"machine_type": schema.StringAttribute{
				Computed: true,
				Description: "Virtual machine type/chipset. `null` to use hypervisor default.",
			},
			"uuid": schema.StringAttribute{
				Computed: true,
				Description: "Unique UUID for the VM. `null` to auto-generate.",
			},
			"display_available": schema.BoolAttribute{
				Computed: true,
				Description: "Whether at least one display device is available for this VM.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Current runtime status information for the VM.",
			},
			"enable_secure_boot": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to enable UEFI Secure Boot for enhanced security.",
			},
					},
				},
			},
		},
	}
}

func (d *VmsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VmsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("vm.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query vms: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]VmsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := VmsItemModel{}
		if v, ok := resultMap["command_line_args"]; ok && v != nil {
			itemModel.CommandLineArgs = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cpu_mode"]; ok && v != nil {
			itemModel.CpuMode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cpu_model"]; ok && v != nil {
			itemModel.CpuModel = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vcpus"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Vcpus = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["cores"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Cores = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["threads"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Threads = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["cpuset"]; ok && v != nil {
			itemModel.Cpuset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["nodeset"]; ok && v != nil {
			itemModel.Nodeset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enable_cpu_topology_extension"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.EnableCpuTopologyExtension = types.BoolValue(bv) }
		}
		if v, ok := resultMap["pin_vcpus"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.PinVcpus = types.BoolValue(bv) }
		}
		if v, ok := resultMap["suspend_on_snapshot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.SuspendOnSnapshot = types.BoolValue(bv) }
		}
		if v, ok := resultMap["trusted_platform_module"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.TrustedPlatformModule = types.BoolValue(bv) }
		}
		if v, ok := resultMap["memory"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Memory = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["min_memory"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.MinMemory = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["hyperv_enlightenments"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.HypervEnlightenments = types.BoolValue(bv) }
		}
		if v, ok := resultMap["bootloader"]; ok && v != nil {
			itemModel.Bootloader = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["bootloader_ovmf"]; ok && v != nil {
			itemModel.BootloaderOvmf = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["autostart"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Autostart = types.BoolValue(bv) }
		}
		if v, ok := resultMap["hide_from_msr"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.HideFromMsr = types.BoolValue(bv) }
		}
		if v, ok := resultMap["ensure_display_device"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.EnsureDisplayDevice = types.BoolValue(bv) }
		}
		if v, ok := resultMap["time"]; ok && v != nil {
			itemModel.Time = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["shutdown_timeout"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.ShutdownTimeout = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["arch_type"]; ok && v != nil {
			itemModel.ArchType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["machine_type"]; ok && v != nil {
			itemModel.MachineType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["uuid"]; ok && v != nil {
			itemModel.Uuid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["display_available"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.DisplayAvailable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["status"]; ok && v != nil {
			itemModel.Status = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enable_secure_boot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.EnableSecureBoot = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"arch_type": types.StringType,
			"autostart": types.BoolType,
			"bootloader": types.StringType,
			"bootloader_ovmf": types.StringType,
			"command_line_args": types.StringType,
			"cores": types.Int64Type,
			"cpu_mode": types.StringType,
			"cpu_model": types.StringType,
			"cpuset": types.StringType,
			"description": types.StringType,
			"display_available": types.BoolType,
			"enable_cpu_topology_extension": types.BoolType,
			"enable_secure_boot": types.BoolType,
			"ensure_display_device": types.BoolType,
			"hide_from_msr": types.BoolType,
			"hyperv_enlightenments": types.BoolType,
			"id": types.StringType,
			"machine_type": types.StringType,
			"memory": types.Int64Type,
			"min_memory": types.Int64Type,
			"name": types.StringType,
			"nodeset": types.StringType,
			"pin_vcpus": types.BoolType,
			"shutdown_timeout": types.Int64Type,
			"status": types.StringType,
			"suspend_on_snapshot": types.BoolType,
			"threads": types.Int64Type,
			"time": types.StringType,
			"trusted_platform_module": types.BoolType,
			"uuid": types.StringType,
			"vcpus": types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
