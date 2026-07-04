package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"time"
)

type VmResource struct {
	client *client.Client
}

type VmResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	StartOnCreate              types.Bool   `tfsdk:"start_on_create"`
	CommandLineArgs            types.String `tfsdk:"command_line_args"`
	CpuMode                    types.String `tfsdk:"cpu_mode"`
	CpuModel                   types.String `tfsdk:"cpu_model"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	Vcpus                      types.Int64  `tfsdk:"vcpus"`
	Cores                      types.Int64  `tfsdk:"cores"`
	Threads                    types.Int64  `tfsdk:"threads"`
	Cpuset                     types.String `tfsdk:"cpuset"`
	Nodeset                    types.String `tfsdk:"nodeset"`
	EnableCpuTopologyExtension types.Bool   `tfsdk:"enable_cpu_topology_extension"`
	PinVcpus                   types.Bool   `tfsdk:"pin_vcpus"`
	SuspendOnSnapshot          types.Bool   `tfsdk:"suspend_on_snapshot"`
	TrustedPlatformModule      types.Bool   `tfsdk:"trusted_platform_module"`
	Memory                     types.Int64  `tfsdk:"memory"`
	MinMemory                  types.Int64  `tfsdk:"min_memory"`
	HypervEnlightenments       types.Bool   `tfsdk:"hyperv_enlightenments"`
	Bootloader                 types.String `tfsdk:"bootloader"`
	BootloaderOvmf             types.String `tfsdk:"bootloader_ovmf"`
	Autostart                  types.Bool   `tfsdk:"autostart"`
	HideFromMsr                types.Bool   `tfsdk:"hide_from_msr"`
	EnsureDisplayDevice        types.Bool   `tfsdk:"ensure_display_device"`
	Time                       types.String `tfsdk:"time"`
	ShutdownTimeout            types.Int64  `tfsdk:"shutdown_timeout"`
	ArchType                   types.String `tfsdk:"arch_type"`
	MachineType                types.String `tfsdk:"machine_type"`
	Uuid                       types.String `tfsdk:"uuid"`
	EnableSecureBoot           types.Bool   `tfsdk:"enable_secure_boot"`
}

func NewVmResource() resource.Resource {
	return &VmResource{}
}

func (r *VmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (r *VmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *VmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Virtual Machine (VM).",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"start_on_create": schema.BoolAttribute{Optional: true, Description: "Start the resource immediately after creation (default: false)"},
			"command_line_args": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Additional command line arguments passed to the VM hypervisor.",
			},
			"cpu_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "CPU virtualization mode.  * `CUSTOM`: Use specified model. * `HOST-MODEL`: Mirror host CPU. * `HOST-",
			},
			"cpu_model": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Specific CPU model to emulate. `null` to use hypervisor default.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Display name of the virtual machine.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional description or notes about the virtual machine.",
			},
			"vcpus": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of virtual CPUs allocated to the VM.",
			},
			"cores": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of CPU cores per socket.",
			},
			"threads": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of threads per CPU core.",
			},
			"cpuset": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set of host CPU cores to pin VM CPUs to. `null` for no pinning.",
			},
			"nodeset": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.",
			},
			"enable_cpu_topology_extension": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to expose detailed CPU topology information to the guest OS.",
			},
			"pin_vcpus": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexib",
			},
			"suspend_on_snapshot": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to suspend the VM when taking snapshots.",
			},
			"trusted_platform_module": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable virtual Trusted Platform Module (TPM) for the VM.",
			},
			"memory": schema.Int64Attribute{
				Required:    true,
				Optional:    false,
				Description: "Amount of memory allocated to the VM in megabytes.",
			},
			"min_memory": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink    ",
			},
			"hyperv_enlightenments": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable Hyper-V enlightenments for improved Windows guest performance.",
			},
			"bootloader": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility.",
			},
			"bootloader_ovmf": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "OVMF firmware file to use for UEFI boot.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"autostart": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to automatically start the VM when the host system boots.",
			},
			"hide_from_msr": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to hide hypervisor signatures from guest OS MSR access.",
			},
			"ensure_display_device": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to ensure at least one display device is configured for the VM.",
			},
			"time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time.",
			},
			"shutdown_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances",
			},
			"arch_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Guest architecture type. `null` to use hypervisor default.",
			},
			"machine_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Virtual machine type/chipset. `null` to use hypervisor default.",
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Unique UUID for the VM. `null` to auto-generate.",
			},
			"enable_secure_boot": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether to enable UEFI Secure Boot for enhanced security.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *VmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
		return
	}
	r.client = client
}

func (r *VmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.CommandLineArgs.IsNull() && !data.CommandLineArgs.IsUnknown() {
		params["command_line_args"] = data.CommandLineArgs.ValueString()
	}
	if !data.CpuMode.IsNull() && !data.CpuMode.IsUnknown() {
		params["cpu_mode"] = data.CpuMode.ValueString()
	}
	if !data.CpuModel.IsNull() && !data.CpuModel.IsUnknown() {
		params["cpu_model"] = data.CpuModel.ValueString()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Vcpus.IsNull() && !data.Vcpus.IsUnknown() {
		params["vcpus"] = data.Vcpus.ValueInt64()
	}
	if !data.Cores.IsNull() && !data.Cores.IsUnknown() {
		params["cores"] = data.Cores.ValueInt64()
	}
	if !data.Threads.IsNull() && !data.Threads.IsUnknown() {
		params["threads"] = data.Threads.ValueInt64()
	}
	if !data.Cpuset.IsNull() && !data.Cpuset.IsUnknown() {
		params["cpuset"] = data.Cpuset.ValueString()
	}
	if !data.Nodeset.IsNull() && !data.Nodeset.IsUnknown() {
		params["nodeset"] = data.Nodeset.ValueString()
	}
	if !data.EnableCpuTopologyExtension.IsNull() && !data.EnableCpuTopologyExtension.IsUnknown() {
		params["enable_cpu_topology_extension"] = data.EnableCpuTopologyExtension.ValueBool()
	}
	if !data.PinVcpus.IsNull() && !data.PinVcpus.IsUnknown() {
		params["pin_vcpus"] = data.PinVcpus.ValueBool()
	}
	if !data.SuspendOnSnapshot.IsNull() && !data.SuspendOnSnapshot.IsUnknown() {
		params["suspend_on_snapshot"] = data.SuspendOnSnapshot.ValueBool()
	}
	if !data.TrustedPlatformModule.IsNull() && !data.TrustedPlatformModule.IsUnknown() {
		params["trusted_platform_module"] = data.TrustedPlatformModule.ValueBool()
	}
	if !data.Memory.IsNull() && !data.Memory.IsUnknown() {
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.MinMemory.IsNull() && !data.MinMemory.IsUnknown() {
		params["min_memory"] = data.MinMemory.ValueInt64()
	}
	if !data.HypervEnlightenments.IsNull() && !data.HypervEnlightenments.IsUnknown() {
		params["hyperv_enlightenments"] = data.HypervEnlightenments.ValueBool()
	}
	if !data.Bootloader.IsNull() && !data.Bootloader.IsUnknown() {
		params["bootloader"] = data.Bootloader.ValueString()
	}
	if !data.BootloaderOvmf.IsNull() && !data.BootloaderOvmf.IsUnknown() {
		params["bootloader_ovmf"] = data.BootloaderOvmf.ValueString()
	}
	if !data.Autostart.IsNull() && !data.Autostart.IsUnknown() {
		params["autostart"] = data.Autostart.ValueBool()
	}
	if !data.HideFromMsr.IsNull() && !data.HideFromMsr.IsUnknown() {
		params["hide_from_msr"] = data.HideFromMsr.ValueBool()
	}
	if !data.EnsureDisplayDevice.IsNull() && !data.EnsureDisplayDevice.IsUnknown() {
		params["ensure_display_device"] = data.EnsureDisplayDevice.ValueBool()
	}
	if !data.Time.IsNull() && !data.Time.IsUnknown() {
		params["time"] = data.Time.ValueString()
	}
	if !data.ShutdownTimeout.IsNull() && !data.ShutdownTimeout.IsUnknown() {
		params["shutdown_timeout"] = data.ShutdownTimeout.ValueInt64()
	}
	if !data.ArchType.IsNull() && !data.ArchType.IsUnknown() {
		params["arch_type"] = data.ArchType.ValueString()
	}
	if !data.MachineType.IsNull() && !data.MachineType.IsUnknown() {
		params["machine_type"] = data.MachineType.ValueString()
	}
	if !data.Uuid.IsNull() && !data.Uuid.IsUnknown() {
		params["uuid"] = data.Uuid.ValueString()
	}
	if !data.EnableSecureBoot.IsNull() && !data.EnableSecureBoot.IsUnknown() {
		params["enable_secure_boot"] = data.EnableSecureBoot.ValueBool()
	}

	result, err := r.client.Call("vm.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create vm: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	startOnCreate := false
	if !data.StartOnCreate.IsNull() {
		startOnCreate = data.StartOnCreate.ValueBool()
	}
	if startOnCreate {
		_, err = r.client.Call("vm.start", func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }())
		if err != nil {
			resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
		}
	}

	// Read back to populate computed fields
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("vm.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back vm: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["cpu_model"]; ok {
		if v == nil {
			data.CpuModel = types.StringNull()
		} else {
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
	if v, ok := resultMap["cpuset"]; ok {
		if v == nil {
			data.Cpuset = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["nodeset"]; ok {
		if v == nil {
			data.Nodeset = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["memory"]; ok {
		switch val := v.(type) {
		case float64:
			data.Memory = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Memory = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["min_memory"]; ok {
		if v == nil {
			data.MinMemory = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.MinMemory = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.MinMemory = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["bootloader_ovmf"]; ok {
		if v == nil {
			data.BootloaderOvmf = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["arch_type"]; ok {
		if v == nil {
			data.ArchType = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["machine_type"]; ok {
		if v == nil {
			data.MachineType = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["uuid"]; ok {
		if v == nil {
			data.Uuid = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["enable_secure_boot"]; ok {
		if bv, ok := v.(bool); ok {
			data.EnableSecureBoot = types.BoolValue(bv)
		}
	}
	if data.CommandLineArgs.IsUnknown() {
		data.CommandLineArgs = types.StringNull()
	}
	if data.CpuMode.IsUnknown() {
		data.CpuMode = types.StringNull()
	}
	if data.CpuModel.IsUnknown() {
		data.CpuModel = types.StringNull()
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Vcpus.IsUnknown() {
		data.Vcpus = types.Int64Null()
	}
	if data.Cores.IsUnknown() {
		data.Cores = types.Int64Null()
	}
	if data.Threads.IsUnknown() {
		data.Threads = types.Int64Null()
	}
	if data.Cpuset.IsUnknown() {
		data.Cpuset = types.StringNull()
	}
	if data.Nodeset.IsUnknown() {
		data.Nodeset = types.StringNull()
	}
	if data.EnableCpuTopologyExtension.IsUnknown() {
		data.EnableCpuTopologyExtension = types.BoolNull()
	}
	if data.PinVcpus.IsUnknown() {
		data.PinVcpus = types.BoolNull()
	}
	if data.SuspendOnSnapshot.IsUnknown() {
		data.SuspendOnSnapshot = types.BoolNull()
	}
	if data.TrustedPlatformModule.IsUnknown() {
		data.TrustedPlatformModule = types.BoolNull()
	}
	if data.MinMemory.IsUnknown() {
		data.MinMemory = types.Int64Null()
	}
	if data.HypervEnlightenments.IsUnknown() {
		data.HypervEnlightenments = types.BoolNull()
	}
	if data.Bootloader.IsUnknown() {
		data.Bootloader = types.StringNull()
	}
	if data.BootloaderOvmf.IsUnknown() {
		data.BootloaderOvmf = types.StringNull()
	}
	if data.Autostart.IsUnknown() {
		data.Autostart = types.BoolNull()
	}
	if data.HideFromMsr.IsUnknown() {
		data.HideFromMsr = types.BoolNull()
	}
	if data.EnsureDisplayDevice.IsUnknown() {
		data.EnsureDisplayDevice = types.BoolNull()
	}
	if data.Time.IsUnknown() {
		data.Time = types.StringNull()
	}
	if data.ShutdownTimeout.IsUnknown() {
		data.ShutdownTimeout = types.Int64Null()
	}
	if data.ArchType.IsUnknown() {
		data.ArchType = types.StringNull()
	}
	if data.MachineType.IsUnknown() {
		data.MachineType = types.StringNull()
	}
	if data.Uuid.IsUnknown() {
		data.Uuid = types.StringNull()
	}
	if data.EnableSecureBoot.IsUnknown() {
		data.EnableSecureBoot = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("vm.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read vm: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["cpu_model"]; ok {
		if v == nil {
			data.CpuModel = types.StringNull()
		} else {
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
	if v, ok := resultMap["cpuset"]; ok {
		if v == nil {
			data.Cpuset = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["nodeset"]; ok {
		if v == nil {
			data.Nodeset = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["memory"]; ok {
		switch val := v.(type) {
		case float64:
			data.Memory = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Memory = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["min_memory"]; ok {
		if v == nil {
			data.MinMemory = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.MinMemory = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.MinMemory = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["bootloader_ovmf"]; ok {
		if v == nil {
			data.BootloaderOvmf = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["arch_type"]; ok {
		if v == nil {
			data.ArchType = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["machine_type"]; ok {
		if v == nil {
			data.MachineType = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["uuid"]; ok {
		if v == nil {
			data.Uuid = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["enable_secure_boot"]; ok {
		if bv, ok := v.(bool); ok {
			data.EnableSecureBoot = types.BoolValue(bv)
		}
	}
	if data.CommandLineArgs.IsUnknown() {
		data.CommandLineArgs = types.StringNull()
	}
	if data.CpuMode.IsUnknown() {
		data.CpuMode = types.StringNull()
	}
	if data.CpuModel.IsUnknown() {
		data.CpuModel = types.StringNull()
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Vcpus.IsUnknown() {
		data.Vcpus = types.Int64Null()
	}
	if data.Cores.IsUnknown() {
		data.Cores = types.Int64Null()
	}
	if data.Threads.IsUnknown() {
		data.Threads = types.Int64Null()
	}
	if data.Cpuset.IsUnknown() {
		data.Cpuset = types.StringNull()
	}
	if data.Nodeset.IsUnknown() {
		data.Nodeset = types.StringNull()
	}
	if data.EnableCpuTopologyExtension.IsUnknown() {
		data.EnableCpuTopologyExtension = types.BoolNull()
	}
	if data.PinVcpus.IsUnknown() {
		data.PinVcpus = types.BoolNull()
	}
	if data.SuspendOnSnapshot.IsUnknown() {
		data.SuspendOnSnapshot = types.BoolNull()
	}
	if data.TrustedPlatformModule.IsUnknown() {
		data.TrustedPlatformModule = types.BoolNull()
	}
	if data.MinMemory.IsUnknown() {
		data.MinMemory = types.Int64Null()
	}
	if data.HypervEnlightenments.IsUnknown() {
		data.HypervEnlightenments = types.BoolNull()
	}
	if data.Bootloader.IsUnknown() {
		data.Bootloader = types.StringNull()
	}
	if data.BootloaderOvmf.IsUnknown() {
		data.BootloaderOvmf = types.StringNull()
	}
	if data.Autostart.IsUnknown() {
		data.Autostart = types.BoolNull()
	}
	if data.HideFromMsr.IsUnknown() {
		data.HideFromMsr = types.BoolNull()
	}
	if data.EnsureDisplayDevice.IsUnknown() {
		data.EnsureDisplayDevice = types.BoolNull()
	}
	if data.Time.IsUnknown() {
		data.Time = types.StringNull()
	}
	if data.ShutdownTimeout.IsUnknown() {
		data.ShutdownTimeout = types.Int64Null()
	}
	if data.ArchType.IsUnknown() {
		data.ArchType = types.StringNull()
	}
	if data.MachineType.IsUnknown() {
		data.MachineType = types.StringNull()
	}
	if data.Uuid.IsUnknown() {
		data.Uuid = types.StringNull()
	}
	if data.EnableSecureBoot.IsUnknown() {
		data.EnableSecureBoot = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.CommandLineArgs.IsNull() && !data.CommandLineArgs.IsUnknown() {
		params["command_line_args"] = data.CommandLineArgs.ValueString()
	}
	if !data.CpuMode.IsNull() && !data.CpuMode.IsUnknown() {
		params["cpu_mode"] = data.CpuMode.ValueString()
	}
	if !data.CpuModel.IsNull() && !data.CpuModel.IsUnknown() {
		params["cpu_model"] = data.CpuModel.ValueString()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Vcpus.IsNull() && !data.Vcpus.IsUnknown() {
		params["vcpus"] = data.Vcpus.ValueInt64()
	}
	if !data.Cores.IsNull() && !data.Cores.IsUnknown() {
		params["cores"] = data.Cores.ValueInt64()
	}
	if !data.Threads.IsNull() && !data.Threads.IsUnknown() {
		params["threads"] = data.Threads.ValueInt64()
	}
	if !data.Cpuset.IsNull() && !data.Cpuset.IsUnknown() {
		params["cpuset"] = data.Cpuset.ValueString()
	}
	if !data.Nodeset.IsNull() && !data.Nodeset.IsUnknown() {
		params["nodeset"] = data.Nodeset.ValueString()
	}
	if !data.EnableCpuTopologyExtension.IsNull() && !data.EnableCpuTopologyExtension.IsUnknown() {
		params["enable_cpu_topology_extension"] = data.EnableCpuTopologyExtension.ValueBool()
	}
	if !data.PinVcpus.IsNull() && !data.PinVcpus.IsUnknown() {
		params["pin_vcpus"] = data.PinVcpus.ValueBool()
	}
	if !data.SuspendOnSnapshot.IsNull() && !data.SuspendOnSnapshot.IsUnknown() {
		params["suspend_on_snapshot"] = data.SuspendOnSnapshot.ValueBool()
	}
	if !data.TrustedPlatformModule.IsNull() && !data.TrustedPlatformModule.IsUnknown() {
		params["trusted_platform_module"] = data.TrustedPlatformModule.ValueBool()
	}
	if !data.Memory.IsNull() && !data.Memory.IsUnknown() {
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.MinMemory.IsNull() && !data.MinMemory.IsUnknown() {
		params["min_memory"] = data.MinMemory.ValueInt64()
	}
	if !data.HypervEnlightenments.IsNull() && !data.HypervEnlightenments.IsUnknown() {
		params["hyperv_enlightenments"] = data.HypervEnlightenments.ValueBool()
	}
	if !data.Bootloader.IsNull() && !data.Bootloader.IsUnknown() {
		params["bootloader"] = data.Bootloader.ValueString()
	}
	if !data.Autostart.IsNull() && !data.Autostart.IsUnknown() {
		params["autostart"] = data.Autostart.ValueBool()
	}
	if !data.HideFromMsr.IsNull() && !data.HideFromMsr.IsUnknown() {
		params["hide_from_msr"] = data.HideFromMsr.ValueBool()
	}
	if !data.EnsureDisplayDevice.IsNull() && !data.EnsureDisplayDevice.IsUnknown() {
		params["ensure_display_device"] = data.EnsureDisplayDevice.ValueBool()
	}
	if !data.Time.IsNull() && !data.Time.IsUnknown() {
		params["time"] = data.Time.ValueString()
	}
	if !data.ShutdownTimeout.IsNull() && !data.ShutdownTimeout.IsUnknown() {
		params["shutdown_timeout"] = data.ShutdownTimeout.ValueInt64()
	}
	if !data.ArchType.IsNull() && !data.ArchType.IsUnknown() {
		params["arch_type"] = data.ArchType.ValueString()
	}
	if !data.MachineType.IsNull() && !data.MachineType.IsUnknown() {
		params["machine_type"] = data.MachineType.ValueString()
	}
	if !data.Uuid.IsNull() && !data.Uuid.IsUnknown() {
		params["uuid"] = data.Uuid.ValueString()
	}

	_, err = r.client.Call("vm.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update vm: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("vm.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back vm: %s", readErr))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["cpu_model"]; ok {
		if v == nil {
			data.CpuModel = types.StringNull()
		} else {
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
	if v, ok := resultMap["cpuset"]; ok {
		if v == nil {
			data.Cpuset = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["nodeset"]; ok {
		if v == nil {
			data.Nodeset = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["memory"]; ok {
		switch val := v.(type) {
		case float64:
			data.Memory = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Memory = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["min_memory"]; ok {
		if v == nil {
			data.MinMemory = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.MinMemory = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.MinMemory = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["bootloader_ovmf"]; ok {
		if v == nil {
			data.BootloaderOvmf = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["arch_type"]; ok {
		if v == nil {
			data.ArchType = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["machine_type"]; ok {
		if v == nil {
			data.MachineType = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["uuid"]; ok {
		if v == nil {
			data.Uuid = types.StringNull()
		} else {
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
	}
	if v, ok := resultMap["enable_secure_boot"]; ok {
		if bv, ok := v.(bool); ok {
			data.EnableSecureBoot = types.BoolValue(bv)
		}
	}
	if data.CommandLineArgs.IsUnknown() {
		data.CommandLineArgs = types.StringNull()
	}
	if data.CpuMode.IsUnknown() {
		data.CpuMode = types.StringNull()
	}
	if data.CpuModel.IsUnknown() {
		data.CpuModel = types.StringNull()
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
	if data.Vcpus.IsUnknown() {
		data.Vcpus = types.Int64Null()
	}
	if data.Cores.IsUnknown() {
		data.Cores = types.Int64Null()
	}
	if data.Threads.IsUnknown() {
		data.Threads = types.Int64Null()
	}
	if data.Cpuset.IsUnknown() {
		data.Cpuset = types.StringNull()
	}
	if data.Nodeset.IsUnknown() {
		data.Nodeset = types.StringNull()
	}
	if data.EnableCpuTopologyExtension.IsUnknown() {
		data.EnableCpuTopologyExtension = types.BoolNull()
	}
	if data.PinVcpus.IsUnknown() {
		data.PinVcpus = types.BoolNull()
	}
	if data.SuspendOnSnapshot.IsUnknown() {
		data.SuspendOnSnapshot = types.BoolNull()
	}
	if data.TrustedPlatformModule.IsUnknown() {
		data.TrustedPlatformModule = types.BoolNull()
	}
	if data.MinMemory.IsUnknown() {
		data.MinMemory = types.Int64Null()
	}
	if data.HypervEnlightenments.IsUnknown() {
		data.HypervEnlightenments = types.BoolNull()
	}
	if data.Bootloader.IsUnknown() {
		data.Bootloader = types.StringNull()
	}
	if data.BootloaderOvmf.IsUnknown() {
		data.BootloaderOvmf = types.StringNull()
	}
	if data.Autostart.IsUnknown() {
		data.Autostart = types.BoolNull()
	}
	if data.HideFromMsr.IsUnknown() {
		data.HideFromMsr = types.BoolNull()
	}
	if data.EnsureDisplayDevice.IsUnknown() {
		data.EnsureDisplayDevice = types.BoolNull()
	}
	if data.Time.IsUnknown() {
		data.Time = types.StringNull()
	}
	if data.ShutdownTimeout.IsUnknown() {
		data.ShutdownTimeout = types.Int64Null()
	}
	if data.ArchType.IsUnknown() {
		data.ArchType = types.StringNull()
	}
	if data.MachineType.IsUnknown() {
		data.MachineType = types.StringNull()
	}
	if data.Uuid.IsUnknown() {
		data.Uuid = types.StringNull()
	}
	if data.EnableSecureBoot.IsUnknown() {
		data.EnableSecureBoot = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	id = []interface{}{id, map[string]interface{}{"force": true, "zvols": true}}

	_, _ = r.client.Call("vm.stop", func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }())
	time.Sleep(2 * time.Second)

	_, err = r.client.Call("vm.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete vm: %s", err))
		return
	}
}
