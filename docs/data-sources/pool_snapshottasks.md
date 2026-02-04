---
page_title: "truenas_pool_snapshottasks Data Source - terraform-provider-truenas"
subcategory: "Storage - Pools"
description: |-
  Query pool.snapshottask
---

# truenas_pool_snapshottasks (Data Source)

Query pool.snapshottask

## Example Usage

```terraform
data "truenas_pool_snapshottasks" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_snapshottasks to retrieve.

### Read-Only

- None
