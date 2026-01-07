---
page_title: "truenas_cloud_backup Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS cloud_backup data
---

# truenas_cloud_backup (Data Source)

Retrieves TrueNAS cloud_backup data

## Example Usage

```terraform
data "truenas_cloud_backup" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cloud_backup to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
