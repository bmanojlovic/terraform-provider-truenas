---
page_title: "truenas_action_cloud_backup_restore Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Restore files to the directory `destination_path` from the `snapshot_id` subfolder `subfolder` created by the cloud backup job `id`.
---

# truenas_action_cloud_backup_restore (Resource)

Restore files to the directory `destination_path` from the `snapshot_id` subfolder `subfolder` created by the cloud backup job `id`.

This is an action resource that executes the `cloud_backup.restore` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_cloud_backup_restore" "example" {
  id = 1
  snapshot_id = "value"
  subfolder = "value"
  destination_path = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the cloud backup task.
- `snapshot_id` (String, Required) ID of the snapshot to restore.
- `subfolder` (String, Required) Path within the snapshot to restore.
- `destination_path` (String, Required) Local path to restore to.
- `options` (String, Optional) Additional restore options.

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
