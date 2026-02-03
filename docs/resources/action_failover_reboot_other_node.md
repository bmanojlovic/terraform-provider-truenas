---
page_title: "truenas_action_failover_reboot_other_node Resource - terraform-provider-truenas"
subcategory: "Actions - Other"
description: |-
  Reboot the other node and wait for it to come back online.  NOTE: This makes very few checks on HA systems. You need to     know what you're doing before calling this.
---

# truenas_action_failover_reboot_other_node (Resource)

Reboot the other node and wait for it to come back online.  NOTE: This makes very few checks on HA systems. You need to     know what you're doing before calling this.

This is an action resource that executes the `failover.reboot.other_node` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_failover_reboot_other_node" "example" {
}
```

## Schema

### Input Parameters

- `options` (String, Optional) Options for rebooting the other node.

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
