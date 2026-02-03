---
page_title: "truenas_action_support_new_ticket Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Creates a new ticket for support. This is done using the support proxy API. For TrueNAS Community Edition it will be created on JIRA and for TrueNAS Enterprise on Salesforce.  For Community Edition, `criticality`, `environment`, `phone`, `name`, and `email` attributes are not required. For Enterprise, `token` and `type` attributes are not required.
---

# truenas_action_support_new_ticket (Resource)

Creates a new ticket for support. This is done using the support proxy API. For TrueNAS Community Edition it will be created on JIRA and for TrueNAS Enterprise on Salesforce.  For Community Edition, `criticality`, `environment`, `phone`, `name`, and `email` attributes are not required. For Enterprise, `token` and `type` attributes are not required.

This is an action resource that executes the `support.new_ticket` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_support_new_ticket" "example" {
  data = "value"
}
```

## Schema

### Input Parameters

- `data` (String, Required) Support ticket data for either enterprise or community support.

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
