---
page_title: "truenas_action_pool_snapshot_rollback Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Execute pool.snapshot.rollback
---

# truenas_action_pool_snapshot_rollback (Resource)

Execute pool.snapshot.rollback

This is an action resource that executes the `pool.snapshot.rollback` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_snapshot_rollback" "example" {
  id = "value"
}
```

## Schema

### Input Parameters

- `id` (String, Required) ID of the snapshot to rollback to.
- `options` (String, Optional) Options for controlling snapshot rollback behavior.

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
