---
page_title: "truenas_action_config_reset Resource - terraform-provider-truenas"
subcategory: "Actions - System - Config"
description: |-
  Reset database to configuration defaults.  If `reboot` is true this job will reboot the system after its completed with a delay of 10 seconds.
---

# truenas_action_config_reset (Resource)

Reset database to configuration defaults.  If `reboot` is true this job will reboot the system after its completed with a delay of 10 seconds.

This is an action resource that executes the `config.reset` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_config_reset" "example" {
}
```

## Schema

### Input Parameters

- `options` (String, Optional) Options controlling the configuration reset behavior.

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
