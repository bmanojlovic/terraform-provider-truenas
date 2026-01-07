---
page_title: "truenas_virt_volume Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  VirtVolumeCreateArgs parameters.
---

# truenas_virt_volume (Resource)

VirtVolumeCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_virt_volume" "example" {
  name = "example-name"
  content_type = "BLOCK"
  size = 1024
}
```

## Schema

### Required

- `name` (Required) - Name for the new virtualization volume (alphanumeric, dashes, dots, underscores).. Type: `string`

### Optional

- `content_type` (Optional) - content_type configuration Valid values: `BLOCK` Default: `BLOCK`. Type: `string`
- `size` (Optional) - Size of volume in MB and it should at least be 512 MB. Default: `1024`. Type: `integer`
- `storage_pool` (Optional) - Storage pool in which to create the volume. This must be one of pools listed     in virt.global.config output under `storage_pools`. If the value is None, then     the pool defined as `pool` in virt.global.config will be used.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_virt_volume.example <id>
```
