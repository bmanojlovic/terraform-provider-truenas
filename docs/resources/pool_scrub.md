---
page_title: "truenas_pool_scrub Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for the new scrub schedule.
---

# truenas_pool_scrub (Resource)

Configuration for the new scrub schedule.

## Example Usage

```terraform
resource "truenas_pool_scrub" "example" {
  pool = 1
  threshold = 35
  description = ""
  enabled = true
}
```

## Schema

### Required

- `pool` (Required) - ID of the pool to scrub.. Type: `integer`

### Optional

- `threshold` (Optional) - Days before a scrub is due when a scrub should automatically start. Default: `35`. Type: `integer`
- `description` (Optional) - Description or notes for this scrub schedule. Default: ``. Type: `string`
- `schedule` (Optional) - Cron schedule for when scrubs should run.. Type: `object`
- `enabled` (Optional) - Whether this scrub schedule is enabled. Default: `True`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_scrub.example <id>
```
