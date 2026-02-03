---
page_title: "truenas_action_pool_attach Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  `target_vdev` is the GUID of the vdev where the disk needs to be attached. In case of STRIPED vdev, this is the STRIPED disk GUID which will be converted to mirror. If `target_vdev` is mirror, it will be converted into a n-way mirror.
---

# truenas_action_pool_attach (Resource)

`target_vdev` is the GUID of the vdev where the disk needs to be attached. In case of STRIPED vdev, this is the STRIPED disk GUID which will be converted to mirror. If `target_vdev` is mirror, it will be converted into a n-way mirror.

This is an action resource that executes the `pool.attach` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_attach" "example" {
  oid = 1
  options = "value"
}
```

## Schema

### Input Parameters

- `oid` (Int64, Required) ID of the pool to attach a disk to.
- `options` (String, Required) Configuration for the disk attachment operation.

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
