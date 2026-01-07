---
page_title: "truenas_sharing_smb Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  SMB share configuration data for the new share.
---

# truenas_sharing_smb (Resource)

SMB share configuration data for the new share.

## Example Usage

```terraform
resource "truenas_sharing_smb" "example" {
  name = "example-name"
  path = "example-path"
  purpose = "DEFAULT_SHARE"
  enabled = true
  comment = ""
  readonly = false
  browsable = true
  access_based_share_enumeration = false
}
```

## Schema

### Required

- `name` (Required) - SMB share name. SMB share names are case-insensitive and must be unique, and are subject     to the following restrictions:

* A share name must be no more than 80 characters in length.

* The following characters are illegal in a share name: `\ / [ ] : | < > + = ; , * ? "`

* Unicode control characters are illegal in a share name.

* The following share names are not allowed: global, printers, homes.. Type: `string`
- `path` (Required) - Local server path to share by using the SMB protocol. The path must start with `/mnt/` and must be in a     ZFS pool.

Use the string `EXTERNAL` if the share works as a DFS proxy.

WARNING: The TrueNAS server does not check if external paths are reachable. . Type: `string`

### Optional

- `purpose` (Optional) - This parameter sets the purpose of the SMB share. It controls how the SMB share behaves and what features are     available through options. The DEFAULT_SHARE setting is best for most applications, and should be used, unless     there is a specific reason to change it.

* `DEFAULT_SHARE`: Set the SMB share for best compatibility with common SMB clients.

* `LEGACY_SHARE`: Set the SMB share for compatibility with older TrueNAS versions. Automated backend migrations       use this to help the administrator move to better-supported share settings. It should not be used for new SMB       shares.

* `TIMEMACHINE_SHARE`: The SMB share is presented to MacOS clients as a time machine target.
  NOTE: `aapl_extensions` must be set in the global `smb.config`.

* `MULTIPROTOCOL_SHARE`: The SMB share is configured for multi-protocol access. Set this if the `path` is shared       through NFS, FTP, or used by containers or apps.
  NOTE: This setting can reduce SMB share performance because it turns off some SMB features for safer       interoperability with external processes.

* `TIME_LOCKED_SHARE`: The SMB share makes files read-only through the SMB protocol after the set grace_period       ends.
  WARNING: This setting does not work if the `path` is accessed locally or if another SMB share without the       `TIME_LOCKED_SHARE` purpose uses the same path.
  WARNING: This setting might not meet regulatory requirements for write-once storage.

* `PRIVATE_DATASETS_SHARE`: The server uses the specified `dataset_naming_schema` in `options` to make a new ZFS       dataset when the client connects. The server uses this dataset as the share path during the SMB session.

* `EXTERNAL_SHARE`: The SMB share is a DFS proxy to a share hosted on an external SMB server.

* `VEEAM_REPOSITORY_SHARE`: The SMB share is a repository for Veeam Backup & Replication and supports Fast Clone.
  NOTE: This feature is available only for TrueNAS Enterprise customers.

* `FCP_SHARE`: The SMB share is a used for Final Cut Pro storage. This feature automatically configures the share       to provide storage according to Apple support guidelines described in https://support.apple.com/en-ca/101919.       NOTE: `aapl_extensions` must be set in the global `smb.config`.       WARNING: This feature forcibly enables `aapl_name_mangling` on the SMB share which may cause unexpected behavior       for data that was written without this feature enabled. Valid values: `DEFAULT_SHARE`, `LEGACY_SHARE`, `TIMEMACHINE_SHARE` Default: `DEFAULT_SHARE`. Type: `string`
- `enabled` (Optional) - If unset, the SMB share is not available over the SMB protocol.  Default: `True`. Type: `boolean`
- `comment` (Optional) - Text field that is seen next to a share when an SMB client requests a list of SMB shares on the TrueNAS     server.  Default: ``. Type: `string`
- `readonly` (Optional) - If set, SMB clients cannot create or change files and directories in the SMB share.

NOTE: If set, the share path is still writeable by local processes or other file sharing protocols.  Default: `False`. Type: `boolean`
- `browsable` (Optional) - If set, the share is included when an SMB client requests a list of SMB shares on the TrueNAS server.  Default: `True`. Type: `boolean`
- `access_based_share_enumeration` (Optional) - If set, the share is only included when an SMB client requests a list of shares on the SMB server if     the share (not filesystem) access control list (see `sharing.smb.getacl`) grants access to the user.  Default: `False`. Type: `boolean`
- `audit` (Optional) - Audit configuration for monitoring SMB share access and operations.. Type: `object`
- `options` (Optional) - Additional configuration related to the configured SMB share purpose. If null, then the default     options related to the share purpose will be applied. . Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_sharing_smb.example <id>
```
