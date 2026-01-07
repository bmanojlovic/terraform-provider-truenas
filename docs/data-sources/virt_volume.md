---
page_title: "truenas_virt_volume Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS virt_volume data
---

# truenas_virt_volume (Data Source)

Retrieves TrueNAS virt_volume data

## Example Usage

```terraform
data "truenas_virt_volume" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the virt_volume to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
