---
page_title: "truenas_action_pool_ddt_prefetch Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Prefetch DDT entries in pool `pool_name`.
---

# truenas_action_pool_ddt_prefetch (Resource)

Prefetch DDT entries in pool `pool_name`.

This is an action resource that executes the `pool.ddt_prefetch` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_ddt_prefetch" "example" {
  pool_name = "value"
}
```

## Schema

### Input Parameters

- `pool_name` (String, Required) Name of the pool to prefetch deduplication table entries for.

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
