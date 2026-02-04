---
page_title: "truenas_virt_instances Data Source - terraform-provider-truenas"
subcategory: "Virtualization"
description: |-
  Query all instances with `query-filters` and `query-options`.
---

# truenas_virt_instances (Data Source)

Query all instances with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_virt_instances" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the virt_instances to retrieve.

### Read-Only

- None
