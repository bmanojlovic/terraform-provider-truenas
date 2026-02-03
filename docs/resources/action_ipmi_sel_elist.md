---
page_title: "truenas_action_ipmi_sel_elist Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Query IPMI System Event Log (SEL) extended list
---

# truenas_action_ipmi_sel_elist (Resource)

Query IPMI System Event Log (SEL) extended list

This is an action resource that executes the `ipmi.sel.elist` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_ipmi_sel_elist" "example" {
}
```

## Schema

### Input Parameters

- `filters` (List, Optional) List of filters for query results. See API documentation for "Query Methods" for more guidance.
- `options` (String, Optional) Query options including pagination, ordering, and additional parameters.

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
