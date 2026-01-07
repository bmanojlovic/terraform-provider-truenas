---
page_title: "truenas_iscsi_extent Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  iSCSI extent configuration data for creation.
---

# truenas_iscsi_extent (Resource)

iSCSI extent configuration data for creation.

## Example Usage

```terraform
resource "truenas_iscsi_extent" "example" {
  name = "example-name"
  type = "DISK"
  filesize = "0"
  blocksize = 512
  pblocksize = false
  comment = ""
  insecure_tpc = true
  xen = false
}
```

## Schema

### Required

- `name` (Required) - Name of the iSCSI extent.. Type: `string`

### Optional

- `type` (Optional) - Type of the extent storage backend. Valid values: `DISK`, `FILE` Default: `DISK`. Type: `string`
- `disk` (Optional) - Disk device to use for the extent or `null` if using a file.. Type: `string`
- `serial` (Optional) - Serial number for the extent or `null` to auto-generate.. Type: `string`
- `path` (Optional) - File path for file-based extents or `null` if using a disk.. Type: `string`
- `filesize` (Optional) - Size of the file-based extent in bytes. Default: `0`. Type: `string`
- `blocksize` (Optional) - Block size for the extent in bytes. Valid values: `512`, `1024`, `2048` Default: `512`. Type: `integer`
- `pblocksize` (Optional) - Whether to use physical block size reporting. Default: `False`. Type: `boolean`
- `avail_threshold` (Optional) - Available space threshold percentage or `null` to disable.. Type: `string`
- `comment` (Optional) - Optional comment describing the extent. Default: ``. Type: `string`
- `insecure_tpc` (Optional) - Whether to enable insecure Third Party Copy (TPC) operations. Default: `True`. Type: `boolean`
- `xen` (Optional) - Whether to enable Xen compatibility mode. Default: `False`. Type: `boolean`
- `rpm` (Optional) - Reported RPM type for the extent. Valid values: `UNKNOWN`, `SSD`, `5400` Default: `SSD`. Type: `string`
- `ro` (Optional) - Whether the extent is read-only. Default: `False`. Type: `boolean`
- `enabled` (Optional) - Whether the extent is enabled and available for use. Default: `True`. Type: `boolean`
- `product_id` (Optional) - Product ID string for the extent or `null` for default.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_extent.example <id>
```
