---
page_title: "truenas_sharing_nfs Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a NFS Share.
---

# truenas_sharing_nfs (Resource)

Create a NFS Share.

## Example Usage

```terraform
resource "truenas_sharing_nfs" "example" {
  path = "example-value"
}
```

## Schema

### Required

- `path` (String) - Local path to be exported.

### Optional

- `aliases` (List) - IGNORED for now. Default: `[]`
- `comment` (String) - User comment associated with share. Default: ``
- `enabled` (Bool) - Enable or disable the share. Default: `True`
- `expose_snapshots` (Bool) - Enterprise feature to enable access to the ZFS snapshot directory for the export. Export path must be the root directory of a ZFS dataset. Default: `False`
- `hosts` (List) - List of IP's/hostnames which are allowed to access the share. No quotes or spaces are allowed. Each entry must be unique. If empty, all IP's/hostnames are allowed. Excessively long lists should be avo Default: `[]`
- `mapall_group` (String) - Map all client groups to a specified group. Default: `None`
- `mapall_user` (String) - Map all client users to a specified user. Default: `None`
- `maproot_group` (String) - Map root group client to a specified group. Default: `None`
- `maproot_user` (String) - Map root user client to a specified user. Default: `None`
- `networks` (List) - List of authorized networks that are allowed to access the share having format     "network/mask" CIDR notation. Each entry must be unique. If empty, all networks are allowed. Excessively long lists s Default: `[]`
- `ro` (Bool) - Export the share as read only. Default: `False`
- `security` (List) - Specify the security schema. Default: `[]`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_sharing_nfs.example <id>
```
