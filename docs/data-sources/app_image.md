---
page_title: "truenas_app_image Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS app_image data
---

# truenas_app_image (Data Source)

Retrieves TrueNAS app_image data

## Example Usage

```terraform
data "truenas_app_image" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the app_image to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
