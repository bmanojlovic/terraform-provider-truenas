---
page_title: "truenas_action_pool_dataset_destroy_snapshots Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Destroy specified snapshots of a given dataset.
---

# truenas_action_pool_dataset_destroy_snapshots (Resource)

Destroy specified snapshots of a given dataset.

This is an action resource that executes the `pool.dataset.destroy_snapshots` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_dataset_destroy_snapshots" "example" {
  name = "value"
  snapshots = "value"
}
```

## Schema

### Input Parameters

- `name` (String, Required) The dataset name to destroy snapshots for.
- `snapshots` (String, Required) Specification of which snapshots to destroy (all, specific ones, or ranges).

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
