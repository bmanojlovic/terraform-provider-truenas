---
page_title: "truenas_vm_device Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS vm_device data
---

# truenas_vm_device (Data Source)

Retrieves TrueNAS vm_device data

## Example Usage

```terraform
data "truenas_vm_device" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vm_device to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
