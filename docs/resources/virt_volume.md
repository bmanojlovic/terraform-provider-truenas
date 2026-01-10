---
page_title: "truenas_virt_volume Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Manages TrueNAS virt.volume resources
---

# truenas_virt_volume (Resource)

Manages TrueNAS virt.volume resources

## Example Usage

```terraform
resource "truenas_virt_volume" "example" {
  name = "example-value"
}
```

## Schema

### Required

- `name` (String) - Name for the new virtualization volume (alphanumeric, dashes, dots, underscores).

### Optional

- `content_type` (String) -  Default: `BLOCK` Valid values: `BLOCK`
- `size` (Int64) - Size of volume in MB and it should at least be 512 MB. Default: `1024`
- `storage_pool` (String) - Storage pool in which to create the volume. This must be one of pools listed     in virt.global.config output under `storage_pools`. If the value is None, then     the pool defined as `pool` in virt.g Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_virt_volume.example <id>
```
