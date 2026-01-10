---
page_title: "truenas_users Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query users with `query-filters` and `query-options`.
---

# truenas_users (Data Source)

Query users with `query-filters` and `query-options`.

## Example Usage

```terraform
# Get all user
data "truenas_users" "all" {}

# Access items
output "user_count" {
  value = length(data.truenas_users.all.items)
}

output "user_names" {
  value = [for item in data.truenas_users.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of user resources

Items have the following attributes:

- None
