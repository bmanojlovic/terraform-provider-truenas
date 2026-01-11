---
page_title: "truenas_pool_snapshottask Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a Periodic Snapshot Task
---

# truenas_pool_snapshottask (Resource)

Create a Periodic Snapshot Task


## Example Usage

```terraform
resource "truenas_pool_snapshottask" "example" {
  dataset = "example-value"
}
```

## Schema

### Required

- `dataset` (String) - The dataset to take snapshots of.

### Optional

- `allow_empty` (Bool) - Whether to take snapshots even if no data has changed. Default: `True`
- `enabled` (Bool) - Whether this periodic snapshot task is enabled. Default: `True`
- `exclude` (List) - Array of dataset patterns to exclude from recursive snapshots. Default: `[]`
- `lifetime_unit` (String) - Unit of time for snapshot retention. Default: `WEEK` Valid values: `HOUR`, `DAY`, `WEEK`, `MONTH`, `YEAR`
- `lifetime_value` (Int64) - Number of time units to retain snapshots. `lifetime_unit` gives the time unit. Default: `2`
- `naming_schema` (String) - Naming pattern for generated snapshots using strftime format. Default: `auto-%Y-%m-%d_%H-%M`
- `recursive` (Bool) - Whether to recursively snapshot child datasets. Default: `False`
- `schedule` (String) - Cron schedule for when snapshots should be taken.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_snapshottask.example <id>
```
