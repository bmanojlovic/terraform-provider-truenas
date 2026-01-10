---
page_title: "truenas_disks Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query disks.
---

# truenas_disks (Data Source)

Query disks.

## Example Usage

```terraform
# Get all disk
data "truenas_disks" "all" {}

# Access items
output "disk_count" {
  value = length(data.truenas_disks.all.items)
}

output "disk_names" {
  value = [for item in data.truenas_disks.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of disk resources

Items have the following attributes:

- None
