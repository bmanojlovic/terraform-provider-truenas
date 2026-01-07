---
page_title: "truenas_iscsi_auth Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  iSCSI authentication credential configuration data for creation.
---

# truenas_iscsi_auth (Resource)

iSCSI authentication credential configuration data for creation.

## Example Usage

```terraform
resource "truenas_iscsi_auth" "example" {
  tag = 1
  user = "example-user"
  secret = "example-secret"
  peeruser = ""
  peersecret = ""
  discovery_auth = "NONE"
}
```

## Schema

### Required

- `tag` (Required) - Numeric tag used to associate this credential with iSCSI targets.. Type: `integer`
- `user` (Required) - Username for iSCSI CHAP authentication.. Type: `string`
- `secret` (Required) - Password/secret for iSCSI CHAP authentication.. Type: `string`

### Optional

- `peeruser` (Optional) - Username for mutual CHAP authentication or empty string if not configured. Default: ``. Type: `string`
- `peersecret` (Optional) - Password/secret for mutual CHAP authentication or empty string if not configured. Default: ``. Type: `string`
- `discovery_auth` (Optional) - Authentication method for target discovery. If "CHAP_MUTUAL" is selected for target discovery, it is only     permitted for a single entry systemwide. Valid values: `NONE`, `CHAP`, `CHAP_MUTUAL` Default: `NONE`. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_auth.example <id>
```
