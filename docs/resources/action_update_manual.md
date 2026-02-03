---
page_title: "truenas_action_update_manual Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Update the system using a manual update file.
---

# truenas_action_update_manual (Resource)

Update the system using a manual update file.

This is an action resource that executes the `update.manual` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_update_manual" "example" {
  path = "value"
}
```

## Schema

### Input Parameters

- `path` (String, Required) The absolute path to the update file.
- `options` (String, Optional) Options for controlling the manual update process.

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
