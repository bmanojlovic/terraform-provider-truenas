---
page_title: "truenas_action_pool_scrub_run Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Initiate a scrub of a pool `name` if last scrub was performed more than `threshold` days before.
---

# truenas_action_pool_scrub_run (Resource)

Initiate a scrub of a pool `name` if last scrub was performed more than `threshold` days before.

This is an action resource that executes the `pool.scrub.run` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_scrub_run" "example" {
  name = "value"
}
```

## Schema

### Input Parameters

- `name` (String, Required) Name of the pool to run scrub on.
- `threshold` (Int64, Optional) Days before a scrub is due when the scrub should start.

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
