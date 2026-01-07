---
page_title: "truenas_app_upgrade_summary_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes upgrade_summary action on app resource
---

# truenas_app_upgrade_summary_action (Resource)

Executes upgrade_summary action on app resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_app" "example" {
  # ... resource configuration
}

resource "truenas_app_upgrade_summary_action" "example" {
  resource_id = truenas_app.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the app resource to perform action on

### Optional

- `app_name` (Optional) - app_name parameter. Type: `string`
- `options` (Optional) - options parameter. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the upgrade_summary action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
