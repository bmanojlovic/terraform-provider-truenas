---
page_title: "truenas_virt_instance Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS virt_instance data
---

# truenas_virt_instance (Data Source)

Retrieves TrueNAS virt_instance data

## Example Usage

```terraform
data "truenas_virt_instance" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the virt_instance to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
