---
page_title: "truenas_cloud_backup Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for the new cloud backup task.
---

# truenas_cloud_backup (Resource)

Configuration for the new cloud backup task.

## Example Usage

```terraform
resource "truenas_cloud_backup" "example" {
  path = "example-path"
  credentials = 1
  attributes = {}
  password = "example-password"
  keep_last = 1
  description = ""
  pre_script = ""
  post_script = ""
}
```

## Schema

### Required

- `path` (Required) - The local path to back up beginning with `/mnt` or `/dev/zvol`.. Type: `string`
- `credentials` (Required) - ID of the cloud credential to use for each backup.. Type: `integer`
- `attributes` (Required) - Additional information for each backup, e.g. bucket name.. Type: `object`
- `password` (Required) - Password for the remote repository.. Type: `string`
- `keep_last` (Required) - How many of the most recent backup snapshots to keep after each backup.. Type: `integer`

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
- `transfer_setting` (Optional) - * DEFAULT:
    * pack size given by `$RESTIC_PACK_SIZE` (default 16 MiB)
    * read concurrency given by `$RESTIC_READ_CONCURRENCY` (default 2 files)

* PERFORMANCE:
    * pack size = 29 MiB
    * read concurrency given by `$RESTIC_READ_CONCURRENCY` (default 2 files)

* FAST_STORAGE:
    * pack size = 58 MiB
    * read concurrency = 100 files Valid values: `DEFAULT`, `PERFORMANCE`, `FAST_STORAGE` Default: `DEFAULT`. Type: `string`
- `absolute_paths` (Optional) - Preserve absolute paths in each backup (cannot be set when `snapshot=True`). Default: `False`. Type: `boolean`
- `cache_path` (Optional) - Cache path. If not set, performance may degrade.. Type: `string`
- `rate_limit` (Optional) - Maximum upload/download rate in KiB/s. Passed to `restic --limit-upload` on `cloud_backup.sync` and     `restic --limit-download` on `cloud_backup.restore`. `null` indicates no rate limit will be imposed.

Can be overridden on a sync or restore call.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cloud_backup.example <id>
```
