---
page_title: "truenas_tunable Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new system tunable.
---

# truenas_tunable (Resource)

Configuration for creating a new system tunable.

## Example Usage

```terraform
resource "truenas_tunable" "example" {
  var = "example-var"
  value = "example-value"
  type = "SYSCTL"
  comment = ""
  enabled = true
  update_initramfs = true
}
```

## Schema

### Required

- `var` (Required) - Name or identifier of the system parameter to tune.. Type: `string`
- `value` (Required) - Value to assign to the tunable parameter.. Type: `string`

### Optional

- `type` (Optional) - * `SYSCTL`: `var` is a sysctl name (e.g. `kernel.watchdog`) and `value` is its corresponding value (e.g. `0`).
* `UDEV`: `var` is a udev rules file name (e.g. `10-disable-usb`, `.rules` suffix will be appended automatically)     and `value` is its contents (e.g. `BUS=="usb", OPTIONS+="ignore_device"`).
* `ZFS`: `var` is a ZFS kernel module parameter name (e.g. `zfs_dirty_data_max_max`) and `value` is its value     (e.g. `783091712`). Valid values: `SYSCTL`, `UDEV`, `ZFS` Default: `SYSCTL`. Type: `string`
- `comment` (Optional) - Optional descriptive comment explaining the purpose of this tunable. Default: ``. Type: `string`
- `enabled` (Optional) - Whether this tunable is active and should be applied. Default: `True`. Type: `boolean`
- `update_initramfs` (Optional) - If `false`, then initramfs will not be updated after creating a ZFS tunable and you will need to run     `system boot update_initramfs` manually. Default: `True`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_tunable.example <id>
```
