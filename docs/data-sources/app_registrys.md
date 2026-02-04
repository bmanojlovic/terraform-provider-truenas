---
page_title: "truenas_app_registrys Data Source - terraform-provider-truenas"
subcategory: "Applications"
description: |-
  Query app.registry
---

# truenas_app_registrys (Data Source)

Query app.registry

## Example Usage

```terraform
data "truenas_app_registrys" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the app_registrys to retrieve.

### Read-Only

- None
