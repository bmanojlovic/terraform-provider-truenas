---
page_title: "truenas_action_boot_attach Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Attach a disk to the boot pool, turning a stripe into a mirror.  `expand` option will determine whether the new disk partition will be          the maximum available or the same size as the current disk.
---

# truenas_action_boot_attach (Resource)

Attach a disk to the boot pool, turning a stripe into a mirror.  `expand` option will determine whether the new disk partition will be          the maximum available or the same size as the current disk.

This is an action resource that executes the `boot.attach` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_attach" "example" {
  dev = "value"
}
```

## Schema

### Input Parameters

- `dev` (String, Required) Device name or path to attach to the boot pool.
- `options` (String, Optional) Options for the attach operation.

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
