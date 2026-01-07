---
page_title: "truenas_pool_snapshot_rollback_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes rollback action on pool_snapshot resource
---

# truenas_pool_snapshot_rollback_action (Resource)

Executes rollback action on pool_snapshot resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_pool_snapshot" "example" {
  # ... resource configuration
}

resource "truenas_pool_snapshot_rollback_action" "example" {
  resource_id = truenas_pool_snapshot.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the pool_snapshot resource to perform action on

### Optional

- `options` (Optional) - options parameter. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the rollback action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
