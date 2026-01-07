---
page_title: "truenas_group Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Group configuration data for the new group.
---

# truenas_group (Resource)

Group configuration data for the new group.

## Example Usage

```terraform
resource "truenas_group" "example" {
  name = "example-name"
  sudo_commands = []
  sudo_commands_nopasswd = []
  smb = true
  users = []
}
```

## Schema

### Required

- `name` (Required) - name configuration. Type: `string`

### Optional

- `gid` (Optional) - If `null`, it is automatically filled with the next one available.. Type: `string`
- `sudo_commands` (Optional) - A list of commands that group members may execute with elevated privileges. User is prompted for password     when executing any command from the list.  Default: `[]`. Type: `array`
- `sudo_commands_nopasswd` (Optional) - A list of commands that group members may execute with elevated privileges. User is not prompted for password     when executing any command from the list.  Default: `[]`. Type: `array`
- `smb` (Optional) - If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT group account     on the TrueNAS SMB server and has a `sid` value.  Default: `True`. Type: `boolean`
- `userns_idmap` (Optional) - Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to all containers. Alternatively, the target GID may be     explicitly specified. If null, then the GID will not be mapped.

**NOTE: This field will be ignored for groups that have been assigned TrueNAS roles.**. Type: `string`
- `users` (Optional) - A list a API user identifiers for local users who are members of this group. These IDs match the `id` field     from `user.query`.

NOTE: This field is empty for groups that come from directory services (`local` is `False`).  Default: `[]`. Type: `array`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_group.example <id>
```
