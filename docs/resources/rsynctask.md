---
page_title: "truenas_rsynctask Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a Rsync Task.
---

# truenas_rsynctask (Resource)

Create a Rsync Task.

## Example Usage

```terraform
resource "truenas_rsynctask" "example" {
  path = "example-value"
  user = "example-value"
}
```

## Schema

### Required

- `path` (String) - Local filesystem path to synchronize.
- `user` (String) - Username to run the rsync task as.

### Optional

- `archive` (Bool) - Make rsync run recursively, preserving symlinks, permissions, modification times, group, and special files. Default: `False`
- `compress` (Bool) - Reduce the size of the data to be transmitted. Default: `True`
- `delayupdates` (Bool) - Delay updating destination files until all transfers are complete. Default: `True`
- `delete` (Bool) - Delete files in the destination directory that do not exist in the source directory. Default: `False`
- `desc` (String) - Description of the rsync task. Default: ``
- `direction` (String) - Specify if data should be PULLED or PUSHED from the remote system. Default: `PUSH` Valid values: `PULL`, `PUSH`
- `enabled` (Bool) - Whether this rsync task is enabled. Default: `True`
- `extra` (List) - Array of additional rsync command-line options.
- `mode` (String) - Operating mechanism for Rsync, i.e. Rsync Module mode or Rsync SSH mode. Default: `MODULE` Valid values: `MODULE`, `SSH`
- `preserveattr` (Bool) - Preserve extended attributes of files. Default: `False`
- `preserveperm` (Bool) - Preserve original file permissions. Default: `False`
- `quiet` (Bool) - Suppress informational messages from rsync. Default: `False`
- `recursive` (Bool) - Recursively transfer subdirectories. Default: `True`
- `remotehost` (String) - IP address or hostname of the remote system. If username differs on the remote host, "username@remote_host"     format should be used. Default: `None`
- `remotemodule` (String) - Name of remote module, this attribute should be specified when `mode` is set to MODULE. Default: `None`
- `remotepath` (String) - Path on the remote system to synchronize with. Default: ``
- `remoteport` (Int64) - Port number for SSH connection. Only applies when `mode` is SSH. Default: `None`
- `schedule` (String) - Cron schedule for when the rsync task should run.
- `ssh_credentials` (Int64) - Keychain credential ID for SSH authentication. `null` to use user's SSH keys. Default: `None`
- `ssh_keyscan` (Bool) - Automatically add remote host key to user's known_hosts file. Default: `False`
- `times` (Bool) - Preserve modification times of files. Default: `True`
- `validate_rpath` (Bool) - Validate the existence of the remote path. Default: `True`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_rsynctask.example <id>
```
