---
page_title: "truenas_vm_devices Data Source - terraform-provider-truenas"
subcategory: "Virtual Machines"
description: |-
  Query vm.device
---

# truenas_vm_devices (Data Source)

Query vm.device

## Example Usage

```terraform
data "truenas_vm_devices" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vm_devices to retrieve.

### Read-Only

- None
