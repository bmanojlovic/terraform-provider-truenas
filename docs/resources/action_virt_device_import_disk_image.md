---
page_title: "truenas_action_virt_device_import_disk_image Resource - terraform-provider-truenas"
subcategory: "Actions - Virtualization"
description: |-
  Imports a specified disk image.  Utilized qemu-img with the auto-detect functionality to auto-convert any supported disk image format to RAW -> ZVOL  As of this implementation it supports:  - QCOW2 - QED - RAW - VDI - VPC - VMDK  `diskimg` is a required parameter for the incoming disk image `zvol` is the required target for the imported disk image
---

# truenas_action_virt_device_import_disk_image (Resource)

Imports a specified disk image.  Utilized qemu-img with the auto-detect functionality to auto-convert any supported disk image format to RAW -> ZVOL  As of this implementation it supports:  - QCOW2 - QED - RAW - VDI - VPC - VMDK  `diskimg` is a required parameter for the incoming disk image `zvol` is the required target for the imported disk image

This is an action resource that executes the `virt.device.import_disk_image` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_virt_device_import_disk_image" "example" {
  virt_device_import_disk_image = "value"
}
```

## Schema

### Input Parameters

- `virt_device_import_disk_image` (String, Required) VirtDeviceImportDiskImageArgs parameters.

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
