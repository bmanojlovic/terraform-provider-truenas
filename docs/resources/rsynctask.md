---
page_title: "truenas_rsynctask Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new rsync task.
---

# truenas_rsynctask (Resource)

Configuration for creating a new rsync task.

## Example Usage

```terraform
resource "truenas_rsynctask" "example" {
  path = "example-path"
  user = "example-user"
  mode = "MODULE"
  remotepath = ""
  direction = "PULL"
  desc = ""
  recursive = true
  times = true
}
```

## Schema

### Required

- `path` (Required) - Local filesystem path to synchronize.. Type: `string`
- `user` (Required) - Username to run the rsync task as.. Type: `string`

### Optional

- `mode` (Optional) - Operating mechanism for Rsync, i.e. Rsync Module mode or Rsync SSH mode. Valid values: `MODULE`, `SSH` Default: `MODULE`. Type: `string`
- `remotehost` (Optional) - IP address or hostname of the remote system. If username differs on the remote host, "username@remote_host"     format should be used.. Type: `string`
- `remoteport` (Optional) - Port number for SSH connection. Only applies when `mode` is SSH.. Type: `string`
- `remotemodule` (Optional) - Name of remote module, this attribute should be specified when `mode` is set to MODULE.. Type: `string`
- `ssh_credentials` (Optional) - Keychain credential ID for SSH authentication. `null` to use user's SSH keys.. Type: `string`
- `remotepath` (Optional) - Path on the remote system to synchronize with. Default: ``. Type: `string`
- `direction` (Optional) - Specify if data should be PULLED or PUSHED from the remote system. Valid values: `PULL`, `PUSH` Default: `PUSH`. Type: `string`
- `desc` (Optional) - Description of the rsync task. Default: ``. Type: `string`
- `schedule` (Optional) - Cron schedule for when the rsync task should run.. Type: `object`
- `recursive` (Optional) - Recursively transfer subdirectories. Default: `True`. Type: `boolean`
- `times` (Optional) - Preserve modification times of files. Default: `True`. Type: `boolean`
- `compress` (Optional) - Reduce the size of the data to be transmitted. Default: `True`. Type: `boolean`
- `archive` (Optional) - Make rsync run recursively, preserving symlinks, permissions, modification times, group, and special files. Default: `False`. Type: `boolean`
- `delete` (Optional) - Delete files in the destination directory that do not exist in the source directory. Default: `False`. Type: `boolean`
- `quiet` (Optional) - Suppress informational messages from rsync. Default: `False`. Type: `boolean`
- `preserveperm` (Optional) - Preserve original file permissions. Default: `False`. Type: `boolean`
- `preserveattr` (Optional) - Preserve extended attributes of files. Default: `False`. Type: `boolean`
- `delayupdates` (Optional) - Delay updating destination files until all transfers are complete. Default: `True`. Type: `boolean`
- `extra` (Optional) - Array of additional rsync command-line options.. Type: `array`
- `enabled` (Optional) - Whether this rsync task is enabled. Default: `True`. Type: `boolean`
- `validate_rpath` (Optional) - Validate the existence of the remote path. Default: `True`. Type: `boolean`
- `ssh_keyscan` (Optional) - Automatically add remote host key to user's known_hosts file. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_rsynctask.example <id>
```
