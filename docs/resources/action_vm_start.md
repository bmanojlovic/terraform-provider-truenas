---
page_title: "truenas_action_vm_start Resource - terraform-provider-truenas"
subcategory: "Actions - Virtual Machines"
description: |-
  Start a VM.  options.overcommit defaults to false, meaning VMs are not allowed to start if there is not enough available memory to hold all configured VMs. If true, VM starts even if there is not enough memory for all configured VMs.  Error codes:      ENOMEM(12): not enough free memory to run the VM without overcommit
---

# truenas_action_vm_start (Resource)

Start a VM.  options.overcommit defaults to false, meaning VMs are not allowed to start if there is not enough available memory to hold all configured VMs. If true, VM starts even if there is not enough memory for all configured VMs.  Error codes:      ENOMEM(12): not enough free memory to run the VM without overcommit

This is an action resource that executes the `vm.start` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_vm_start" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the virtual machine to start.
- `options` (String, Optional) Options controlling the VM start process.

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
