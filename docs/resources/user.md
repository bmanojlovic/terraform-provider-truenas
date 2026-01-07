---
page_title: "truenas_user Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new user account.
---

# truenas_user (Resource)

Configuration for creating a new user account.

## Example Usage

```terraform
resource "truenas_user" "example" {
  username = "example-username"
  full_name = "example-full_name"
  home = "/var/empty"
  shell = "/usr/bin/zsh"
  smb = true
  password_disabled = false
  ssh_password_enabled = false
  locked = false
}
```

## Schema

### Required

- `username` (Required) - String used to uniquely identify the user on the server. In order to be portable across     systems, local user names must be composed of characters from the POSIX portable filename     character set (IEEE Std 1003.1-2024 section 3.265). This means alphanumeric characters,     hyphens, underscores, and periods. Usernames also may not begin with a hyphen or a period.. Type: `string`
- `full_name` (Required) - Comment field to provide additional information about the user account. Typically, this is     the full name of the user or a short description of a service account. There are no character set restrictions     for this field. This field is for information only. . Type: `string`

### Optional

- `uid` (Optional) - UNIX UID. If not provided, it is automatically filled with the next one available.. Type: `string`
- `home` (Optional) - The local file system path for the user account's home directory.
Typically, this is required only if the account has shell access (local or SSH) to TrueNAS.
This is not required for accounts used only for SMB share access.  Default: `/var/empty`. Type: `string`
- `shell` (Optional) - Available choices can be retrieved with `user.shell_choices`. Default: `/usr/bin/zsh`. Type: `string`
- `smb` (Optional) - The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash of the     user account's password for local accounts. This feature is unavailable for local accounts when General Purpose OS     STIG compatibility mode is enabled. If set to `true` the user is automatically added to the `builtin_users`     group. Default: `True`. Type: `boolean`
- `userns_idmap` (Optional) - Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to all containers. Alternatively, the target UID may be     explicitly specified. If `null`, then the UID will not be mapped.

NOTE: This field will be ignored for users that have been assigned TrueNAS roles.. Type: `string`
- `group` (Optional) - The group entry `id` for the user's primary group. This is not the same as the Unix group `gid` value.     This is required if `group_create` is `false`. . Type: `string`
- `groups` (Optional) - Array of additional groups to which the user belongs. NOTE: Groups are identified by their group entry `id`,     not their Unix group ID (`gid`). . Type: `array`
- `password_disabled` (Optional) - If set to `true` password authentication for the user account is disabled.

NOTE: Users with password authentication disabled may still authenticate to the TrueNAS server by other methods,     such as SSH key-based authentication.

NOTE: Password authentication is required for `smb` users. Default: `False`. Type: `boolean`
- `ssh_password_enabled` (Optional) - Allow the user to authenticate to the TrueNAS SSH server using a password.

WARNING: The established best practice is to use only key-based authentication for SSH servers.  Default: `False`. Type: `boolean`
- `sshpubkey` (Optional) - SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server. . Type: `string`
- `locked` (Optional) - If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS server.  Default: `False`. Type: `boolean`
- `sudo_commands` (Optional) - An array of commands the user may execute with elevated privileges. User is prompted for password     when executing any command from the array. . Type: `array`
- `sudo_commands_nopasswd` (Optional) - An array of commands the user may execute with elevated privileges. User is *not* prompted for password     when executing any command from the array. . Type: `array`
- `email` (Optional) - Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and     notifications. . Type: `string`
- `group_create` (Optional) - If set to `true`, the TrueNAS server automatically creates a new local group as the user's primary group.  Default: `False`. Type: `boolean`
- `home_create` (Optional) - Create a new home directory for the user in the specified `home` path.  Default: `False`. Type: `boolean`
- `home_mode` (Optional) - Filesystem permission to set on the user's home directory.  Default: `700`. Type: `string`
- `password` (Optional) - The password for the user account. This is required if `random_password` is not set. . Type: `string`
- `random_password` (Optional) - Generate a random 20 character password for the user. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_user.example <id>
```
