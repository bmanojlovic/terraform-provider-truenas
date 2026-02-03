---
page_title: "truenas_action_boot_set_scrub_interval Resource - terraform-provider-truenas"
subcategory: "Actions - System - Boot"
description: |-
  Set Automatic Scrub Interval value in days.
---

# truenas_action_boot_set_scrub_interval (Resource)

Set Automatic Scrub Interval value in days.

This is an action resource that executes the `boot.set_scrub_interval` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_boot_set_scrub_interval" "example" {
  interval = 1
}
```

## Schema

### Input Parameters

- `interval` (Int64, Required) Scrub interval in days (must be a positive integer).

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
