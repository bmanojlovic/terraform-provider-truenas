---
page_title: "truenas_cronjob Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new cron job.
---

# truenas_cronjob (Resource)

Create a new cron job.


## Example Usage

```terraform
resource "truenas_cronjob" "example" {
  command = "example-value"
  user = "example-value"
}
```

## Schema

### Required

- `command` (String) - Shell command or script to execute.
- `user` (String) - System user account to run the command as.

### Optional

- `description` (String) - Human-readable description of what this cron job does. Default: ``
- `enabled` (Bool) - Whether the cron job is active and will be executed. Default: `True`
- `schedule` (String) - Cron schedule configuration for when the job runs. Default: `{'minute': '00', 'hour': '*', 'dom': '*', 'month': '*', 'dow': '*'}`
- `stderr` (Bool) - Whether to IGNORE standard error (if `false`, it will be added to email). Default: `False`
- `stdout` (Bool) - Whether to IGNORE standard output (if `false`, it will be added to email). Default: `True`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cronjob.example <id>
```
