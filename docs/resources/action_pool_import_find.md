---
page_title: "truenas_action_pool_import_find Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Returns a job id which can be used to retrieve a list of pools available for import with the following details as a result of the job: name, guid, status, hostname.
---

# truenas_action_pool_import_find (Resource)

Returns a job id which can be used to retrieve a list of pools available for import with the following details as a result of the job: name, guid, status, hostname.

This is an action resource that executes the `pool.import_find` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_import_find" "example" {
}
```

## Schema

### Input Parameters

None

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
