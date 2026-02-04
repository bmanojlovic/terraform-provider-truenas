---
page_title: "truenas_app_ix_volumes Data Source - terraform-provider-truenas"
subcategory: "Applications"
description: |-
  Query ix-volumes with `filters` and `options`.
---

# truenas_app_ix_volumes (Data Source)

Query ix-volumes with `filters` and `options`.

## Example Usage

```terraform
data "truenas_app_ix_volumes" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the app_ix_volumes to retrieve.

### Read-Only

- None
