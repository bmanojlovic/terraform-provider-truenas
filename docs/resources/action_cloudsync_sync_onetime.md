---
page_title: "truenas_action_cloudsync_sync_onetime Resource - terraform-provider-truenas"
subcategory: "Actions - Cloud Sync"
description: |-
  Run cloud sync task without creating it.
---

# truenas_action_cloudsync_sync_onetime (Resource)

Run cloud sync task without creating it.

This is an action resource that executes the `cloudsync.sync_onetime` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_cloudsync_sync_onetime" "example" {
  cloud_sync_sync_onetime = "value"
}
```

## Schema

### Input Parameters

- `cloud_sync_sync_onetime` (String, Required) Cloud sync task configuration for one-time execution.
- `cloud_sync_sync_onetime_options` (String, Optional) Options for the one-time sync operation.

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
