---
page_title: "truenas_action_pool_dataset_export_key Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Export own encryption key for dataset `id`. If `download` is `true`, key will be downloaded in a json file where the same file can be used to unlock the dataset, otherwise it will be returned as string.  Please refer to websocket documentation for downloading the file.
---

# truenas_action_pool_dataset_export_key (Resource)

Export own encryption key for dataset `id`. If `download` is `true`, key will be downloaded in a json file where the same file can be used to unlock the dataset, otherwise it will be returned as string.  Please refer to websocket documentation for downloading the file.

This is an action resource that executes the `pool.dataset.export_key` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_dataset_export_key" "example" {
  id = "value"
}
```

## Schema

### Input Parameters

- `id` (String, Required) The dataset ID (full path) to export the encryption key from.
- `download` (Bool, Optional) Whether to prepare the key for download as a file.

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
