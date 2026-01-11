---
page_title: "truenas_pool_snapshot Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Take a snapshot from a given dataset.
---

# truenas_pool_snapshot (Resource)

Take a snapshot from a given dataset.


## Example Usage

```terraform
resource "truenas_pool_snapshot" "example" {
  dataset = "example-value"
}
```

## Schema

### Required

- `dataset` (String) - Name of the dataset to create a snapshot of.

### Optional

- `exclude` (List) - Array of dataset patterns to exclude from recursive snapshots. Default: `[]`
- `name` (String) - Explicit name for the snapshot.
- `naming_schema` (String) - Naming schema pattern to generate the snapshot name automatically.
- `properties` (String) - Object mapping ZFS property names to values to set on the snapshot. Default: `{}`
- `recursive` (Bool) - Whether to recursively snapshot child datasets. Default: `False`
- `vmware_sync` (Bool) - Whether to sync VMware VMs before taking the snapshot. Default: `False`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_snapshot.example <id>
```
