package provider

import (
	"context"
	"fmt"

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
				Computed: true,
				Description: "Automatically start after creation (default: true)",
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

	params := map[string]interface{}{
		"command_line_args": data.CommandLineArgs.ValueString(),
		"cpu_mode": data.CpuMode.ValueString(),
		"cpu_model": data.CpuModel.ValueString(),
		"name": data.Name.ValueString(),
		"description": data.Description.ValueString(),
		"vcpus": data.Vcpus.ValueInt64(),
		"cores": data.Cores.ValueInt64(),
		"threads": data.Threads.ValueInt64(),
		"cpuset": data.Cpuset.ValueString(),
		"nodeset": data.Nodeset.ValueString(),
		"enable_cpu_topology_extension": data.EnableCpuTopologyExtension.ValueBool(),
		"pin_vcpus": data.PinVcpus.ValueBool(),
		"suspend_on_snapshot": data.SuspendOnSnapshot.ValueBool(),
		"trusted_platform_module": data.TrustedPlatformModule.ValueBool(),
		"memory": data.Memory.ValueInt64(),
		"min_memory": data.MinMemory.ValueString(),
		"hyperv_enlightenments": data.HypervEnlightenments.ValueBool(),
		"bootloader": data.Bootloader.ValueString(),
		"bootloader_ovmf": data.BootloaderOvmf.ValueString(),
		"autostart": data.Autostart.ValueBool(),
		"hide_from_msr": data.HideFromMsr.ValueBool(),
		"ensure_display_device": data.EnsureDisplayDevice.ValueBool(),
		"time": data.Time.ValueString(),
		"shutdown_timeout": data.ShutdownTimeout.ValueInt64(),
		"arch_type": data.ArchType.ValueString(),
		"machine_type": data.MachineType.ValueString(),
		"enable_secure_boot": data.EnableSecureBoot.ValueBool(),
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
	startOnCreate := true  // default
	if !data.StartOnCreate.IsNull() {
		startOnCreate = data.StartOnCreate.ValueBool()
	}
	if startOnCreate {
		_, err = r.client.Call("vm.start", data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
		}
	}
	// Set default for start_on_create if not specified
	if data.StartOnCreate.IsNull() {
		data.StartOnCreate = types.BoolValue(true)
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

	params := map[string]interface{}{
		"command_line_args": data.CommandLineArgs.ValueString(),
		"cpu_mode": data.CpuMode.ValueString(),
		"cpu_model": data.CpuModel.ValueString(),
		"name": data.Name.ValueString(),
		"description": data.Description.ValueString(),
		"vcpus": data.Vcpus.ValueInt64(),
		"cores": data.Cores.ValueInt64(),
		"threads": data.Threads.ValueInt64(),
		"cpuset": data.Cpuset.ValueString(),
		"nodeset": data.Nodeset.ValueString(),
		"enable_cpu_topology_extension": data.EnableCpuTopologyExtension.ValueBool(),
		"pin_vcpus": data.PinVcpus.ValueBool(),
		"suspend_on_snapshot": data.SuspendOnSnapshot.ValueBool(),
		"trusted_platform_module": data.TrustedPlatformModule.ValueBool(),
		"memory": data.Memory.ValueInt64(),
		"min_memory": data.MinMemory.ValueString(),
		"hyperv_enlightenments": data.HypervEnlightenments.ValueBool(),
		"bootloader": data.Bootloader.ValueString(),
		"bootloader_ovmf": data.BootloaderOvmf.ValueString(),
		"autostart": data.Autostart.ValueBool(),
		"hide_from_msr": data.HideFromMsr.ValueBool(),
		"ensure_display_device": data.EnsureDisplayDevice.ValueBool(),
		"time": data.Time.ValueString(),
		"shutdown_timeout": data.ShutdownTimeout.ValueInt64(),
		"arch_type": data.ArchType.ValueString(),
		"machine_type": data.MachineType.ValueString(),
		"enable_secure_boot": data.EnableSecureBoot.ValueBool(),
	}

	_, err := r.client.Call("vm.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
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
