---
page_title: "truenas_interfaces Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query Interfaces with `query-filters` and `query-options`
---

# truenas_interfaces (Data Source)

Query Interfaces with `query-filters` and `query-options`

## Example Usage

```terraform
# Get all interface
data "truenas_interfaces" "all" {}

# Access items
output "interface_count" {
  value = length(data.truenas_interfaces.all.items)
}

output "interface_names" {
  value = [for item in data.truenas_interfaces.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of interface resources

Items have the following attributes:

- None
