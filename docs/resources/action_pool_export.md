---
page_title: "truenas_action_pool_export Resource - terraform-provider-truenas"
subcategory: "Actions - Storage - Pools"
description: |-
  Export pool of `id`.  `cascade` will delete all attachments of the given pool (`pool.attachments`). `restart_services` will restart services that have open files on given pool. `destroy` will also PERMANENTLY destroy the pool/data.
---

# truenas_action_pool_export (Resource)

Export pool of `id`.  `cascade` will delete all attachments of the given pool (`pool.attachments`). `restart_services` will restart services that have open files on given pool. `destroy` will also PERMANENTLY destroy the pool/data.

This is an action resource that executes the `pool.export` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_pool_export" "example" {
  id = 1
}
```

## Schema

### Input Parameters

- `id` (Int64, Required) ID of the pool to export.
- `options` (String, Optional) Options for controlling the pool export process.

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
