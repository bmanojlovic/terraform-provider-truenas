---
page_title: "truenas_services Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Query all system services with `query-filters` and `query-options`.
---

# truenas_services (Data Source)

Query all system services with `query-filters` and `query-options`.

## Example Usage

```terraform
# Get all service
data "truenas_services" "all" {}

# Access items
output "service_count" {
  value = length(data.truenas_services.all.items)
}

output "service_names" {
  value = [for item in data.truenas_services.all.items : item.name]
}
```

## Schema

### Read-Only

- `items` (List of Object) - List of service resources

Items have the following attributes:

- None
