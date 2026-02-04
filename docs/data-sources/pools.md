---
page_title: "truenas_pools Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query pool
---

# truenas_pools (Data Source)

Query pool

## Example Usage

```terraform
data "truenas_pools" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pools to retrieve.

### Read-Only

- None
