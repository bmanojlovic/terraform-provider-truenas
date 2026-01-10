---
page_title: "truenas_group Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_group (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_group" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the group to retrieve.

### Read-Only

- `builtin` (Bool) - If `True`, the group is an internal system account for the TrueNAS server. Typically, one should     create dedicated groups for access to the TrueNAS server webui and shares.
- `gid` (Int64) - A non-negative integer used to identify a group. TrueNAS uses this value for permission checks and many other     system purposes.
- `group` (String) - A string used to identify a group. Identical to the `name` key.
- `immutable` (Bool) - This is a read-only field showing if the group entry can be changed. If `True`, the group is immutable and     cannot be changed. If `False`, the group can be changed.
- `local` (Bool) - If `True`, the group is local to the TrueNAS server. If `False`, the group is provided by a directory service.
- `name` (String) - A string used to identify a group.
- `roles` (List) - List of roles assigned to this groups. Roles control administrative access to TrueNAS through the web UI and     API. You can change group roles by using `privilege.create`, `privilege.update`, and `p
- `sid` (String) - The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses this value to     check share access and for other purposes.
- `smb` (Bool) - If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT group account     on the TrueNAS SMB server and has a `sid` value.
- `sudo_commands` (List) - A list of commands that group members may execute with elevated privileges. User is prompted for password     when executing any command from the list.
- `sudo_commands_nopasswd` (List) - A list of commands that group members may execute with elevated privileges. User is not prompted for password     when executing any command from the list.
- `userns_idmap` (Int64) - Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to all containers. Alternatively, the target GID may be     explicitly specified. If null, then the GID 
- `users` (List) - A list a API user identifiers for local users who are members of this group. These IDs match the `id` field     from `user.query`.  NOTE: This field is empty for groups that come from directory servic
