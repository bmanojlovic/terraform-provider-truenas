---
page_title: "truenas_virt_volumes Data Source - terraform-provider-truenas"
subcategory: "Virtualization"
description: |-
  Query virt.volume
---

# truenas_virt_volumes (Data Source)

Query virt.volume

## Example Usage

```terraform
data "truenas_virt_volumes" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the virt_volumes to retrieve.

### Read-Only

- None
