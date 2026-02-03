---
page_title: "truenas_action_system_reboot Resource - terraform-provider-truenas"
subcategory: "Actions - System"
description: |-
  Reboots the operating system.  Emits an "added" event of name "system" and id "reboot".
---

# truenas_action_system_reboot (Resource)

Reboots the operating system.  Emits an "added" event of name "system" and id "reboot".

This is an action resource that executes the `system.reboot` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_system_reboot" "example" {
  reason = "value"
}
```

## Schema

### Input Parameters

- `reason` (String, Required) Reason for the system reboot.
- `options` (String, Optional) Options for controlling the reboot process.

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
