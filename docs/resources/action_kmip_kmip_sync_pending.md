---
page_title: "truenas_action_kmip_kmip_sync_pending Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Returns true or false based on if there are keys which are to be synced from local database to remote KMIP server or vice versa.
---

# truenas_action_kmip_kmip_sync_pending (Resource)

Returns true or false based on if there are keys which are to be synced from local database to remote KMIP server or vice versa.

This is an action resource that executes the `kmip.kmip_sync_pending` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_kmip_kmip_sync_pending" "example" {
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
