---
page_title: "truenas_replication_restore_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes restore action on replication resource
---

# truenas_replication_restore_action (Resource)

Executes restore action on replication resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_replication" "example" {
  # ... resource configuration
}

resource "truenas_replication_restore_action" "example" {
  resource_id = truenas_replication.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the replication resource to perform action on

### Optional

- `replication_restore` (Optional) - replication_restore parameter. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the restore action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
