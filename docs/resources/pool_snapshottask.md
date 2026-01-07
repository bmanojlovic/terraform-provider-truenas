---
page_title: "truenas_pool_snapshottask Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for the new periodic snapshot task.
---

# truenas_pool_snapshottask (Resource)

Configuration for the new periodic snapshot task.

## Example Usage

```terraform
resource "truenas_pool_snapshottask" "example" {
  dataset = "example-dataset"
  recursive = false
  lifetime_value = 2
  lifetime_unit = "HOUR"
  enabled = true
  exclude = []
  naming_schema = "auto-%Y-%m-%d_%H-%M"
  allow_empty = true
}
```

## Schema

### Required

- `dataset` (Required) - The dataset to take snapshots of.. Type: `string`

### Optional

- `recursive` (Optional) - Whether to recursively snapshot child datasets. Default: `False`. Type: `boolean`
- `lifetime_value` (Optional) - Number of time units to retain snapshots. `lifetime_unit` gives the time unit. Default: `2`. Type: `integer`
- `lifetime_unit` (Optional) - Unit of time for snapshot retention. Valid values: `HOUR`, `DAY`, `WEEK` Default: `WEEK`. Type: `string`
- `enabled` (Optional) - Whether this periodic snapshot task is enabled. Default: `True`. Type: `boolean`
- `exclude` (Optional) - Array of dataset patterns to exclude from recursive snapshots. Default: `[]`. Type: `array`
- `naming_schema` (Optional) - Naming pattern for generated snapshots using strftime format. Default: `auto-%Y-%m-%d_%H-%M`. Type: `string`
- `allow_empty` (Optional) - Whether to take snapshots even if no data has changed. Default: `True`. Type: `boolean`
- `schedule` (Optional) - Cron schedule for when snapshots should be taken.. Type: `object`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_snapshottask.example <id>
```
