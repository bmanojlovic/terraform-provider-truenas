---
page_title: "truenas_vms Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query vm resources
---

# truenas_vms (Data Source)

Query vm resources

## Example Usage

```terraform
# Get all vm
data "truenas_vms" "all" {}

# Access items
output "vm_count" {
  value = length(data.truenas_vms.all.items)
}

output "vm_names" {
  value = [for item in data.truenas_vms.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of vm resources

Items have the following attributes:

- None
