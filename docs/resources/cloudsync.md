---
page_title: "truenas_cloudsync Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Cloud sync task configuration data.
---

# truenas_cloudsync (Resource)

Cloud sync task configuration data.

## Example Usage

```terraform
resource "truenas_cloudsync" "example" {
  path = "example-path"
  credentials = 1
  attributes = {}
  direction = "PUSH"
  transfer_mode = "SYNC"
  description = ""
  pre_script = ""
  post_script = ""
}
```

## Schema

### Required

- `path` (Required) - The local path to back up beginning with `/mnt` or `/dev/zvol`.. Type: `string`
- `credentials` (Required) - ID of the cloud credential.. Type: `integer`
- `attributes` (Required) - Additional information for each backup, e.g. bucket name.. Type: `object`
- `direction` (Required) - Direction of the cloud sync operation.

* `PUSH`: Upload local files to cloud storage
* `PULL`: Download files from cloud storage to local storage Valid values: `PUSH`, `PULL`. Type: `string`
- `transfer_mode` (Required) - How files are transferred between local and cloud storage.

* `SYNC`: Synchronize directories (add new, update changed, remove deleted)
* `COPY`: Copy files without removing any existing files
* `MOVE`: Move files (copy then delete from source) Valid values: `SYNC`, `COPY`, `MOVE`. Type: `string`

### Optional

- `description` (Optional) - The name of the task to display in the UI. Default: ``. Type: `string`
- `schedule` (Optional) - Cron schedule dictating when the task should run.. Type: `object`
- `pre_script` (Optional) - A Bash script to run immediately before every backup. Default: ``. Type: `string`
- `post_script` (Optional) - A Bash script to run immediately after every backup if it succeeds. Default: ``. Type: `string`
- `snapshot` (Optional) - Whether to create a temporary snapshot of the dataset before every backup. Default: `False`. Type: `boolean`
- `include` (Optional) - Paths to pass to `restic backup --include`.. Type: `array`
- `exclude` (Optional) - Paths to pass to `restic backup --exclude`.. Type: `array`
- `args` (Optional) - (Slated for removal). Default: ``. Type: `string`
- `enabled` (Optional) - Can enable/disable the task. Default: `True`. Type: `boolean`
- `bwlimit` (Optional) - Schedule of bandwidth limits.. Type: `array`
- `transfers` (Optional) - Maximum number of parallel file transfers. `null` for default.. Type: `string`
- `encryption` (Optional) - Whether to encrypt files before uploading to cloud storage. Default: `False`. Type: `boolean`
- `filename_encryption` (Optional) - Whether to encrypt filenames in addition to file contents. Default: `False`. Type: `boolean`
- `encryption_password` (Optional) - Password for client-side encryption. Empty string if encryption is disabled. Default: ``. Type: `string`
- `encryption_salt` (Optional) - Salt value for encryption key derivation. Empty string if encryption is disabled. Default: ``. Type: `string`
- `create_empty_src_dirs` (Optional) - Whether to create empty directories in the destination that exist in the source. Default: `False`. Type: `boolean`
- `follow_symlinks` (Optional) - Whether to follow symbolic links and sync the files they point to. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cloudsync.example <id>
```
