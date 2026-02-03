---
page_title: "truenas_action_docker_backup Resource - terraform-provider-truenas"
subcategory: "Actions - Applications"
description: |-
  Create a backup of existing apps.  This creates a backup of existing apps on the same pool in which docker is initialized.
---

# truenas_action_docker_backup (Resource)

Create a backup of existing apps.  This creates a backup of existing apps on the same pool in which docker is initialized.

This is an action resource that executes the `docker.backup` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_docker_backup" "example" {
}
```

## Schema

### Input Parameters

- `backup_name` (String, Optional) Name for the backup or `null` to generate a timestamp-based name.

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
