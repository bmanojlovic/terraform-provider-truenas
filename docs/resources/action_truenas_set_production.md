---
page_title: "truenas_action_truenas_set_production Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Sets system production state and optionally sends initial debug.
---

# truenas_action_truenas_set_production (Resource)

Sets system production state and optionally sends initial debug.

This is an action resource that executes the `truenas.set_production` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_truenas_set_production" "example" {
  production = true
}
```

## Schema

### Input Parameters

- `production` (Bool, Required) Whether to configure the system for production use.
- `attach_debug` (Bool, Optional) Whether to attach debug information when transitioning to production mode.

### Computed Outputs

- `action_id` (String) Unique identifier for this action execution
- `job_id` (Int64) Background job ID (if applicable)
- `state` (String) Job state: SUCCESS, FAILED, or RUNNING
- `progress` (Float64) Job progress percentage (0-100)
- `result` (String) Action result data
- `error` (String) Error message if action failed

## Notes

- Actions execute immediately when the resource is created
- Background jobs are monitored until completion
- Progress updates are logged during execution
- The resource cannot be updated - changes force recreation
- Destroying the resource does not undo the action
