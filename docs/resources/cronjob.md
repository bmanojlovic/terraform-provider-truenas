---
page_title: "truenas_cronjob Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Cron job configuration data for the new job.
---

# truenas_cronjob (Resource)

Cron job configuration data for the new job.

## Example Usage

```terraform
resource "truenas_cronjob" "example" {
  command = "example-command"
  user = "example-user"
  enabled = true
  stderr = false
  stdout = true
  schedule = {'minute': '00', 'hour': '*', 'dom': '*', 'month': '*', 'dow': '*'}
  description = ""
}
```

## Schema

### Required

- `command` (Required) - Shell command or script to execute.. Type: `string`
- `user` (Required) - System user account to run the command as.. Type: `string`

### Optional

- `enabled` (Optional) - Whether the cron job is active and will be executed. Default: `True`. Type: `boolean`
- `stderr` (Optional) - Whether to IGNORE standard error (if `false`, it will be added to email). Default: `False`. Type: `boolean`
- `stdout` (Optional) - Whether to IGNORE standard output (if `false`, it will be added to email). Default: `True`. Type: `boolean`
- `schedule` (Optional) - Cron schedule configuration for when the job runs. Default: `{'minute': '00', 'hour': '*', 'dom': '*', 'month': '*', 'dow': '*'}`. Type: `object`
- `description` (Optional) - Human-readable description of what this cron job does. Default: ``. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cronjob.example <id>
```
