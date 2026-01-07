---
page_title: "truenas_cloud_backup_restore_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes restore action on cloud_backup resource
---

# truenas_cloud_backup_restore_action (Resource)

Executes restore action on cloud_backup resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_cloud_backup" "example" {
  # ... resource configuration
}

resource "truenas_cloud_backup_restore_action" "example" {
  resource_id = truenas_cloud_backup.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the cloud_backup resource to perform action on

### Optional

- `snapshot_id` (Optional) - snapshot_id parameter. Type: `string`
- `subfolder` (Optional) - subfolder parameter. Type: `string`
- `destination_path` (Optional) - destination_path parameter. Type: `string`
- `options` (Optional) - options parameter. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the restore action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
