---
page_title: "truenas_action_system_general_ui_restart Resource - terraform-provider-truenas"
subcategory: "Actions - System"
description: |-
  Restart HTTP server to use latest UI settings.  HTTP server will be restarted after `delay` seconds.
---

# truenas_action_system_general_ui_restart (Resource)

Restart HTTP server to use latest UI settings.  HTTP server will be restarted after `delay` seconds.

This is an action resource that executes the `system.general.ui_restart` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_system_general_ui_restart" "example" {
}
```

## Schema

### Input Parameters

- `delay` (Int64, Optional) How long to wait before the UI is restarted.

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
