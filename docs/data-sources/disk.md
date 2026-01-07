---
page_title: "truenas_disk Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS disk data
---

# truenas_disk (Data Source)

Retrieves TrueNAS disk data

## Example Usage

```terraform
data "truenas_disk" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the disk to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
