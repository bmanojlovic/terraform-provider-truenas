---
page_title: "truenas_action_pool_online Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Bring a pool disk online.
---

# truenas_action_pool_online (Resource)

Bring a pool disk online.

This is an action resource that executes the `pool.online` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_online" "example" {
  id = 1
  options = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the pool to bring a disk online in.
- `options` (String, Required) Disk identifier to bring online.

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
