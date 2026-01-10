---
page_title: "truenas_vm_device Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new device for the VM of id `vm`.
---

# truenas_vm_device (Resource)

Create a new device for the VM of id `vm`.

## Example Usage

```terraform
resource "truenas_vm_device" "example" {
  attributes = "example-value"
  vm = 1
}
```

## Schema

### Required

- `attributes` (String) - Device-specific configuration attributes.
- `vm` (Int64) - ID of the virtual machine this device belongs to.

### Optional

- `order` (Int64) - Boot order priority for this device. `null` for automatic assignment. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_vm_device.example <id>
```
