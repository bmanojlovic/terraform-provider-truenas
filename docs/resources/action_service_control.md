---
page_title: "truenas_action_service_control Resource - terraform-provider-truenas"
subcategory: "Actions - Services"
description: |-
  Execute service.control
---

# truenas_action_service_control (Resource)

Execute service.control

This is an action resource that executes the `service.control` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_service_control" "example" {
  verb = "value"
  service = "value"
}
```

## Schema

### Input Parameters

- `verb` (String, Required) The service operation to perform.
- `service` (String, Required) Name of the service to control.
- `options` (String, Optional) Options for controlling the service operation behavior.

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
