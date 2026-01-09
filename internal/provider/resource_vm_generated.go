package provider

import (
	"context"
	"fmt"
	"strconv"

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
	MinMemory types.String `tfsdk:"min_memory"`
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
	EnableSecureBoot types.Bool `tfsdk:"enable_secure_boot"`
}

func NewVmResource() resource.Resource {
	return &VmResource{}
}

func (r *VmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (r *VmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS vm resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"start_on_create": schema.BoolAttribute{
				Optional: true,
				Description: "Start the resource immediately after creation (default: true if not specified)",
			},
			"command_line_args": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"cpu_mode": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"cpu_model": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"vcpus": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"cores": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"threads": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"cpuset": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"nodeset": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enable_cpu_topology_extension": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"pin_vcpus": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"suspend_on_snapshot": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"trusted_platform_module": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"memory": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"min_memory": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"hyperv_enlightenments": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"bootloader": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"bootloader_ovmf": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"autostart": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"hide_from_msr": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"ensure_display_device": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"time": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"shutdown_timeout": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"arch_type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"machine_type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enable_secure_boot": schema.BoolAttribute{
				Required: false,
				Optional: true,
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
	params["name"] = data.Name.ValueString()
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
	params["memory"] = data.Memory.ValueInt64()
	if !data.MinMemory.IsNull() {
		params["min_memory"] = data.MinMemory.ValueString()
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
	if !data.EnableSecureBoot.IsNull() {
		params["enable_secure_boot"] = data.EnableSecureBoot.ValueBool()
	}

	result, err := r.client.Call("vm.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

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
		// Convert string ID to integer for TrueNAS API
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

	_, err := r.client.Call("vm.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state VmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
	params["name"] = data.Name.ValueString()
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
	params["memory"] = data.Memory.ValueInt64()
	if !data.MinMemory.IsNull() {
		params["min_memory"] = data.MinMemory.ValueString()
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
	if !data.EnableSecureBoot.IsNull() {
		params["enable_secure_boot"] = data.EnableSecureBoot.ValueBool()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("vm.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("vm.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
