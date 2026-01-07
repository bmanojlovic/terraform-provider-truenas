---
page_title: "truenas_app Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS app data
---

# truenas_app (Data Source)

Retrieves TrueNAS app data

## Example Usage

```terraform
data "truenas_app" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the app to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
