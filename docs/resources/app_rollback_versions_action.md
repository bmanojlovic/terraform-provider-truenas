---
page_title: "truenas_app_rollback_versions_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Name of the application to get rollback versions for.
---

# truenas_app_rollback_versions_action (Resource)

Name of the application to get rollback versions for.

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_app" "example" {
  # ... resource configuration
}

resource "truenas_app_rollback_versions_action" "example" {
  resource_id = truenas_app.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the app resource to perform action on

### Optional

- None

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the rollback_versions action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
