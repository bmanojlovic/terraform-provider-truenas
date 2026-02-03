---
page_title: "truenas_action_vm_log_file_download Resource - terraform-provider-truenas"
subcategory: "Actions - Virtual Machines"
description: |-
  Retrieve log file contents of `id` VM.  It will download empty file if log file does not exist.
---

# truenas_action_vm_log_file_download (Resource)

Retrieve log file contents of `id` VM.  It will download empty file if log file does not exist.

This is an action resource that executes the `vm.log_file_download` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_vm_log_file_download" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the virtual machine to download log file for.

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
