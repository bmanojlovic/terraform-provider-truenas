---
page_title: "truenas_cloudsync Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS cloudsync data
---

# truenas_cloudsync (Data Source)

Retrieves TrueNAS cloudsync data

## Example Usage

```terraform
data "truenas_cloudsync" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cloudsync to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
