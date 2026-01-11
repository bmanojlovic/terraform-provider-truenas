---
page_title: "truenas_pool Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new ZFS Pool.
---

# truenas_pool (Resource)

Create a new ZFS Pool.


## Example Usage

```terraform
resource "truenas_pool" "example" {
  name = "example-value"
  topology = "example-value"
}
```

## Schema

### Required

- `name` (String) - Name for the new storage pool.
- `topology` (String) - Physical layout and configuration of vdevs in the pool.

### Optional

- `allow_duplicate_serials` (Bool) - Whether to allow disks with duplicate serial numbers in the pool. Default: `False`
- `checksum` (String) - Checksum algorithm to use for data integrity verification. Default: `None` Valid values: `ON`, `OFF`, `FLETCHER2`, `FLETCHER4`, `SHA256`, `SHA512`, `SKEIN`, `EDONR`, `BLAKE3`, `None`
- `dedup_table_quota` (String) - How to manage the deduplication table quota allocation. Default: `AUTO` Valid values: `AUTO`, `CUSTOM`, `None`
- `dedup_table_quota_value` (Int64) - Custom quota value in bytes when `dedup_table_quota` is set to CUSTOM. Default: `None`
- `deduplication` (String) - Make sure no block of data is duplicated in the pool. If set to `VERIFY` and two blocks have similar     signatures, byte-to-byte comparison is performed to ensure that the blcoks are identical. This  Default: `None` Valid values: `ON`, `VERIFY`, `OFF`, `None`
- `encryption` (Bool) - If set, create a ZFS encrypted root dataset for this pool. Default: `False`
- `encryption_options` (String) - Specify configuration for encryption of root dataset.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool.example <id>
```
