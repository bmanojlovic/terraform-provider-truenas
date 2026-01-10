---
page_title: "truenas_tunable Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a tunable.
---

# truenas_tunable (Resource)

Create a tunable.

## Example Usage

```terraform
resource "truenas_tunable" "example" {
  value = "example-value"
  var = "example-value"
}
```

## Schema

### Required

- `value` (String) - Value to assign to the tunable parameter.
- `var` (String) - Name or identifier of the system parameter to tune.

### Optional

- `comment` (String) - Optional descriptive comment explaining the purpose of this tunable. Default: ``
- `enabled` (Bool) - Whether this tunable is active and should be applied. Default: `True`
- `type` (String) - * `SYSCTL`: `var` is a sysctl name (e.g. `kernel.watchdog`) and `value` is its corresponding value (e.g. `0`). * `UDEV`: `var` is a udev rules file name (e.g. `10-disable-usb`, `.rules` suffix will be Default: `SYSCTL` Valid values: `SYSCTL`, `UDEV`, `ZFS`
- `update_initramfs` (Bool) - If `false`, then initramfs will not be updated after creating a ZFS tunable and you will need to run     `system boot update_initramfs` manually. Default: `True`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_tunable.example <id>
```
