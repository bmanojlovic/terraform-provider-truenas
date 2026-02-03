---
page_title: "truenas_action_docker_delete_backup Resource - terraform-provider-truenas"
subcategory: "Actions - Applications"
description: |-
  Delete `backup_name` app backup.
---

# truenas_action_docker_delete_backup (Resource)

Delete `backup_name` app backup.

This is an action resource that executes the `docker.delete_backup` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_docker_delete_backup" "example" {
  backup_name = "value"
}
```

## Schema

### Input Parameters

- `backup_name` (String, Required) Name of the backup to delete.

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
