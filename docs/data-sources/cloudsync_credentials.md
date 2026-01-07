---
page_title: "truenas_cloudsync_credentials Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS cloudsync_credentials data
---

# truenas_cloudsync_credentials (Data Source)

Retrieves TrueNAS cloudsync_credentials data

## Example Usage

```terraform
data "truenas_cloudsync_credentials" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cloudsync_credentials to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
