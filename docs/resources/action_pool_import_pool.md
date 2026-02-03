---
page_title: "truenas_action_pool_import_pool Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Import a pool found with `pool.import_find`.  Errors:     ENOENT - Pool not found
---

# truenas_action_pool_import_pool (Resource)

Import a pool found with `pool.import_find`.  Errors:     ENOENT - Pool not found

This is an action resource that executes the `pool.import_pool` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_import_pool" "example" {
  pool_import = "value"
}
```

## Schema

### Input Parameters

- `pool_import` (String, Required) PoolImportPoolArgs parameters.

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
