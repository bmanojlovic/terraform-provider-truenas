---
page_title: "truenas_action_kmip_clear_sync_pending_keys Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Clear all keys which are pending to be synced between KMIP server and TN database.  For ZFS/SED keys, we remove the UID from local database with which we are able to retrieve ZFS/SED keys. It should be used with caution.
---

# truenas_action_kmip_clear_sync_pending_keys (Resource)

Clear all keys which are pending to be synced between KMIP server and TN database.  For ZFS/SED keys, we remove the UID from local database with which we are able to retrieve ZFS/SED keys. It should be used with caution.

This is an action resource that executes the `kmip.clear_sync_pending_keys` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_kmip_clear_sync_pending_keys" "example" {
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
