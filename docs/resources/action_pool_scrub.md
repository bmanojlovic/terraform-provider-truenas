---
page_title: "truenas_action_pool_scrub Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Performs a scrub action to pool of `id`.  `action` can be either of "START", "STOP" or "PAUSE".
---

# truenas_action_pool_scrub (Resource)

Performs a scrub action to pool of `id`.  `action` can be either of "START", "STOP" or "PAUSE".

This is an action resource that executes the `pool.scrub` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_scrub" "example" {
  id = 1
  action = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the pool to perform scrub action on.
- `action` (String, Required) The scrub action to perform.

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
