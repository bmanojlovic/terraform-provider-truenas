---
page_title: "truenas_pool_snapshottask_run_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes run action on pool_snapshottask resource
---

# truenas_pool_snapshottask_run_action (Resource)

Executes run action on pool_snapshottask resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_pool_snapshottask" "example" {
  # ... resource configuration
}

resource "truenas_pool_snapshottask_run_action" "example" {
  resource_id = truenas_pool_snapshottask.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the pool_snapshottask resource to perform action on

### Optional

- None

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the run action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
