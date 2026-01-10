---
page_title: "truenas_user Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_user (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_user" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the user to retrieve.

### Read-Only

- `api_keys` (List) - Array of API key IDs associated with this user account for programmatic access.
- `builtin` (Bool) - If `true`, the user account is an internal system account for the TrueNAS server. Typically, one should     create dedicated user accounts for access to the TrueNAS server webui and shares.
- `email` (String) - Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and     notifications.
- `full_name` (String) - Comment field to provide additional information about the user account. Typically, this is     the full name of the user or a short description of a service account. There are no character set restric
- `group` (String) - The primary group of the user account.
- `groups` (List) - Array of additional groups to which the user belongs. NOTE: Groups are identified by their group entry `id`,     not their Unix group ID (`gid`).
- `home` (String) - The local file system path for the user account's home directory. Typically, this is required only if the account has shell access (local or SSH) to TrueNAS. This is not required for accounts used onl
- `immutable` (Bool) - If `true`, the account is system-provided and most fields related to it may not be changed.
- `last_password_change` (String) - The date of the last password change for local user accounts.
- `local` (Bool) - If `true`, the account is local to the TrueNAS server. If `false`, the account is provided by a directory     service.
- `locked` (Bool) - If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS server.
- `password_age` (Int64) - The age in days of the password for local user accounts.
- `password_change_required` (Bool) - Password change for local user account is required on next login.
- `password_disabled` (Bool) - If set to `true` password authentication for the user account is disabled.  NOTE: Users with password authentication disabled may still authenticate to the TrueNAS server by other methods,     such as
- `password_history` (String) - This contains hashes of the ten most recent passwords used by local user accounts, and is     for enforcing password history requirements as defined in system.security.
- `roles` (List) - Array of roles assigned to this user's groups. Roles control administrative access to TrueNAS through     the web UI and API.
- `shell` (String) - Available choices can be retrieved with `user.shell_choices`.
- `sid` (String) - The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses     this value to check share access and for other purposes.
- `smb` (Bool) - The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash of the     user account's password for local accounts. This feature is unavailable for local accounts
- `smbhash` (String) - NT hash of the local account password for `smb` users. This value is `null` for accounts provided by directory     services or non-SMB accounts.
- `ssh_password_enabled` (Bool) - Allow the user to authenticate to the TrueNAS SSH server using a password.  WARNING: The established best practice is to use only key-based authentication for SSH servers.
- `sshpubkey` (String) - SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server.
- `sudo_commands` (List) - An array of commands the user may execute with elevated privileges. User is prompted for password     when executing any command from the array.
- `sudo_commands_nopasswd` (List) - An array of commands the user may execute with elevated privileges. User is *not* prompted for password     when executing any command from the array.
- `twofactor_auth_configured` (Bool) - If `true`, the account has been configured for two-factor authentication. Users are prompted for a     second factor when authenticating to the TrueNAS web UI and API. They may also be prompted when s
- `uid` (Int64) - A non-negative integer used to identify a system user. TrueNAS uses this value for permission     checks and many other system purposes.
- `unixhash` (String) - Hashed password for local accounts. This value is `null` for accounts provided by directory services.
- `username` (String) - A string used to identify a user. Local accounts must use characters from the POSIX portable filename     character set.
- `userns_idmap` (Int64) - Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to all containers. Alternatively, the target UID may be     explicitly specified. If `null`, then the UID
