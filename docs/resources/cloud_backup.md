---
page_title: "truenas_cloud_backup Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new cloud backup task
---

# truenas_cloud_backup (Resource)

Create a new cloud backup task


## Example Usage

```terraform
resource "truenas_cloud_backup" "example" {
  attributes = "example-value"
  credentials = 1
  keep_last = 1
  password = "example-value"
  path = "example-value"
}
```

## Schema

### Required

- `attributes` (String) - Additional information for each backup, e.g. bucket name.
- `credentials` (Int64) - ID of the cloud credential to use for each backup.
- `keep_last` (Int64) - How many of the most recent backup snapshots to keep after each backup.
- `password` (String) - Password for the remote repository.
- `path` (String) - The local path to back up beginning with `/mnt` or `/dev/zvol`.

### Optional

- `absolute_paths` (Bool) - Preserve absolute paths in each backup (cannot be set when `snapshot=True`). Default: `False`
- `args` (String) - (Slated for removal). Default: ``
- `cache_path` (String) - Cache path. If not set, performance may degrade. Default: `None`
- `description` (String) - The name of the task to display in the UI. Default: ``
- `enabled` (Bool) - Can enable/disable the task. Default: `True`
- `exclude` (List) - Paths to pass to `restic backup --exclude`.
- `include` (List) - Paths to pass to `restic backup --include`.
- `post_script` (String) - A Bash script to run immediately after every backup if it succeeds. Default: ``
- `pre_script` (String) - A Bash script to run immediately before every backup. Default: ``
- `rate_limit` (Int64) - Maximum upload/download rate in KiB/s. Passed to `restic --limit-upload` on `cloud_backup.sync` and     `restic --limit-download` on `cloud_backup.restore`. `null` indicates no rate limit will be impo Default: `None`
- `schedule` (String) - Cron schedule dictating when the task should run.
- `snapshot` (Bool) - Whether to create a temporary snapshot of the dataset before every backup. Default: `False`
- `transfer_setting` (String) - * DEFAULT:     * pack size given by `$RESTIC_PACK_SIZE` (default 16 MiB)     * read concurrency given by `$RESTIC_READ_CONCURRENCY` (default 2 files)  * PERFORMANCE:     * pack size = 29 MiB     * rea Default: `DEFAULT` Valid values: `DEFAULT`, `PERFORMANCE`, `FAST_STORAGE`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cloud_backup.example <id>
```
