---
page_title: "truenas_pool_dataset Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_pool_dataset (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_pool_dataset" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_dataset to retrieve.

### Read-Only

- `aclmode` (String) - How Access Control Lists (ACLs) are handled when chmod is used.
- `acltype` (String) - The type of Access Control List system used (NFSV4, POSIX, or OFF).
- `atime` (String) - Whether file access times are updated when files are accessed.
- `available` (String) - Amount of disk space available to this dataset and its children.
- `casesensitivity` (String) - File name case sensitivity setting (sensitive/insensitive).
- `checksum` (String) - Data integrity checksum algorithm used for this dataset.
- `children` (List) - Array of child dataset objects nested under this dataset.
- `comments` (String) - ZFS comments property for storing descriptive text about the dataset.
- `compression` (String) - Compression algorithm and level applied to data in this dataset.
- `compressratio` (String) - The achieved compression ratio as a decimal (e.g., '2.50x').
- `copies` (String) - Number of copies of data blocks to maintain for redundancy (1-3).
- `creation` (String) - Timestamp when this dataset was created.
- `deduplication` (String) - ZFS deduplication setting - whether identical data blocks are stored only once.
- `encrypted` (Bool) - Whether the dataset is encrypted.
- `encryption_algorithm` (String) - Encryption algorithm used (e.g., AES-256-GCM). Only relevant for encrypted datasets.
- `encryption_root` (String) - The root dataset where encryption is enabled. `null` if the dataset is not encrypted.
- `exec` (String) - Whether files in this dataset can be executed.
- `key_format` (String) - Format of the encryption key (hex/raw/passphrase). Only relevant for encrypted datasets.
- `key_loaded` (String) - Whether the encryption key is currently loaded for encrypted datasets. `null` for unencrypted datasets.
- `locked` (Bool) - Whether an encrypted dataset is currently locked (key not loaded).
- `managedby` (String) - Identifies which service or system manages this dataset.
- `mountpoint` (String) - Filesystem path where this dataset is mounted. Null for unmounted datasets or volumes.
- `name` (String) - The dataset name without the pool prefix.
- `origin` (String) - The snapshot from which this clone was created. Empty for non-clone datasets.
- `pbkdf2iters` (String) - Number of PBKDF2 iterations used for passphrase-based encryption keys.
- `pool` (String) - The name of the ZFS pool containing this dataset.
- `quota` (String) - Maximum amount of disk space this dataset and its children can consume.
- `quota_critical` (String) - ZFS quota critical threshold property as a percentage.
- `quota_warning` (String) - ZFS quota warning threshold property as a percentage.
- `readonly` (String) - Whether the dataset is read-only.
- `recordsize` (String) - The suggested block size for files in this filesystem dataset.
- `refquota` (String) - Maximum amount of disk space this dataset itself can consume (excluding children).
- `refquota_critical` (String) - ZFS reference quota critical threshold property as a percentage.
- `refquota_warning` (String) - ZFS reference quota warning threshold property as a percentage.
- `refreservation` (String) - Minimum amount of disk space guaranteed to be available to this dataset itself.
- `reservation` (String) - Minimum amount of disk space guaranteed to be available to this dataset and its children.
- `snapdev` (String) - Controls visibility of volume snapshots under /dev/zvol/<pool>/.
- `snapdir` (String) - Visibility of the .zfs/snapshot directory (visible/hidden).
- `sparse` (String) - For volumes, whether to use sparse (thin) provisioning.
- `special_small_block_size` (String) - Size threshold below which blocks are stored on special vdevs if configured.
- `sync` (String) - Synchronous write behavior (standard/always/disabled).
- `type` (String) - The dataset type.
- `used` (String) - Total amount of disk space consumed by this dataset and all its children.
- `usedbychildren` (String) - Amount of disk space consumed by child datasets.
- `usedbydataset` (String) - Amount of disk space consumed by this dataset itself, excluding children and snapshots.
- `usedbyrefreservation` (String) - Amount of disk space consumed by the refreservation of this dataset.
- `usedbysnapshots` (String) - Amount of disk space consumed by snapshots of this dataset.
- `user_properties` (String) - Custom user-defined ZFS properties set on this dataset as key-value pairs.
- `volblocksize` (String) - For volumes, the block size used by the volume.
- `volsize` (String) - For volumes, the logical size of the volume.
- `xattr` (String) - Extended attributes storage method (on/off).
