---
page_title: "truenas_virt_instance Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  VirtInstanceCreateArgs parameters.
---

# truenas_virt_instance (Resource)

VirtInstanceCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_virt_instance" "example" {
  name = "example-name"
  image = "example-image"
  start_on_create = true
  source_type = "IMAGE"
  root_disk_size = 10
  root_disk_io_bus = "NVME"
  remote = "LINUX_CONTAINERS"
  instance_type = "CONTAINER"
}
```

## Schema

### Required

- `name` (Required) - Name for the new virtual instance.. Type: `string`
- `image` (Required) - Image identifier to use for creating the instance.. Type: `string`

### Optional

- `start_on_create` (Optional) - Start the resource immediately after creation. Default behavior: starts if not specified. Type: `boolean`
- `source_type` (Optional) - Source type for instance creation. Valid values: `IMAGE` Default: `IMAGE`. Type: `string`
- `storage_pool` (Optional) - Storage pool under which to allocate root filesystem. Must be one of the pools     listed in virt.global.config output under "storage_pools". If None (default) then the pool     specified in the global configuration will be used.. Type: `string`
- `root_disk_size` (Optional) - This can be specified when creating VMs so the root device's size can be configured. Root device for VMs     is a sparse zvol and the field specifies space in GBs and defaults to 10 GBs. Default: `10`. Type: `integer`
- `root_disk_io_bus` (Optional) - I/O bus type for the root disk (defaults to NVME for best performance). Valid values: `NVME`, `VIRTIO-BLK`, `VIRTIO-SCSI` Default: `NVME`. Type: `string`
- `remote` (Optional) - Remote image source to use. Valid values: `LINUX_CONTAINERS` Default: `LINUX_CONTAINERS`. Type: `string`
- `instance_type` (Optional) - Type of instance to create. Valid values: `CONTAINER` Default: `CONTAINER`. Type: `string`
- `environment` (Optional) - Environment variables to set inside the instance.. Type: `string`
- `autostart` (Optional) - Whether the instance should automatically start when the host boots. Default: `True`. Type: `string`
- `cpu` (Optional) - CPU allocation specification or `null` for automatic allocation.. Type: `string`
- `devices` (Optional) - Array of devices to attach to the instance.. Type: `string`
- `memory` (Optional) - Memory allocation in bytes or `null` for automatic allocation.. Type: `string`
- `privileged_mode` (Optional) - This is only valid for containers and should only be set when container instance which is to be deployed is to     run in a privileged mode. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_virt_instance.example <id>
```
