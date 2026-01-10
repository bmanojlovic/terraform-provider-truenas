---
page_title: "truenas_disk Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_disk (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_disk" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the disk to retrieve.

### Read-Only

- `advpowermgmt` (String) - Advanced power management level or `DISABLED` to turn off power management.
- `bus` (String) - System bus type the disk is connected to.
- `description` (String) - Human-readable description of the disk device.
- `devname` (String) - Device name in the operating system.
- `enclosure` (String) - Physical enclosure information or `null` if not in an enclosure.
- `expiretime` (String) - Expiration timestamp for disk data or `null` if not applicable.
- `hddstandby` (String) - Hard disk standby timer in minutes or `ALWAYS ON` to disable standby.
- `identifier` (String) - Unique identifier for the disk device.
- `kmip_uid` (String) - KMIP (Key Management Interoperability Protocol) unique identifier or `null`.
- `lunid` (String) - Logical unit number identifier or `null` if not applicable.
- `model` (String) - Manufacturer model name/number of the disk. `null` if not available.
- `name` (String) - System name of the disk device.
- `number` (Int64) - Numeric identifier assigned to the disk.
- `passwd` (String) - Disk encryption password (masked for security).
- `pool` (String) - Name of the storage pool this disk belongs to. `null` if not part of any pool.
- `rotationrate` (Int64) - Disk rotation speed in RPM or `null` for SSDs and unknown devices.
- `serial` (String) - Manufacturer serial number of the disk.
- `size` (Int64) - Total size of the disk in bytes. `null` if not available.
- `subsystem` (String) - Storage subsystem type.
- `transfermode` (String) - Data transfer mode and capabilities of the disk.
- `type` (String) - Disk type classification or `null` if not determined.
- `zfs_guid` (String) - ZFS globally unique identifier for this disk or `null` if not used in ZFS.
