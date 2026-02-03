---
page_title: "truenas_action_directoryservices_cache_refresh Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  This method refreshes the directory services cache for users and groups that is used as a backing for `user.query` and `group.query` methods. The first cache fill in an Active Directory domain may take a significant amount of time to complete and so it is performed as within a job. The most likely situation in which a user may desire to refresh the directory services cache is after new users or groups  to a remote directory server with the intention to have said users or groups appear in the results of the aforementioned account-related methods.  A cache refresh is not required in order to use newly-added users and groups for in permissions and ACL related methods. Likewise, a cache refresh will not resolve issues with users being unable to authenticate to shares.
---

# truenas_action_directoryservices_cache_refresh (Resource)

This method refreshes the directory services cache for users and groups that is used as a backing for `user.query` and `group.query` methods. The first cache fill in an Active Directory domain may take a significant amount of time to complete and so it is performed as within a job. The most likely situation in which a user may desire to refresh the directory services cache is after new users or groups  to a remote directory server with the intention to have said users or groups appear in the results of the aforementioned account-related methods.  A cache refresh is not required in order to use newly-added users and groups for in permissions and ACL related methods. Likewise, a cache refresh will not resolve issues with users being unable to authenticate to shares.

This is an action resource that executes the `directoryservices.cache_refresh` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_directoryservices_cache_refresh" "example" {
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
