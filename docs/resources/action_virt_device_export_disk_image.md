---
page_title: "truenas_action_virt_device_export_disk_image Resource - terraform-provider-truenas"
subcategory: "Actions - Virtualization"
description: |-
  Exports a zvol to a formatted VM disk image.  Utilized qemu-img with the conversion functionality to export a zvol to any supported disk image format, from RAW -> ${OTHER}. The resulting file will be set to inherit the permissions of the target directory.  As of this implementation it supports the following {format} options :  - QCOW2 - QED - RAW - VDI - VPC - VMDK  `format` is a required parameter for the exported disk image `directory` is a required parameter for the export disk image `zvol` is the source for the disk image
---

# truenas_action_virt_device_export_disk_image (Resource)

Exports a zvol to a formatted VM disk image.  Utilized qemu-img with the conversion functionality to export a zvol to any supported disk image format, from RAW -> ${OTHER}. The resulting file will be set to inherit the permissions of the target directory.  As of this implementation it supports the following {format} options :  - QCOW2 - QED - RAW - VDI - VPC - VMDK  `format` is a required parameter for the exported disk image `directory` is a required parameter for the export disk image `zvol` is the source for the disk image

This is an action resource that executes the `virt.device.export_disk_image` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_virt_device_export_disk_image" "example" {
  virt_device_export_disk_image = "value"
}
```

## Schema

### Input Parameters

- `virt_device_export_disk_image` (String, Required) VirtDeviceExportDiskImageArgs parameters.

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
