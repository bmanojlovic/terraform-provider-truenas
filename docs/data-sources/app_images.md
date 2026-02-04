---
page_title: "truenas_app_images Data Source - terraform-provider-truenas"
subcategory: "Applications"
description: |-
  Query all docker images with `query-filters` and `query-options`.
---

# truenas_app_images (Data Source)

Query all docker images with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_app_images" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the app_images to retrieve.

### Read-Only

- None
