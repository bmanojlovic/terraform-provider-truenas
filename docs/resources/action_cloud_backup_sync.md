---
page_title: "truenas_action_cloud_backup_sync Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Run the cloud backup job `id`.
---

# truenas_action_cloud_backup_sync (Resource)

Run the cloud backup job `id`.

This is an action resource that executes the `cloud_backup.sync` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_cloud_backup_sync" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) The cloud backup task ID.
- `options` (String, Optional) Sync options.

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
