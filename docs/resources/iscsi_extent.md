---
page_title: "truenas_iscsi_extent Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an iSCSI Extent.
---

# truenas_iscsi_extent (Resource)

Create an iSCSI Extent.


## Example Usage

```terraform
resource "truenas_iscsi_extent" "example" {
  name = "example-value"
}
```

## Schema

### Required

- `name` (String) - Name of the iSCSI extent.

### Optional

- `avail_threshold` (Int64) - Available space threshold percentage or `null` to disable. Default: `None`
- `blocksize` (Int64) - Block size for the extent in bytes. Default: `512` Valid values: `512`, `1024`, `2048`, `4096`
- `comment` (String) - Optional comment describing the extent. Default: ``
- `disk` (String) - Disk device to use for the extent or `null` if using a file. Default: `None`
- `enabled` (Bool) - Whether the extent is enabled and available for use. Default: `True`
- `filesize` (Int64) - Size of the file-based extent in bytes. Default: `0`
- `insecure_tpc` (Bool) - Whether to enable insecure Third Party Copy (TPC) operations. Default: `True`
- `path` (String) - File path for file-based extents or `null` if using a disk. Default: `None`
- `pblocksize` (Bool) - Whether to use physical block size reporting. Default: `False`
- `product_id` (String) - Product ID string for the extent or `null` for default. Default: `None`
- `ro` (Bool) - Whether the extent is read-only. Default: `False`
- `rpm` (String) - Reported RPM type for the extent. Default: `SSD` Valid values: `UNKNOWN`, `SSD`, `5400`, `7200`, `10000`, `15000`
- `serial` (String) - Serial number for the extent or `null` to auto-generate. Default: `None`
- `type` (String) - Type of the extent storage backend. Default: `DISK` Valid values: `DISK`, `FILE`
- `xen` (Bool) - Whether to enable Xen compatibility mode. Default: `False`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_extent.example <id>
```
