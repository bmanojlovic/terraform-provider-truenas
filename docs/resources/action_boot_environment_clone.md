---
page_title: "truenas_action_boot_environment_clone Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Clone current boot environment.
---

# truenas_action_boot_environment_clone (Resource)

Clone current boot environment.

This is an action resource that executes the `boot.environment.clone` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_environment_clone" "backup" {
  boot_environment_clone = jsonencode({
    id     = "24.10-MASTER-20241201-015106"
    target = "backup-before-upgrade"
  })
}
```

## Schema

### Input Parameters

- `boot_environment_clone` (String, Required) BootEnvironmentCloneArgs parameters.

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
