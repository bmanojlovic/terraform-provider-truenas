---
page_title: "truenas_action_pool_replace Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Replace a disk on a pool.  `label` is the ZFS guid or a device name `disk` is the identifier of a disk If `preserve_settings` is true, then settings (power management, S.M.A.R.T., etc.) of a disk being replaced will be applied to a new disk.
---

# truenas_action_pool_replace (Resource)

Replace a disk on a pool.  `label` is the ZFS guid or a device name `disk` is the identifier of a disk If `preserve_settings` is true, then settings (power management, S.M.A.R.T., etc.) of a disk being replaced will be applied to a new disk.

This is an action resource that executes the `pool.replace` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_replace" "example" {
  id = 1
  options = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the pool to replace a disk in.
- `options` (String, Required) Configuration for the disk replacement operation.

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
