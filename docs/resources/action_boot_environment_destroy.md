---
page_title: "truenas_action_boot_environment_destroy Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Delete a boot environment.
---

# truenas_action_boot_environment_destroy (Resource)

Delete a boot environment.

This is an action resource that executes the `boot.environment.destroy` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_environment_destroy" "example" {
  boot_environment_destroy = "value"
}
```

## Schema

### Input Parameters

- `boot_environment_destroy` (String, Required) BootEnvironmentDestroyArgs parameters.

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
