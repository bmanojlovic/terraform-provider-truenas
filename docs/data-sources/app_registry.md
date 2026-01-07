---
page_title: "truenas_app_registry Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS app_registry data
---

# truenas_app_registry (Data Source)

Retrieves TrueNAS app_registry data

## Example Usage

```terraform
data "truenas_app_registry" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the app_registry to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
