---
page_title: "truenas_action_service_started_or_enabled Resource - terraform-provider-truenas"
subcategory: "Actions - Services"
description: |-
  Test if service specified by `service` is started or enabled to start automatically.
---

# truenas_action_service_started_or_enabled (Resource)

Test if service specified by `service` is started or enabled to start automatically.

This is an action resource that executes the `service.started_or_enabled` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_service_started_or_enabled" "example" {
  service = "value"
}
```

## Schema

### Input Parameters

- `service` (String, Required) Name of the service to check if running or enabled.

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
