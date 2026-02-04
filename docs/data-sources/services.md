---
page_title: "truenas_services Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query all system services with `query-filters` and `query-options`.
---

# truenas_services (Data Source)

Query all system services with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_services" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the services to retrieve.

### Read-Only

- None
