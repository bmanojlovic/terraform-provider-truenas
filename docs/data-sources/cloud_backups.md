---
page_title: "truenas_cloud_backups Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query cloud_backup
---

# truenas_cloud_backups (Data Source)

Query cloud_backup

## Example Usage

```terraform
data "truenas_cloud_backups" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cloud_backups to retrieve.

### Read-Only

- None
