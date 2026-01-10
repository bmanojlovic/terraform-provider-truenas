---
page_title: "truenas_groups Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query groups with `query-filters` and `query-options`.
---

# truenas_groups (Data Source)

Query groups with `query-filters` and `query-options`.

## Example Usage

```terraform
# Get all group
data "truenas_groups" "all" {}

# Access items
output "group_count" {
  value = length(data.truenas_groups.all.items)
}

output "group_names" {
  value = [for item in data.truenas_groups.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of group resources

Items have the following attributes:

- None
