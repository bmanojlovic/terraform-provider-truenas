---
page_title: "truenas_pool_datasets Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query Pool Datasets with `query-filters` and `query-options`.
---

# truenas_pool_datasets (Data Source)

Query Pool Datasets with `query-filters` and `query-options`.

## Example Usage

```terraform
# Get all pool_dataset
data "truenas_pool_datasets" "all" {}

# Access items
output "pool_dataset_count" {
  value = length(data.truenas_pool_datasets.all.items)
}

output "pool_dataset_names" {
  value = [for item in data.truenas_pool_datasets.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of pool_dataset resources

Items have the following attributes:

- None
