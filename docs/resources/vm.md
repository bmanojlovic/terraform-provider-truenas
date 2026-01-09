---
page_title: "truenas_vm Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  VMCreateArgs parameters.
---

# truenas_vm (Resource)

VMCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_vm" "example" {
  name = "example-name"
  memory = 1
  start_on_create = true
  command_line_args = ""
  cpu_mode = "CUSTOM"
  description = ""
  vcpus = 1
  cores = 1
}
```

## Schema

### Required

- `name` (Required) - Display name of the virtual machine.. Type: `string`
- `memory` (Required) - Amount of memory allocated to the VM in megabytes.. Type: `integer`

### Optional

- `start_on_create` (Optional) - Start the resource immediately after creation. Default behavior: starts if not specified. Type: `boolean`
- `command_line_args` (Optional) - Additional command line arguments passed to the VM hypervisor. Default: ``. Type: `string`
- `cpu_mode` (Optional) - CPU virtualization mode.

* `CUSTOM`: Use specified model.
* `HOST-MODEL`: Mirror host CPU.
* `HOST-PASSTHROUGH`: Provide direct access to host CPU features. Valid values: `CUSTOM`, `HOST-MODEL`, `HOST-PASSTHROUGH` Default: `CUSTOM`. Type: `string`
- `cpu_model` (Optional) - Specific CPU model to emulate. `null` to use hypervisor default.. Type: `string`
- `description` (Optional) - Optional description or notes about the virtual machine. Default: ``. Type: `string`
- `vcpus` (Optional) - Number of virtual CPUs allocated to the VM. Default: `1`. Type: `integer`
- `cores` (Optional) - Number of CPU cores per socket. Default: `1`. Type: `integer`
- `threads` (Optional) - Number of threads per CPU core. Default: `1`. Type: `integer`
- `cpuset` (Optional) - Set of host CPU cores to pin VM CPUs to. `null` for no pinning.. Type: `string`
- `nodeset` (Optional) - Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.. Type: `string`
- `enable_cpu_topology_extension` (Optional) - Whether to expose detailed CPU topology information to the guest OS. Default: `False`. Type: `boolean`
- `pin_vcpus` (Optional) - Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexibility. Default: `False`. Type: `boolean`
- `suspend_on_snapshot` (Optional) - Whether to suspend the VM when taking snapshots. Default: `False`. Type: `boolean`
- `trusted_platform_module` (Optional) - Whether to enable virtual Trusted Platform Module (TPM) for the VM. Default: `False`. Type: `boolean`
- `min_memory` (Optional) - Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink     during low usage but guarantees this minimum. `null` to disable ballooning.. Type: `string`
- `hyperv_enlightenments` (Optional) - Whether to enable Hyper-V enlightenments for improved Windows guest performance. Default: `False`. Type: `boolean`
- `bootloader` (Optional) - Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility. Valid values: `UEFI_CSM`, `UEFI` Default: `UEFI`. Type: `string`
- `bootloader_ovmf` (Optional) - OVMF firmware file to use for UEFI boot.. Type: `string`
- `autostart` (Optional) - Whether to automatically start the VM when the host system boots. Default: `True`. Type: `boolean`
- `hide_from_msr` (Optional) - Whether to hide hypervisor signatures from guest OS MSR access. Default: `False`. Type: `boolean`
- `ensure_display_device` (Optional) - Whether to ensure at least one display device is configured for the VM. Default: `True`. Type: `boolean`
- `time` (Optional) - Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time. Valid values: `LOCAL`, `UTC` Default: `LOCAL`. Type: `string`
- `shutdown_timeout` (Optional) - Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances     allowing sufficient time for clean shutdown while avoiding indefinite hangs. Default: `90`. Type: `integer`
- `arch_type` (Optional) - Guest architecture type. `null` to use hypervisor default.. Type: `string`
- `machine_type` (Optional) - Virtual machine type/chipset. `null` to use hypervisor default.. Type: `string`
- `enable_secure_boot` (Optional) - Whether to enable UEFI Secure Boot for enhanced security. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_vm.example <id>
```
