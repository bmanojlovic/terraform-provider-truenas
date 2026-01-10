---
page_title: "truenas_user Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new user.
---

# truenas_user (Resource)

Create a new user.

## Example Usage

```terraform
resource "truenas_user" "example" {
  full_name = "example-value"
  username = "example-value"
}
```

## Schema

### Required

- `full_name` (String) - Comment field to provide additional information about the user account. Typically, this is     the full name of the user or a short description of a service account. There are no character set restric
- `username` (String) - String used to uniquely identify the user on the server. In order to be portable across     systems, local user names must be composed of characters from the POSIX portable filename     character set 

### Optional

- `email` (String) - Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and     notifications. Default: `None`
- `group` (Int64) - The group entry `id` for the user's primary group. This is not the same as the Unix group `gid` value.     This is required if `group_create` is `false`. Default: `None`
- `group_create` (Bool) - If set to `true`, the TrueNAS server automatically creates a new local group as the user's primary group. Default: `False`
- `groups` (List) - Array of additional groups to which the user belongs. NOTE: Groups are identified by their group entry `id`,     not their Unix group ID (`gid`).
- `home` (String) - The local file system path for the user account's home directory. Typically, this is required only if the account has shell access (local or SSH) to TrueNAS. This is not required for accounts used onl Default: `/var/empty`
- `home_create` (Bool) - Create a new home directory for the user in the specified `home` path. Default: `False`
- `home_mode` (String) - Filesystem permission to set on the user's home directory. Default: `700`
- `locked` (Bool) - If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS server. Default: `False`
- `password` (String) - The password for the user account. This is required if `random_password` is not set. Default: `None`
- `password_disabled` (Bool) - If set to `true` password authentication for the user account is disabled.  NOTE: Users with password authentication disabled may still authenticate to the TrueNAS server by other methods,     such as Default: `False`
- `random_password` (Bool) - Generate a random 20 character password for the user. Default: `False`
- `shell` (String) - Available choices can be retrieved with `user.shell_choices`. Default: `/usr/bin/zsh`
- `smb` (Bool) - The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash of the     user account's password for local accounts. This feature is unavailable for local accounts Default: `True`
- `ssh_password_enabled` (Bool) - Allow the user to authenticate to the TrueNAS SSH server using a password.  WARNING: The established best practice is to use only key-based authentication for SSH servers. Default: `False`
- `sshpubkey` (String) - SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server. Default: `None`
- `sudo_commands` (List) - An array of commands the user may execute with elevated privileges. User is prompted for password     when executing any command from the array.
- `sudo_commands_nopasswd` (List) - An array of commands the user may execute with elevated privileges. User is *not* prompted for password     when executing any command from the array.
- `uid` (Int64) - UNIX UID. If not provided, it is automatically filled with the next one available. Default: `None`
- `userns_idmap` (Int64) - Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to all containers. Alternatively, the target UID may be     explicitly specified. If `null`, then the UID Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_user.example <id>
```
