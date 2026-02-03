---
page_title: "truenas_action_alert_restore Resource - terraform-provider-truenas"
subcategory: "Actions - Alerts"
description: |-
  Restore `id` alert which had been dismissed.
---

# truenas_action_alert_restore (Resource)

Restore `id` alert which had been dismissed.

This is an action resource that executes the `alert.restore` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_alert_restore" "example" {
  uuid = "value"
}
```

## Schema

### Input Parameters

- `uuid` (String, Required) UUID of the dismissed alert to restore.

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
