---
page_title: "truenas_sharing_nfs Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  NFS share configuration data for the new share.
---

# truenas_sharing_nfs (Resource)

NFS share configuration data for the new share.

## Example Usage

```terraform
resource "truenas_sharing_nfs" "example" {
  path = "example-path"
  aliases = []
  comment = ""
  networks = []
  hosts = []
  ro = false
  security = []
  enabled = true
}
```

## Schema

### Required

- `path` (Required) - Local path to be exported. . Type: `string`

### Optional

- `aliases` (Optional) - IGNORED for now.  Default: `[]`. Type: `array`
- `comment` (Optional) - User comment associated with share.  Default: ``. Type: `string`
- `networks` (Optional) - List of authorized networks that are allowed to access the share having format     "network/mask" CIDR notation. Each entry must be unique. If empty, all networks are allowed.
Excessively long lists should be avoided. Default: `[]`. Type: `array`
- `hosts` (Optional) - List of IP's/hostnames which are allowed to access the share. No quotes or spaces are allowed.
Each entry must be unique. If empty, all IP's/hostnames are allowed.
Excessively long lists should be avoided. Default: `[]`. Type: `array`
- `ro` (Optional) - Export the share as read only.  Default: `False`. Type: `boolean`
- `maproot_user` (Optional) - Map root user client to a specified user. . Type: `string`
- `maproot_group` (Optional) - Map root group client to a specified group. . Type: `string`
- `mapall_user` (Optional) - Map all client users to a specified user. . Type: `string`
- `mapall_group` (Optional) - Map all client groups to a specified group. . Type: `string`
- `security` (Optional) - Specify the security schema.  Default: `[]`. Type: `array`
- `enabled` (Optional) - Enable or disable the share.  Default: `True`. Type: `boolean`
- `expose_snapshots` (Optional) - Enterprise feature to enable access to the ZFS snapshot directory for the export.
Export path must be the root directory of a ZFS dataset. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_sharing_nfs.example <id>
```
