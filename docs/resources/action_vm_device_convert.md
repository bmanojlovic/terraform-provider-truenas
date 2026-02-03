---
page_title: "truenas_action_vm_device_convert Resource - terraform-provider-truenas"
subcategory: "Actions - Virtual Machines"
description: |-
  Convert between disk images and ZFS volumes. Supported disk image formats         are qcow2, qed, raw, vdi, vhdx, and vmdk. The conversion direction is determined         automatically based on file extension.
---

# truenas_action_vm_device_convert (Resource)

Convert between disk images and ZFS volumes. Supported disk image formats         are qcow2, qed, raw, vdi, vhdx, and vmdk. The conversion direction is determined         automatically based on file extension.

This is an action resource that executes the `vm.device.convert` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_vm_device_convert" "example" {
  vm_convert = "value"
}
```

## Schema

### Input Parameters

- `vm_convert` (String, Required) VMDeviceConvertArgs parameters.

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
