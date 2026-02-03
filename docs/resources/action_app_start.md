---
page_title: "truenas_action_app_start Resource - terraform-provider-truenas"
subcategory: "Actions - Applications"
description: |-
  Start `app_name` app.
---

# truenas_action_app_start (Resource)

Start `app_name` app.

This is an action resource that executes the `app.start` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_app_start" "example" {
  app_name = "value"
}
```

## Schema

### Input Parameters

- `app_name` (String, Required) Name of the application to start.

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
