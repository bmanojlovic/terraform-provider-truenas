---
page_title: "truenas_action_pool_dataset_lock Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Locks `id` dataset. It will unmount the dataset and its children before locking.  After the dataset has been unmounted, system will set immutable flag on the dataset's mountpoint where the dataset was mounted before it was locked making sure that the path cannot be modified. Once the dataset is unlocked, it will not be affected by this change and consumers can continue consuming it.
---

# truenas_action_pool_dataset_lock (Resource)

Locks `id` dataset. It will unmount the dataset and its children before locking.  After the dataset has been unmounted, system will set immutable flag on the dataset's mountpoint where the dataset was mounted before it was locked making sure that the path cannot be modified. Once the dataset is unlocked, it will not be affected by this change and consumers can continue consuming it.

This is an action resource that executes the `pool.dataset.lock` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_dataset_lock" "example" {
  id = "value"
}
```

## Schema

### Input Parameters

- `id` (String, Required) The dataset ID (full path) to lock.
- `options` (String, Optional) Options for locking the dataset, such as force unmount settings.

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
