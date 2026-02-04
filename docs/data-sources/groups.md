---
page_title: "truenas_groups Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query groups with `query-filters` and `query-options`.
---

# truenas_groups (Data Source)

Query groups with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_groups" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the groups to retrieve.

### Read-Only

- None
