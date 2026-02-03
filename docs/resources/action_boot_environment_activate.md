---
page_title: "truenas_action_boot_environment_activate Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Activate a boot environment for next reboot.
---

# truenas_action_boot_environment_activate (Resource)

Activate a boot environment for next reboot.

This is an action resource that executes the `boot.environment.activate` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_environment_activate" "switch" {
  boot_environment_activate = jsonencode({
    id = "24.10-MASTER-20241201-015106"
  })
}
```

## Schema

### Input Parameters

- `boot_environment_activate` (String, Required) BootEnvironmentActivateArgs parameters.

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
