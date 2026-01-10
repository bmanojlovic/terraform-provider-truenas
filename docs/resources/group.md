---
page_title: "truenas_group Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new group.
---

# truenas_group (Resource)

Create a new group.

## Example Usage

```terraform
resource "truenas_group" "example" {
  name = "example-value"
}
```

## Schema

### Required

- `name` (String) - 

### Optional

- `gid` (Int64) - If `null`, it is automatically filled with the next one available. Default: `None`
- `smb` (Bool) - If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT group account     on the TrueNAS SMB server and has a `sid` value. Default: `True`
- `sudo_commands` (List) - A list of commands that group members may execute with elevated privileges. User is prompted for password     when executing any command from the list. Default: `[]`
- `sudo_commands_nopasswd` (List) - A list of commands that group members may execute with elevated privileges. User is not prompted for password     when executing any command from the list. Default: `[]`
- `userns_idmap` (Int64) - Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to all containers. Alternatively, the target GID may be     explicitly specified. If null, then the GID  Default: `None`
- `users` (List) - A list a API user identifiers for local users who are members of this group. These IDs match the `id` field     from `user.query`.  NOTE: This field is empty for groups that come from directory servic Default: `[]`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_group.example <id>
```
