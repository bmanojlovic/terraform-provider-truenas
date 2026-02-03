---
page_title: "truenas_action_virt_instance_stop Resource - terraform-provider-truenas"
subcategory: "Actions - Virtualization"
description: |-
  Stop an instance.  Timeout is how long it should wait for the instance to shutdown cleanly.
---

# truenas_action_virt_instance_stop (Resource)

Stop an instance.  Timeout is how long it should wait for the instance to shutdown cleanly.

This is an action resource that executes the `virt.instance.stop` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_virt_instance_stop" "example" {
  id = "value"
}
```

## Schema

### Input Parameters

- `id` (String, Required) ID of the virtual instance to stop.
- `stop_args` (String, Optional) Arguments controlling how the instance is stopped.

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
