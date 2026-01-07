---
page_title: "truenas_iscsi_extent Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS iscsi_extent data
---

# truenas_iscsi_extent (Data Source)

Retrieves TrueNAS iscsi_extent data

## Example Usage

```terraform
data "truenas_iscsi_extent" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_extent to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
