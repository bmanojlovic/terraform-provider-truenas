---
page_title: "truenas_action_boot_replace Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Replace device `label` on boot pool with `dev`.
---

# truenas_action_boot_replace (Resource)

Replace device `label` on boot pool with `dev`.

This is an action resource that executes the `boot.replace` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_replace" "example" {
  label = "value"
  dev = "value"
}
```

## Schema

### Input Parameters

- `label` (String, Required) Label of the disk in the boot pool to replace.
- `dev` (String, Required) Device name or path of the replacement disk.

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
