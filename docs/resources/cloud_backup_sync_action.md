---
page_title: "truenas_cloud_backup_sync_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Sync options.
---

# truenas_cloud_backup_sync_action (Resource)

Sync options.

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_cloud_backup" "example" {
  # ... resource configuration
}

resource "truenas_cloud_backup_sync_action" "example" {
  resource_id = truenas_cloud_backup.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the cloud_backup resource to perform action on

### Optional

- `dry_run` (Optional) - Simulate the backup without actually writing to the remote repository.. Type: `boolean`
- `rate_limit` (Optional) - Maximum upload rate in KiB/s. Passed to `restic --limit-upload`.

If provided, overrides the task's rate limit.. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the sync action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
