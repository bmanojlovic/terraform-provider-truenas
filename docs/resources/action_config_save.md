---
page_title: "truenas_action_config_save Resource - terraform-provider-truenas"
subcategory: "Actions - System - Config"
description: |-
  Create a tar file of security-sensitive information. These options select which information is included in the tar file:  `secretseed` bool: When true, include password secret seed. `pool_keys` bool: IGNORED and DEPRECATED as it does not apply on SCALE systems. `root_authorized_keys` bool: When true, include "/root/.ssh/authorized_keys" file for the root user.  If none of these options are set, the tar file is not generated and the database file is returned.
---

# truenas_action_config_save (Resource)

Create a tar file of security-sensitive information. These options select which information is included in the tar file:  `secretseed` bool: When true, include password secret seed. `pool_keys` bool: IGNORED and DEPRECATED as it does not apply on SCALE systems. `root_authorized_keys` bool: When true, include "/root/.ssh/authorized_keys" file for the root user.  If none of these options are set, the tar file is not generated and the database file is returned.

This is an action resource that executes the `config.save` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_config_save" "example" {
}
```

## Schema

### Input Parameters

- `options` (String, Optional) Options controlling what data to include in the configuration backup.

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
