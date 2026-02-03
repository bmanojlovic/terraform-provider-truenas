---
page_title: "truenas_action_cloudsync_sync Resource - terraform-provider-truenas"
subcategory: "Actions - Cloud Sync"
description: |-
  Run the cloud_sync job `id`, syncing the local data to remote.
---

# truenas_action_cloudsync_sync (Resource)

Run the cloud_sync job `id`, syncing the local data to remote.

This is an action resource that executes the `cloudsync.sync` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_cloudsync_sync" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the cloud sync task to run.
- `cloud_sync_sync_options` (String, Optional) Options for the sync operation.

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
