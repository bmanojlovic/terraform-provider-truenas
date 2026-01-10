---
page_title: "truenas_sharing_smb Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Manages TrueNAS sharing.smb resources
---

# truenas_sharing_smb (Resource)

Manages TrueNAS sharing.smb resources

## Example Usage

```terraform
resource "truenas_sharing_smb" "example" {
  name = "example-value"
  path = "example-value"
}
```

## Schema

### Required

- `name` (String) - SMB share name. SMB share names are case-insensitive and must be unique, and are subject     to the following restrictions:  * A share name must be no more than 80 characters in length.  * The followi
- `path` (String) - Local server path to share by using the SMB protocol. The path must start with `/mnt/` and must be in a     ZFS pool.  Use the string `EXTERNAL` if the share works as a DFS proxy.  WARNING: The TrueNA

### Optional

- `access_based_share_enumeration` (Bool) - If set, the share is only included when an SMB client requests a list of shares on the SMB server if     the share (not filesystem) access control list (see `sharing.smb.getacl`) grants access to the  Default: `False`
- `audit` (String) - Audit configuration for monitoring SMB share access and operations.
- `browsable` (Bool) - If set, the share is included when an SMB client requests a list of SMB shares on the TrueNAS server. Default: `True`
- `comment` (String) - Text field that is seen next to a share when an SMB client requests a list of SMB shares on the TrueNAS     server. Default: ``
- `enabled` (Bool) - If unset, the SMB share is not available over the SMB protocol. Default: `True`
- `options` (String) - Additional configuration related to the configured SMB share purpose. If null, then the default     options related to the share purpose will be applied. Default: `None`
- `purpose` (String) - This parameter sets the purpose of the SMB share. It controls how the SMB share behaves and what features are     available through options. The DEFAULT_SHARE setting is best for most applications, an Default: `DEFAULT_SHARE` Valid values: `DEFAULT_SHARE`, `LEGACY_SHARE`, `TIMEMACHINE_SHARE`, `MULTIPROTOCOL_SHARE`, `TIME_LOCKED_SHARE`
- `readonly` (Bool) - If set, SMB clients cannot create or change files and directories in the SMB share.  NOTE: If set, the share path is still writeable by local processes or other file sharing protocols. Default: `False`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_sharing_smb.example <id>
```
