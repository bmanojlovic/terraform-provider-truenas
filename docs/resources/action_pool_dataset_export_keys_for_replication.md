---
page_title: "truenas_action_pool_dataset_export_keys_for_replication Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Export keys for replication task `id` for source dataset(s) which are stored in the system. The exported file is a JSON file which has a dictionary containing dataset names as keys and their keys as the value.  Please refer to websocket documentation for downloading the file.
---

# truenas_action_pool_dataset_export_keys_for_replication (Resource)

Export keys for replication task `id` for source dataset(s) which are stored in the system. The exported file is a JSON file which has a dictionary containing dataset names as keys and their keys as the value.  Please refer to websocket documentation for downloading the file.

This is an action resource that executes the `pool.dataset.export_keys_for_replication` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_dataset_export_keys_for_replication" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) The pool ID to export dataset keys for replication purposes.

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
