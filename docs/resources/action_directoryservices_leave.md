---
page_title: "truenas_action_directoryservices_leave Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Leave an Active Directory or IPA domain. Calling this endpoint when the directory services status is `HEALTHY` will cause TrueNAS to remove its account from the domain and then reset the local directory services configuration on TrueNAS.
---

# truenas_action_directoryservices_leave (Resource)

Leave an Active Directory or IPA domain. Calling this endpoint when the directory services status is `HEALTHY` will cause TrueNAS to remove its account from the domain and then reset the local directory services configuration on TrueNAS.

This is an action resource that executes the `directoryservices.leave` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_directoryservices_leave" "example" {
  credential = "value"
}
```

## Schema

### Input Parameters

- `credential` (String, Required) DirectoryServicesLeaveArgs parameters.

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
