---
page_title: "truenas_action_pool_remove Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Remove a disk from pool of id `id`.  `label` is the vdev guid or device name.  Error codes:      EZFS_NOSPC(2032): out of space to remove a device     EZFS_NODEVICE(2017): no such device in pool     EZFS_NOREPLICAS(2019): no valid replicas
---

# truenas_action_pool_remove (Resource)

Remove a disk from pool of id `id`.  `label` is the vdev guid or device name.  Error codes:      EZFS_NOSPC(2032): out of space to remove a device     EZFS_NODEVICE(2017): no such device in pool     EZFS_NOREPLICAS(2019): no valid replicas

This is an action resource that executes the `pool.remove` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_remove" "example" {
  id = 1
  options = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the pool to remove a disk from.
- `options` (String, Required) Disk identifier to remove from the pool.

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
