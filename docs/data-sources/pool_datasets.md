---
page_title: "truenas_pool_datasets Data Source - terraform-provider-truenas"
subcategory: "Storage - Pools"
description: |-
  Query Pool Datasets with `query-filters` and `query-options`.
---

# truenas_pool_datasets (Data Source)

Query Pool Datasets with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_pool_datasets" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_datasets to retrieve.

### Read-Only

- None
