package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type VmResource struct {
	client *client.Client
}

type VmResourceModel struct {
	ID types.String `tfsdk:"id"`
	StartOnCreate types.Bool `tfsdk:"start_on_create"`
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
	EnableSecureBoot types.Bool `tfsdk:"enable_secure_boot"`
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
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"start_on_create": schema.BoolAttribute{Optional: true, Description: "Start the resource immediately after creation (default: true)"},
			"command_line_args": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Additional command line arguments passed to the VM hypervisor.",
			},
			"cpu_mode": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "CPU virtualization mode.  * `CUSTOM`: Use specified model. * `HOST-MODEL`: Mirror host CPU. * `HOST-",
			},
			"cpu_model": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Specific CPU model to emulate. `null` to use hypervisor default.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Display name of the virtual machine.",
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional description or notes about the virtual machine.",
			},
			"vcpus": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of virtual CPUs allocated to the VM.",
			},
			"cores": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of CPU cores per socket.",
			},
			"threads": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of threads per CPU core.",
			},
			"cpuset": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Set of host CPU cores to pin VM CPUs to. `null` for no pinning.",
			},
			"nodeset": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.",
			},
			"enable_cpu_topology_extension": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to expose detailed CPU topology information to the guest OS.",
			},
			"pin_vcpus": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexib",
			},
			"suspend_on_snapshot": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to suspend the VM when taking snapshots.",
			},
			"trusted_platform_module": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable virtual Trusted Platform Module (TPM) for the VM.",
			},
			"memory": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "Amount of memory allocated to the VM in megabytes.",
			},
			"min_memory": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink    ",
			},
			"hyperv_enlightenments": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable Hyper-V enlightenments for improved Windows guest performance.",
			},
			"bootloader": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility.",
			},
			"bootloader_ovmf": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "OVMF firmware file to use for UEFI boot.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"autostart": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to automatically start the VM when the host system boots.",
			},
			"hide_from_msr": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to hide hypervisor signatures from guest OS MSR access.",
			},
			"ensure_display_device": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to ensure at least one display device is configured for the VM.",
			},
			"time": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time.",
			},
			"shutdown_timeout": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances",
			},
			"arch_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Guest architecture type. `null` to use hypervisor default.",
			},
			"machine_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Virtual machine type/chipset. `null` to use hypervisor default.",
			},
			"uuid": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Unique UUID for the VM. `null` to auto-generate.",
			},
			"enable_secure_boot": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable UEFI Secure Boot for enhanced security.",
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
	if !data.CommandLineArgs.IsNull() {
		params["command_line_args"] = data.CommandLineArgs.ValueString()
	}
	if !data.CpuMode.IsNull() {
		params["cpu_mode"] = data.CpuMode.ValueString()
	}
	if !data.CpuModel.IsNull() {
		params["cpu_model"] = data.CpuModel.ValueString()
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Vcpus.IsNull() {
		params["vcpus"] = data.Vcpus.ValueInt64()
	}
	if !data.Cores.IsNull() {
		params["cores"] = data.Cores.ValueInt64()
	}
	if !data.Threads.IsNull() {
		params["threads"] = data.Threads.ValueInt64()
	}
	if !data.Cpuset.IsNull() {
		params["cpuset"] = data.Cpuset.ValueString()
	}
	if !data.Nodeset.IsNull() {
		params["nodeset"] = data.Nodeset.ValueString()
	}
	if !data.EnableCpuTopologyExtension.IsNull() {
		params["enable_cpu_topology_extension"] = data.EnableCpuTopologyExtension.ValueBool()
	}
	if !data.PinVcpus.IsNull() {
		params["pin_vcpus"] = data.PinVcpus.ValueBool()
	}
	if !data.SuspendOnSnapshot.IsNull() {
		params["suspend_on_snapshot"] = data.SuspendOnSnapshot.ValueBool()
	}
	if !data.TrustedPlatformModule.IsNull() {
		params["trusted_platform_module"] = data.TrustedPlatformModule.ValueBool()
	}
	if !data.Memory.IsNull() {
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.MinMemory.IsNull() {
		params["min_memory"] = data.MinMemory.ValueInt64()
	}
	if !data.HypervEnlightenments.IsNull() {
		params["hyperv_enlightenments"] = data.HypervEnlightenments.ValueBool()
	}
	if !data.Bootloader.IsNull() {
		params["bootloader"] = data.Bootloader.ValueString()
	}
	if !data.BootloaderOvmf.IsNull() {
		params["bootloader_ovmf"] = data.BootloaderOvmf.ValueString()
	}
	if !data.Autostart.IsNull() {
		params["autostart"] = data.Autostart.ValueBool()
	}
	if !data.HideFromMsr.IsNull() {
		params["hide_from_msr"] = data.HideFromMsr.ValueBool()
	}
	if !data.EnsureDisplayDevice.IsNull() {
		params["ensure_display_device"] = data.EnsureDisplayDevice.ValueBool()
	}
	if !data.Time.IsNull() {
		params["time"] = data.Time.ValueString()
	}
	if !data.ShutdownTimeout.IsNull() {
		params["shutdown_timeout"] = data.ShutdownTimeout.ValueInt64()
	}
	if !data.ArchType.IsNull() {
		params["arch_type"] = data.ArchType.ValueString()
	}
	if !data.MachineType.IsNull() {
		params["machine_type"] = data.MachineType.ValueString()
	}
	if !data.Uuid.IsNull() {
		params["uuid"] = data.Uuid.ValueString()
	}
	if !data.EnableSecureBoot.IsNull() {
		params["enable_secure_boot"] = data.EnableSecureBoot.ValueBool()
	}

	result, err := r.client.Call("vm.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create vm: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Handle lifecycle action - start on create if requested
	startOnCreate := true  // default when not specified
	if !data.StartOnCreate.IsNull() {
		startOnCreate = data.StartOnCreate.ValueBool()
	}
	if startOnCreate {
		vmID, err := strconv.Atoi(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
			return
		}
		_, err = r.client.Call("vm.start", vmID)
		if err != nil {
			resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
		}
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
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	result, err := r.client.Call("vm.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read vm: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["command_line_args"]; ok && v != nil {
			data.CommandLineArgs = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cpu_mode"]; ok && v != nil {
			data.CpuMode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cpu_model"]; ok && v != nil {
			data.CpuModel = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["vcpus"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Vcpus = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["cores"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Cores = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["threads"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Threads = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["cpuset"]; ok && v != nil {
			data.Cpuset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["nodeset"]; ok && v != nil {
			data.Nodeset = types.StringValue(fmt.Sprintf("%v", v))
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
			if fv, ok := v.(float64); ok { data.Memory = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["min_memory"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.MinMemory = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["hyperv_enlightenments"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.HypervEnlightenments = types.BoolValue(bv) }
		}
		if v, ok := resultMap["bootloader"]; ok && v != nil {
			data.Bootloader = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["bootloader_ovmf"]; ok && v != nil {
			data.BootloaderOvmf = types.StringValue(fmt.Sprintf("%v", v))
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
			data.Time = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["shutdown_timeout"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.ShutdownTimeout = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["arch_type"]; ok && v != nil {
			data.ArchType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["machine_type"]; ok && v != nil {
			data.MachineType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["uuid"]; ok && v != nil {
			data.Uuid = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enable_secure_boot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.EnableSecureBoot = types.BoolValue(bv) }
		}
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
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	params := map[string]interface{}{}
	if !data.CommandLineArgs.IsNull() {
		params["command_line_args"] = data.CommandLineArgs.ValueString()
	}
	if !data.CpuMode.IsNull() {
		params["cpu_mode"] = data.CpuMode.ValueString()
	}
	if !data.CpuModel.IsNull() {
		params["cpu_model"] = data.CpuModel.ValueString()
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Vcpus.IsNull() {
		params["vcpus"] = data.Vcpus.ValueInt64()
	}
	if !data.Cores.IsNull() {
		params["cores"] = data.Cores.ValueInt64()
	}
	if !data.Threads.IsNull() {
		params["threads"] = data.Threads.ValueInt64()
	}
	if !data.Cpuset.IsNull() {
		params["cpuset"] = data.Cpuset.ValueString()
	}
	if !data.Nodeset.IsNull() {
		params["nodeset"] = data.Nodeset.ValueString()
	}
	if !data.EnableCpuTopologyExtension.IsNull() {
		params["enable_cpu_topology_extension"] = data.EnableCpuTopologyExtension.ValueBool()
	}
	if !data.PinVcpus.IsNull() {
		params["pin_vcpus"] = data.PinVcpus.ValueBool()
	}
	if !data.SuspendOnSnapshot.IsNull() {
		params["suspend_on_snapshot"] = data.SuspendOnSnapshot.ValueBool()
	}
	if !data.TrustedPlatformModule.IsNull() {
		params["trusted_platform_module"] = data.TrustedPlatformModule.ValueBool()
	}
	if !data.Memory.IsNull() {
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.MinMemory.IsNull() {
		params["min_memory"] = data.MinMemory.ValueInt64()
	}
	if !data.HypervEnlightenments.IsNull() {
		params["hyperv_enlightenments"] = data.HypervEnlightenments.ValueBool()
	}
	if !data.Bootloader.IsNull() {
		params["bootloader"] = data.Bootloader.ValueString()
	}
	if !data.Autostart.IsNull() {
		params["autostart"] = data.Autostart.ValueBool()
	}
	if !data.HideFromMsr.IsNull() {
		params["hide_from_msr"] = data.HideFromMsr.ValueBool()
	}
	if !data.EnsureDisplayDevice.IsNull() {
		params["ensure_display_device"] = data.EnsureDisplayDevice.ValueBool()
	}
	if !data.Time.IsNull() {
		params["time"] = data.Time.ValueString()
	}
	if !data.ShutdownTimeout.IsNull() {
		params["shutdown_timeout"] = data.ShutdownTimeout.ValueInt64()
	}
	if !data.ArchType.IsNull() {
		params["arch_type"] = data.ArchType.ValueString()
	}
	if !data.MachineType.IsNull() {
		params["machine_type"] = data.MachineType.ValueString()
	}
	if !data.Uuid.IsNull() {
		params["uuid"] = data.Uuid.ValueString()
	}

	_, err = r.client.Call("vm.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update vm: %s", err))
		return
	}

	data.ID = state.ID
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
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	// Stop VM before deletion if running
	vmID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}
	_, _ = r.client.Call("vm.stop", vmID)  // Ignore errors - VM might already be stopped
	time.Sleep(2 * time.Second)  // Wait for VM to stop

	_, err = r.client.Call("vm.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete vm: %s", err))
		return
	}
}
