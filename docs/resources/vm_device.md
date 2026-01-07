---
page_title: "truenas_vm_device Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  VMDeviceCreateArgs parameters.
---

# truenas_vm_device (Resource)

VMDeviceCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_vm_device" "example" {
  attributes = "example-attributes"
  vm = 1
}
```

## Schema

### Required

- `attributes` (Required) - Device-specific configuration attributes.. Type: `string`
- `vm` (Required) - ID of the virtual machine this device belongs to.. Type: `integer`

### Optional

- `order` (Optional) - Boot order priority for this device. `null` for automatic assignment.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_vm_device.example <id>
```
