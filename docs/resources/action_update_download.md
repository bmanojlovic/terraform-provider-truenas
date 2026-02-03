---
page_title: "truenas_action_update_download Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Download updates.
---

# truenas_action_update_download (Resource)

Download updates.

This is an action resource that executes the `update.download` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_update_download" "example" {
}
```

## Schema

### Input Parameters

- `train` (String, Optional) Specifies the train from which to download the update. If both `train` and `version` are `null``, the most     recent version that matches the currently selected update profile is used.
- `version` (String, Optional) Specific version to download. `null` to download the latest version from the specified train.

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
