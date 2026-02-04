---
page_title: "truenas_pool_snapshots Data Source - terraform-provider-truenas"
subcategory: "Storage - Pools"
description: |-
  Query all ZFS Snapshots with `query-filters` and `query-options`.
---

# truenas_pool_snapshots (Data Source)

Query all ZFS Snapshots with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_pool_snapshots" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_snapshots to retrieve.

### Read-Only

- None
