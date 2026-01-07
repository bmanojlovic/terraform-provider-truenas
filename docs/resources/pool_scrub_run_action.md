---
page_title: "truenas_pool_scrub_run_action Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Executes run action on pool_scrub resource
---

# truenas_pool_scrub_run_action (Resource)

Executes run action on pool_scrub resource

This is an action resource that executes an operation when created or updated. It cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_pool_scrub" "example" {
  # ... resource configuration
}

resource "truenas_pool_scrub_run_action" "example" {
  resource_id = truenas_pool_scrub.example.id
}
```

## Schema

### Required

- `resource_id` (String) ID of the pool_scrub resource to perform action on

### Optional

- `name` (Optional) - name parameter. Type: `string`
- `threshold` (Optional) - threshold parameter. Type: `string`

### Read-Only

- `id` (String) Action execution ID (timestamp-based)

## Notes

- This resource executes the run action when created
- Updates will re-execute the action
- Deletion removes from state but cannot undo the action
- Use with caution as actions are immediate and irreversible
