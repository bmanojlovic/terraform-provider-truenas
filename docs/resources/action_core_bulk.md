---
page_title: "truenas_action_core_bulk Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Will sequentially call `method` with arguments from the `params` list. For example, running      call("core.bulk", "zfs.snapshot.delete", [["tank@snap-1", true], ["tank@snap-2", false]])  will call      call("zfs.snapshot.delete", "tank@snap-1", true)     call("zfs.snapshot.delete", "tank@snap-2", false)  If the first call fails and the seconds succeeds (returning `true`), the result of the overall call will be:      [         {"result": null, "error": "Error deleting snapshot"},         {"result": true, "error": null}     ]  Important note: the execution status of `core.bulk` will always be a `SUCCESS` (unless an unlikely internal error occurs). Caller must check for individual call results to ensure the absence of any call errors.
---

# truenas_action_core_bulk (Resource)

Will sequentially call `method` with arguments from the `params` list. For example, running      call("core.bulk", "zfs.snapshot.delete", [["tank@snap-1", true], ["tank@snap-2", false]])  will call      call("zfs.snapshot.delete", "tank@snap-1", true)     call("zfs.snapshot.delete", "tank@snap-2", false)  If the first call fails and the seconds succeeds (returning `true`), the result of the overall call will be:      [         {"result": null, "error": "Error deleting snapshot"},         {"result": true, "error": null}     ]  Important note: the execution status of `core.bulk` will always be a `SUCCESS` (unless an unlikely internal error occurs). Caller must check for individual call results to ensure the absence of any call errors.

This is an action resource that executes the `core.bulk` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_core_bulk" "example" {
  method = "value"
  params = "value"
}
```

## Schema

### Input Parameters

- `method` (String, Required) Method name to execute for each parameter set.
- `params` (List, Required) Array of parameter arrays, each representing one method call.
- `description` (String, Optional) Format string for job progress (e.g. "Deleting snapshot {0[dataset]}@{0[name]}").

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
