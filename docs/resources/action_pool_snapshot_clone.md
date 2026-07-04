---
page_title: "truenas_action_pool_snapshot_clone Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Clone a snapshot to a new dataset/zvol.
---

# truenas_action_pool_snapshot_clone (Resource)

Clone a snapshot to a new dataset/zvol.

This is an action resource that executes the `pool.snapshot.clone` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_snapshot_clone" "clone" {
  snapshot    = "tank/source@snap1"
  dataset_dst = "tank/destination"
}
```

## Schema

### Input Parameters

- `data` (String, Required) PoolSnapshotCloneArgs parameters.

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
