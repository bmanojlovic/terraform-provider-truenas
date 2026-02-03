---
page_title: "truenas_action_vm_poweroff Resource - terraform-provider-truenas"
subcategory: "Actions - Virtual Machines"
description: |-
  Force power off a VM.
---

# truenas_action_vm_poweroff (Resource)

Force power off a VM.

This is an action resource that executes the `vm.poweroff` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_vm_poweroff" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the virtual machine to power off forcefully.

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
