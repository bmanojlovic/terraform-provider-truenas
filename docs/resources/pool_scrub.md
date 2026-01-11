---
page_title: "truenas_pool_scrub Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a scrub task for a pool.
---

# truenas_pool_scrub (Resource)

Create a scrub task for a pool.


## Example Usage

```terraform
resource "truenas_pool_scrub" "example" {
  pool = 1
}
```

## Schema

### Required

- `pool` (Int64) - ID of the pool to scrub.

### Optional

- `description` (String) - Description or notes for this scrub schedule. Default: ``
- `enabled` (Bool) - Whether this scrub schedule is enabled. Default: `True`
- `schedule` (String) - Cron schedule for when scrubs should run.
- `threshold` (Int64) - Days before a scrub is due when a scrub should automatically start. Default: `35`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_scrub.example <id>
```
