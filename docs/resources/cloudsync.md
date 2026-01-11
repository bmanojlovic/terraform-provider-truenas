---
page_title: "truenas_cloudsync Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Creates a new cloud_sync entry.
---

# truenas_cloudsync (Resource)

Creates a new cloud_sync entry.


## Example Usage

```terraform
resource "truenas_cloudsync" "example" {
  attributes = "example-value"
  credentials = 1
  direction = "example-value"
  path = "example-value"
  transfer_mode = "example-value"
}
```

## Schema

### Required

- `attributes` (String) - Additional information for each backup, e.g. bucket name.
- `credentials` (Int64) - ID of the cloud credential.
- `direction` (String) - Direction of the cloud sync operation.  * `PUSH`: Upload local files to cloud storage * `PULL`: Download files from cloud storage to local storage Valid values: `PUSH`, `PULL`
- `path` (String) - The local path to back up beginning with `/mnt` or `/dev/zvol`.
- `transfer_mode` (String) - How files are transferred between local and cloud storage.  * `SYNC`: Synchronize directories (add new, update changed, remove deleted) * `COPY`: Copy files without removing any existing files * `MOVE Valid values: `SYNC`, `COPY`, `MOVE`

### Optional

- `args` (String) - (Slated for removal). Default: ``
- `bwlimit` (List) - Schedule of bandwidth limits.
- `create_empty_src_dirs` (Bool) - Whether to create empty directories in the destination that exist in the source. Default: `False`
- `description` (String) - The name of the task to display in the UI. Default: ``
- `enabled` (Bool) - Can enable/disable the task. Default: `True`
- `encryption` (Bool) - Whether to encrypt files before uploading to cloud storage. Default: `False`
- `encryption_password` (String) - Password for client-side encryption. Empty string if encryption is disabled. Default: ``
- `encryption_salt` (String) - Salt value for encryption key derivation. Empty string if encryption is disabled. Default: ``
- `exclude` (List) - Paths to pass to `restic backup --exclude`.
- `filename_encryption` (Bool) - Whether to encrypt filenames in addition to file contents. Default: `False`
- `follow_symlinks` (Bool) - Whether to follow symbolic links and sync the files they point to. Default: `False`
- `include` (List) - Paths to pass to `restic backup --include`.
- `post_script` (String) - A Bash script to run immediately after every backup if it succeeds. Default: ``
- `pre_script` (String) - A Bash script to run immediately before every backup. Default: ``
- `schedule` (String) - Cron schedule dictating when the task should run.
- `snapshot` (Bool) - Whether to create a temporary snapshot of the dataset before every backup. Default: `False`
- `transfers` (Int64) - Maximum number of parallel file transfers. `null` for default. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cloudsync.example <id>
```
