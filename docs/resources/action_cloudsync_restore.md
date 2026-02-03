---
page_title: "truenas_action_cloudsync_restore Resource - terraform-provider-truenas"
subcategory: "Actions - Cloud Sync"
description: |-
  Create the opposite of cloud sync task `id` (PULL if it was PUSH and vice versa).
---

# truenas_action_cloudsync_restore (Resource)

Create the opposite of cloud sync task `id` (PULL if it was PUSH and vice versa).

This is an action resource that executes the `cloudsync.restore` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_cloudsync_restore" "example" {
  id = 1
  opts = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the cloud sync task to restore from.
- `opts` (String, Required) Restore operation configuration options.

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
