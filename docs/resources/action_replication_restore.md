---
page_title: "truenas_action_replication_restore Resource - terraform-provider-truenas"
subcategory: "Actions - Replication"
description: |-
  Create the opposite of replication task `id` (PULL if it was PUSH and vice versa).
---

# truenas_action_replication_restore (Resource)

Create the opposite of replication task `id` (PULL if it was PUSH and vice versa).

This is an action resource that executes the `replication.restore` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_replication_restore" "example" {
  id = 1
  replication_restore = "value"
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the replication task to restore.
- `replication_restore` (String, Required) Configuration options for restoring the replication task.

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
