---
page_title: "truenas_pools Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query pool resources
---

# truenas_pools (Data Source)

Query pool resources

## Example Usage

```terraform
# Get all pool
data "truenas_pools" "all" {}

# Access items
output "pool_count" {
  value = length(data.truenas_pools.all.items)
}

output "pool_names" {
  value = [for item in data.truenas_pools.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of pool resources

Items have the following attributes:

- None
