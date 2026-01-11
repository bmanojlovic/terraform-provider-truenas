---
page_title: "truenas_vm Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a Virtual Machine (VM).
---

# truenas_vm (Resource)

Create a Virtual Machine (VM).


## Example Usage

```terraform
resource "truenas_vm" "example" {
  memory = 1
  name = "example-value"
  start_on_create = true
}
```

## Schema

### Required

- `memory` (Int64) - Amount of memory allocated to the VM in megabytes.
- `name` (String) - Display name of the virtual machine.

### Optional

- `start_on_create` (Bool) - Start the resource immediately after creation. Default: `true`
- `arch_type` (String) - Guest architecture type. `null` to use hypervisor default. Default: `None`
- `autostart` (Bool) - Whether to automatically start the VM when the host system boots. Default: `True`
- `bootloader` (String) - Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility. Default: `UEFI` Valid values: `UEFI_CSM`, `UEFI`
- `bootloader_ovmf` (String) - OVMF firmware file to use for UEFI boot. Default: `None`
- `command_line_args` (String) - Additional command line arguments passed to the VM hypervisor. Default: ``
- `cores` (Int64) - Number of CPU cores per socket. Default: `1`
- `cpu_mode` (String) - CPU virtualization mode.  * `CUSTOM`: Use specified model. * `HOST-MODEL`: Mirror host CPU. * `HOST-PASSTHROUGH`: Provide direct access to host CPU features. Default: `CUSTOM` Valid values: `CUSTOM`, `HOST-MODEL`, `HOST-PASSTHROUGH`
- `cpu_model` (String) - Specific CPU model to emulate. `null` to use hypervisor default. Default: `None`
- `cpuset` (String) - Set of host CPU cores to pin VM CPUs to. `null` for no pinning. Default: `None`
- `description` (String) - Optional description or notes about the virtual machine. Default: ``
- `enable_cpu_topology_extension` (Bool) - Whether to expose detailed CPU topology information to the guest OS. Default: `False`
- `enable_secure_boot` (Bool) - Whether to enable UEFI Secure Boot for enhanced security. Default: `False`
- `ensure_display_device` (Bool) - Whether to ensure at least one display device is configured for the VM. Default: `True`
- `hide_from_msr` (Bool) - Whether to hide hypervisor signatures from guest OS MSR access. Default: `False`
- `hyperv_enlightenments` (Bool) - Whether to enable Hyper-V enlightenments for improved Windows guest performance. Default: `False`
- `machine_type` (String) - Virtual machine type/chipset. `null` to use hypervisor default. Default: `None`
- `min_memory` (Int64) - Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink     during low usage but guarantees this minimum. `null` to disable ballooning. Default: `None`
- `nodeset` (String) - Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints. Default: `None`
- `pin_vcpus` (Bool) - Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexibility. Default: `False`
- `shutdown_timeout` (Int64) - Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances     allowing sufficient time for clean shutdown while avoiding indefinite hangs. Default: `90`
- `suspend_on_snapshot` (Bool) - Whether to suspend the VM when taking snapshots. Default: `False`
- `threads` (Int64) - Number of threads per CPU core. Default: `1`
- `time` (String) - Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time. Default: `LOCAL` Valid values: `LOCAL`, `UTC`
- `trusted_platform_module` (Bool) - Whether to enable virtual Trusted Platform Module (TPM) for the VM. Default: `False`
- `vcpus` (Int64) - Number of virtual CPUs allocated to the VM. Default: `1`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_vm.example <id>
```
