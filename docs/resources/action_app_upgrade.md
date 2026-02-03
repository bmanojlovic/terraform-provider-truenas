---
page_title: "truenas_action_app_upgrade Resource - terraform-provider-truenas"
subcategory: "Actions - Applications"
description: |-
  Upgrade `app_name` app to `app_version`.
---

# truenas_action_app_upgrade (Resource)

Upgrade `app_name` app to `app_version`.

This is an action resource that executes the `app.upgrade` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_app_upgrade" "example" {
  app_name = "value"
}
```

## Schema

### Input Parameters

- `app_name` (String, Required) Name of the application to upgrade.
- `options` (String, Optional) Options controlling the upgrade process including target version and snapshot behavior.

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
