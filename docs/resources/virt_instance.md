---
page_title: "truenas_virt_instance Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new virtualized instance.
---

# truenas_virt_instance (Resource)

Create a new virtualized instance.


## Example Usage

```terraform
resource "truenas_virt_instance" "example" {
  image = "example-value"
  name = "example-value"
  start_on_create = true
}
```

## Schema

### Required

- `image` (String) - Image identifier to use for creating the instance.
- `name` (String) - Name for the new virtual instance.

### Optional

- `start_on_create` (Bool) - Start the resource immediately after creation. Default: `true`
- `autostart` (String) - Whether the instance should automatically start when the host boots. Default: `True`
- `cpu` (String) - CPU allocation specification or `null` for automatic allocation. Default: `None`
- `devices` (String) - Array of devices to attach to the instance. Default: `None`
- `environment` (String) - Environment variables to set inside the instance. Default: `None`
- `instance_type` (String) - Type of instance to create. Default: `CONTAINER` Valid values: `CONTAINER`
- `memory` (Int64) - Memory allocation in bytes or `null` for automatic allocation. Default: `None`
- `privileged_mode` (Bool) - This is only valid for containers and should only be set when container instance which is to be deployed is to     run in a privileged mode. Default: `False`
- `remote` (String) - Remote image source to use. Default: `LINUX_CONTAINERS` Valid values: `LINUX_CONTAINERS`
- `root_disk_io_bus` (String) - I/O bus type for the root disk (defaults to NVME for best performance). Default: `NVME` Valid values: `NVME`, `VIRTIO-BLK`, `VIRTIO-SCSI`
- `root_disk_size` (Int64) - This can be specified when creating VMs so the root device's size can be configured. Root device for VMs     is a sparse zvol and the field specifies space in GBs and defaults to 10 GBs. Default: `10`
- `source_type` (String) - Source type for instance creation. Default: `IMAGE` Valid values: `IMAGE`
- `storage_pool` (String) - Storage pool under which to allocate root filesystem. Must be one of the pools     listed in virt.global.config output under "storage_pools". If None (default) then the pool     specified in the globa Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_virt_instance.example <id>
```
