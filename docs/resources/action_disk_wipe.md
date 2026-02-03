---
page_title: "truenas_action_disk_wipe Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Disks"
description: |-
  Performs a wipe of a disk `dev`.
---

# truenas_action_disk_wipe (Resource)

Performs a wipe of a disk `dev`.

This is an action resource that executes the `disk.wipe` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_disk_wipe" "example" {
  dev = "value"
  mode = "value"
}
```

## Schema

### Input Parameters

- `dev` (String, Required) The device to perform the disk wipe operation on. May be passed as /dev/sda or just sda.
- `mode` (String, Required) * QUICK: Write zeros to the first and last 32MB of device. * FULL: Write whole disk with zeros. * FULL_RANDOM: Write whole disk with random bytes.
- `synccache` (Bool, Optional) Synchronize the device with the database.

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
