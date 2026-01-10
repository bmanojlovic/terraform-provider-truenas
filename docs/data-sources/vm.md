---
page_title: "truenas_vm Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_vm (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_vm" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vm to retrieve.

### Read-Only

- `arch_type` (String) - Guest architecture type. `null` to use hypervisor default.
- `autostart` (Bool) - Whether to automatically start the VM when the host system boots.
- `bootloader` (String) - Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility.
- `bootloader_ovmf` (String) - OVMF firmware file to use for UEFI boot.
- `command_line_args` (String) - Additional command line arguments passed to the VM hypervisor.
- `cores` (Int64) - Number of CPU cores per socket.
- `cpu_mode` (String) - CPU virtualization mode.  * `CUSTOM`: Use specified model. * `HOST-MODEL`: Mirror host CPU. * `HOST-PASSTHROUGH`: Provide direct access to host CPU features.
- `cpu_model` (String) - Specific CPU model to emulate. `null` to use hypervisor default.
- `cpuset` (String) - Set of host CPU cores to pin VM CPUs to. `null` for no pinning.
- `description` (String) - Optional description or notes about the virtual machine.
- `devices` (List) - Array of virtual devices attached to this VM.
- `display_available` (Bool) - Whether at least one display device is available for this VM.
- `enable_cpu_topology_extension` (Bool) - Whether to expose detailed CPU topology information to the guest OS.
- `enable_secure_boot` (Bool) - Whether to enable UEFI Secure Boot for enhanced security.
- `ensure_display_device` (Bool) - Whether to ensure at least one display device is configured for the VM.
- `hide_from_msr` (Bool) - Whether to hide hypervisor signatures from guest OS MSR access.
- `hyperv_enlightenments` (Bool) - Whether to enable Hyper-V enlightenments for improved Windows guest performance.
- `machine_type` (String) - Virtual machine type/chipset. `null` to use hypervisor default.
- `memory` (Int64) - Amount of memory allocated to the VM in megabytes.
- `min_memory` (Int64) - Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink     during low usage but guarantees this minimum. `null` to disable ballooning.
- `name` (String) - Display name of the virtual machine.
- `nodeset` (String) - Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.
- `pin_vcpus` (Bool) - Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexibility.
- `shutdown_timeout` (Int64) - Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances     allowing sufficient time for clean shutdown while avoiding indefinite hangs.
- `status` (String) - Current runtime status information for the VM.
- `suspend_on_snapshot` (Bool) - Whether to suspend the VM when taking snapshots.
- `threads` (Int64) - Number of threads per CPU core.
- `time` (String) - Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time.
- `trusted_platform_module` (Bool) - Whether to enable virtual Trusted Platform Module (TPM) for the VM.
- `uuid` (String) - Unique UUID for the VM. `null` to auto-generate.
- `vcpus` (Int64) - Number of virtual CPUs allocated to the VM.
