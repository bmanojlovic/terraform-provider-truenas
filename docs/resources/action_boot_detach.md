---
page_title: "truenas_action_boot_detach Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Detach a disk from the boot pool.
---

# truenas_action_boot_detach (Resource)

Detach a disk from the boot pool.

This is an action resource that executes the `boot.detach` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_detach" "example" {
  dev = "value"
}
```

## Schema

### Input Parameters

- `dev` (String, Required) Device name or path to detach from the boot pool.

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
