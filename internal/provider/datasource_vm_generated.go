package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &VmDataSource{}

func NewVmDataSource() datasource.DataSource {
	return &VmDataSource{}
}

type VmDataSource struct {
	client *client.Client
}

type VmDataSourceModel struct {
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
	Devices types.List `tfsdk:"devices"`
	DisplayAvailable types.Bool `tfsdk:"display_available"`
	Status types.String `tfsdk:"status"`
	EnableSecureBoot types.Bool `tfsdk:"enable_secure_boot"`
}

func (d *VmDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (d *VmDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
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
			"devices": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "Array of virtual devices attached to this VM.",
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
	}
}

func (d *VmDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VmDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("vm.get_instance", func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read vm: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["command_line_args"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.CommandLineArgs = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.CommandLineArgs = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.CommandLineArgs = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["cpu_mode"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.CpuMode = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.CpuMode = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.CpuMode = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["cpu_model"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.CpuModel = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.CpuModel = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.CpuModel = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["name"]; ok && v != nil {
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
		if v, ok := resultMap["description"]; ok && v != nil {
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
		if v, ok := resultMap["vcpus"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				data.Vcpus = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Vcpus = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["cores"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				data.Cores = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Cores = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["threads"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				data.Threads = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Threads = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["cpuset"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Cpuset = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Cpuset = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Cpuset = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["nodeset"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Nodeset = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Nodeset = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Nodeset = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["enable_cpu_topology_extension"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.EnableCpuTopologyExtension = types.BoolValue(bv) }
		}
		if v, ok := resultMap["pin_vcpus"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.PinVcpus = types.BoolValue(bv) }
		}
		if v, ok := resultMap["suspend_on_snapshot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.SuspendOnSnapshot = types.BoolValue(bv) }
		}
		if v, ok := resultMap["trusted_platform_module"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.TrustedPlatformModule = types.BoolValue(bv) }
		}
		if v, ok := resultMap["memory"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				data.Memory = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Memory = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["min_memory"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				data.MinMemory = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.MinMemory = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["hyperv_enlightenments"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.HypervEnlightenments = types.BoolValue(bv) }
		}
		if v, ok := resultMap["bootloader"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Bootloader = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Bootloader = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Bootloader = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["bootloader_ovmf"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.BootloaderOvmf = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.BootloaderOvmf = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.BootloaderOvmf = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["autostart"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Autostart = types.BoolValue(bv) }
		}
		if v, ok := resultMap["hide_from_msr"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.HideFromMsr = types.BoolValue(bv) }
		}
		if v, ok := resultMap["ensure_display_device"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.EnsureDisplayDevice = types.BoolValue(bv) }
		}
		if v, ok := resultMap["time"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Time = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Time = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Time = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["shutdown_timeout"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				data.ShutdownTimeout = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.ShutdownTimeout = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["arch_type"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.ArchType = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ArchType = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ArchType = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["machine_type"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.MachineType = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.MachineType = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.MachineType = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["uuid"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Uuid = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Uuid = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Uuid = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["devices"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Devices, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["display_available"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.DisplayAvailable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["status"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Status = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Status = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Status = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["enable_secure_boot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.EnableSecureBoot = types.BoolValue(bv) }
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
