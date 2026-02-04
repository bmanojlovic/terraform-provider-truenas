---
page_title: "truenas_interfaces Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query Interfaces with `query-filters` and `query-options`
---

# truenas_interfaces (Data Source)

Query Interfaces with `query-filters` and `query-options`

## Example Usage

```terraform
data "truenas_interfaces" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the interfaces to retrieve.

### Read-Only

- None
