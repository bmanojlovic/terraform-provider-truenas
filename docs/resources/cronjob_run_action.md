---
page_title: "truenas_cronjob_run_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes run action on cronjob resource
---

# truenas_cronjob_run_action (Resource)

Executes run action on cronjob resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_cronjob" "example" {
  # ... resource configuration
}

resource "truenas_cronjob_run_action" "example" {
  resource_id = truenas_cronjob.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the cronjob resource to perform action on

### Optional

- `skip_disabled` (Optional) - skip_disabled parameter. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the run action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
