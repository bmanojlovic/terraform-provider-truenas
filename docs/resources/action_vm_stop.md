---
page_title: "truenas_action_vm_stop Resource - terraform-provider-truenas"
subcategory: "Actions - Virtual Machines"
description: |-
  Stops a VM.  For unresponsive guests who have exceeded the `shutdown_timeout` defined by the user and have become unresponsive, they required to be powered down using `vm.poweroff`. `vm.stop` is only going to send a shutdown signal to the guest and wait the desired `shutdown_timeout` value before tearing down guest vmemory.  `force_after_timeout` when supplied, it will initiate poweroff for the VM forcing it to exit if it has not already stopped within the specified `shutdown_timeout`.
---

# truenas_action_vm_stop (Resource)

Stops a VM.  For unresponsive guests who have exceeded the `shutdown_timeout` defined by the user and have become unresponsive, they required to be powered down using `vm.poweroff`. `vm.stop` is only going to send a shutdown signal to the guest and wait the desired `shutdown_timeout` value before tearing down guest vmemory.  `force_after_timeout` when supplied, it will initiate poweroff for the VM forcing it to exit if it has not already stopped within the specified `shutdown_timeout`.

This is an action resource that executes the `vm.stop` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_vm_stop" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the virtual machine to stop.
- `options` (String, Optional) Options controlling the VM stop process.

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
