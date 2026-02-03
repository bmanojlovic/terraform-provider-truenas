---
page_title: "truenas_action_update_file Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Updates the system using the uploaded .tar file.
---

# truenas_action_update_file (Resource)

Updates the system using the uploaded .tar file.

This is an action resource that executes the `update.file` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_update_file" "example" {
}
```

## Schema

### Input Parameters

- `options` (String, Optional) Options for controlling the manual update file upload process.

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
