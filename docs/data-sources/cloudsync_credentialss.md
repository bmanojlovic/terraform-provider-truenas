---
page_title: "truenas_cloudsync_credentialss Data Source - terraform-provider-truenas"
subcategory: "Cloud Sync"
description: |-
  Query cloudsync.credentials
---

# truenas_cloudsync_credentialss (Data Source)

Query cloudsync.credentials

## Example Usage

```terraform
data "truenas_cloudsync_credentialss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cloudsync_credentialss to retrieve.

### Read-Only

- None
