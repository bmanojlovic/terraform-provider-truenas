---
page_title: "truenas_action_cronjob_run Resource - terraform-provider-truenas"
subcategory: "Actions - Scheduled Tasks"
description: |-
  Job to run cronjob task of `id`.
---

# truenas_action_cronjob_run (Resource)

Job to run cronjob task of `id`.

This is an action resource that executes the `cronjob.run` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_cronjob_run" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the cron job to run immediately.
- `skip_disabled` (Bool, Optional) Whether to skip execution if the cron job is disabled.

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
