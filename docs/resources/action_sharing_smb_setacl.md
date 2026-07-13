---
page_title: "truenas_action_sharing_smb_setacl Resource - terraform-provider-truenas"
subcategory: "Actions - Sharing"
description: |-
  Set SMB share ACL entries to control user/group access.
---

# truenas_action_sharing_smb_setacl (Resource)

Set SMB share ACL entries to control user/group access.

This is an action resource that executes the `sharing.smb.setacl` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_sharing_smb_setacl" "myshare" {
  smb_setacl = jsonencode({
    share_name = "myshare"
    share_acl  = [
      { ae_perm = "FULL", ae_type = "ALLOWED", ae_who_str = "myuser" }
    ]
  })
}
```

## Schema

### Input Parameters

- `smb_setacl` (String, Required) SharingSMBSetaclArgs parameters.

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
