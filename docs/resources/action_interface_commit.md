---
page_title: "truenas_action_interface_commit Resource - terraform-provider-truenas"
subcategory: "Actions - Network"
description: |-
  Commit pending network changes. Starts rollback timer.
---

# truenas_action_interface_commit (Resource)

Commit pending network changes. Starts rollback timer.

## Workflow

1. Make interface changes with truenas_interface
2. Call commit (starts checkin_timeout rollback timer)
3. Call checkin within timeout to confirm changes

This is an action resource that executes the `interface.commit` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_interface_commit" "apply" {
  depends_on       = [truenas_interface.vlan10]
  checkin_timeout  = 60
}
```

## Schema

### Input Parameters

- `options` (String, Optional) Options for committing interface changes.

### Computed Outputs

- `action_id` (String) Unique identifier for this action execution
- `job_id` (Int64) Background job ID (if applicable)
- `state` (String) Job state: SUCCESS, FAILED, or RUNNING
- `progress` (Float64) Job progress percentage (0-100)
- `result` (String) Action result data
- `error` (String) Error message if action failed

## Related Actions

- `truenas_action_interface_checkin`
- `truenas_action_interface_rollback`

## Notes

- Actions execute immediately when the resource is created
- Background jobs are monitored until completion
- Progress updates are logged during execution
- The resource cannot be updated - changes force recreation
- Destroying the resource does not undo the action
