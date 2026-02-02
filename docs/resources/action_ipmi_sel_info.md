---
page_title: "truenas_action_ipmi_sel_info Resource - terraform-provider-truenas"
subcategory: "Actions"
description: |-
  Query General information about the IPMI System Event Log
---

# truenas_action_ipmi_sel_info (Resource)

Query General information about the IPMI System Event Log

This is an action resource that executes the `ipmi.sel.info` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_ipmi_sel_info" "example" {
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
