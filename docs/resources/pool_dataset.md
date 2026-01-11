---
page_title: "truenas_pool_dataset Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Creates a dataset/zvol.
---

# truenas_pool_dataset (Resource)

Creates a dataset/zvol.


## Variants

This resource has **2 variants** controlled by the `type` field. Choose the appropriate variant for your use case:

### FILESYSTEM

```terraform
resource "truenas_pool_dataset" "example" {
  type = "FILESYSTEM"
  name = "value"
}
```

**Required fields:** `name`

**Key optional fields:** `aclmode`, `acltype`, `atime`, `casesensitivity`, `quota`, `recordsize`, `refquota`

### VOLUME

```terraform
resource "truenas_pool_dataset" "example" {
  type = "VOLUME"
  name = "value"
  volsize = 1024
}
```

**Required fields:** `name`, `volsize`

**Key optional fields:** `force_size`, `sparse`, `volblocksize`



## Schema

### Required

- `name` (String) - The name of the dataset to create.

### Optional

- `aclmode` (String) - How Access Control Lists are handled when chmod is used. Valid values: `PASSTHROUGH`, `RESTRICTED`, `DISCARD`, `INHERIT` **Applies to:** `FILESYSTEM`
- `acltype` (String) - The type of Access Control List system to use. Valid values: `OFF`, `NFSV4`, `POSIX`, `INHERIT` **Applies to:** `FILESYSTEM`
- `atime` (String) - Whether file access times are updated when files are accessed. Valid values: `ON`, `OFF`, `INHERIT` **Applies to:** `FILESYSTEM`
- `casesensitivity` (String) - File name case sensitivity setting. Valid values: `SENSITIVE`, `INSENSITIVE`, `INHERIT` **Applies to:** `FILESYSTEM`
- `checksum` (String) - Checksum algorithm to verify data integrity. Higher security algorithms like SHA256 provide better     protection but use more CPU. Default: `INHERIT` Valid values: `ON`, `OFF`, `FLETCHER2`, `FLETCHER4`, `SHA256`, `SHA512`, `SKEIN`, `EDONR`, `BLAKE3`, `INHERIT`
- `comments` (String) - Comments or description for the dataset. Default: `INHERIT`
- `compression` (String) - Compression algorithm to use for the dataset. Higher numbered variants provide better compression     but use more CPU. 'INHERIT' uses the parent dataset's setting. Default: `INHERIT` Valid values: `ON`, `OFF`, `LZ4`, `GZIP`, `GZIP-1`, `GZIP-9`, `ZSTD`, `ZSTD-FAST`, `ZLE`, `LZJB`
- `copies` (Int64) - Number of copies of data blocks to maintain for redundancy. Default: `INHERIT`
- `create_ancestors` (Bool) - Whether to create any missing parent datasets. Default: `False`
- `deduplication` (String) - Deduplication setting. 'ON' enables dedup, 'VERIFY' enables with checksum verification, 'OFF' disables. Default: `INHERIT` Valid values: `ON`, `VERIFY`, `OFF`, `INHERIT`
- `encryption` (Bool) - Create a ZFS encrypted root dataset for `name` pool. There is 1 case where ZFS encryption is not allowed for a dataset: 1) If the parent dataset is encrypted with a passphrase and `name` is being crea Default: `False`
- `encryption_options` (String) - Configuration for encryption of dataset for `name` pool.
- `exec` (String) - Whether files in this dataset can be executed. Default: `INHERIT` Valid values: `ON`, `OFF`, `INHERIT`
- `force_size` (Bool) - Force creation even if the size is not optimal. **Applies to:** `VOLUME`
- `inherit_encryption` (Bool) - Whether to inherit encryption settings from the parent dataset. Default: `True`
- `managedby` (String) - Identifies which service or system manages this dataset. Default: `INHERIT`
- `quota` (Int64) - Maximum disk space this dataset and its children can consume in bytes. **Applies to:** `FILESYSTEM`
- `quota_critical` (Int64) - Percentage of dataset quota at which to issue a critical alert. 0-100 or 'INHERIT'. Default: `INHERIT`
- `quota_warning` (Int64) - Percentage of dataset quota at which to issue a warning. 0-100 or 'INHERIT'. Default: `INHERIT`
- `readonly` (String) - Whether the dataset is read-only. Default: `INHERIT` Valid values: `ON`, `OFF`, `INHERIT`
- `recordsize` (String) - The suggested block size for files in this filesystem dataset. **Applies to:** `FILESYSTEM`
- `refquota` (Int64) - Maximum disk space this dataset itself can consume in bytes. **Applies to:** `FILESYSTEM`
- `refquota_critical` (Int64) - Percentage of reference quota at which to issue a critical alert. 0-100 or 'INHERIT'. Default: `INHERIT`
- `refquota_warning` (Int64) - Percentage of reference quota at which to issue a warning. 0-100 or 'INHERIT'. Default: `INHERIT`
- `refreservation` (Int64) - Minimum disk space guaranteed to this dataset itself in bytes.
- `reservation` (Int64) - Minimum disk space guaranteed to this dataset and its children in bytes.
- `share_type` (String) - Optimization type for the dataset based on its intended use. Default: `GENERIC` Valid values: `GENERIC`, `MULTIPROTOCOL`, `NFS`, `SMB`, `APPS`
- `snapdev` (String) - Controls visibility of volume snapshots under /dev/zvol/. Valid values: `HIDDEN`, `VISIBLE`, `INHERIT`
- `snapdir` (String) - Controls visibility of the `.zfs/snapshot` directory. 'DISABLED' hides snapshots, 'VISIBLE' shows them,     'HIDDEN' makes them accessible but not listed. Default: `INHERIT` Valid values: `DISABLED`, `VISIBLE`, `HIDDEN`, `INHERIT`
- `sparse` (Bool) - Whether to use sparse (thin) provisioning for the volume. **Applies to:** `VOLUME`
- `special_small_block_size` (Int64) - Size threshold below which blocks are stored on special vdevs.
- `sync` (String) - Synchronous write behavior for the dataset. Default: `INHERIT` Valid values: `STANDARD`, `ALWAYS`, `DISABLED`, `INHERIT`
- `type` (String) - Type of dataset to create - volume (zvol). Default: `VOLUME` Valid values: `FILESYSTEM`, `VOLUME`
- `user_properties` (List) - Custom user-defined properties to set on the dataset. Default: `[]`
- `volblocksize` (String) - Defaults to `128K` if the parent pool is a DRAID pool or `16K` otherwise. Valid values: `512`, `512B`, `1K`, `2K`, `4K`, `8K`, `16K`, `32K`, `64K`, `128K` **Applies to:** `VOLUME`
- `volsize` (Int64) - The volume size in bytes; supposed to be a multiple of the block size. **Applies to:** `VOLUME`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_dataset.example <id>
```
