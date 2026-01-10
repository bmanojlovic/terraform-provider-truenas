---
page_title: "truenas_pool Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_pool (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_pool" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool to retrieve.

### Read-Only

- `allocated` (Int64) - Amount of space currently allocated in the pool in bytes. `null` if not available.
- `allocated_str` (String) - Human-readable string representation of allocated space. `null` if not available.
- `autotrim` (String) - Auto-trim configuration for the pool indicating whether automatic TRIM operations are enabled.
- `dedup_table_quota` (String) - Quota limit for the deduplication table. `null` if no quota is set.
- `dedup_table_size` (Int64) - Size of the deduplication table in bytes. `null` if deduplication is not enabled.
- `expand` (String) - Information about any active pool expansion operation. `null` if no expansion is running.
- `fragmentation` (String) - Percentage of pool fragmentation as a string. `null` if not available.
- `free` (Int64) - Amount of free space available in the pool in bytes. `null` if not available.
- `free_str` (String) - Human-readable string representation of free space. `null` if not available.
- `freeing` (Int64) - Amount of space being freed (in bytes) by ongoing operations. `null` if not available.
- `freeing_str` (String) - Human-readable string representation of space being freed. `null` if not available.
- `guid` (String) - Globally unique identifier (GUID) for this pool.
- `healthy` (Bool) - Whether the pool is in a healthy state with no errors or warnings.
- `is_upgraded` (Bool) - Whether this pool has been upgraded to the latest feature flags.
- `name` (String) - Name of the storage pool.
- `path` (String) - Filesystem path where the pool is mounted.
- `scan` (String) - Information about any active scrub or resilver operation. `null` if no operation is running.
- `size` (Int64) - Total size of the pool in bytes. `null` if not available.
- `size_str` (String) - Human-readable string representation of the pool size. `null` if not available.
- `status` (String) - Current status of the pool.
- `status_code` (String) - Detailed status code for the pool condition. `null` if not applicable.
- `status_detail` (String) - Human-readable description of the pool status. `null` if not available.
- `topology` (String) - Physical topology and structure of the pool including vdevs. `null` if not available.
- `warning` (Bool) - Whether the pool has warning conditions that require attention.
