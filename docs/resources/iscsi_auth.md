---
page_title: "truenas_iscsi_auth Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an iSCSI Authorized Access.
---

# truenas_iscsi_auth (Resource)

Create an iSCSI Authorized Access.

## Example Usage

```terraform
resource "truenas_iscsi_auth" "example" {
  secret = "example-value"
  tag = 1
  user = "example-value"
}
```

## Schema

### Required

- `secret` (String) - Password/secret for iSCSI CHAP authentication.
- `tag` (Int64) - Numeric tag used to associate this credential with iSCSI targets.
- `user` (String) - Username for iSCSI CHAP authentication.

### Optional

- `discovery_auth` (String) - Authentication method for target discovery. If "CHAP_MUTUAL" is selected for target discovery, it is only     permitted for a single entry systemwide. Default: `NONE` Valid values: `NONE`, `CHAP`, `CHAP_MUTUAL`
- `peersecret` (String) - Password/secret for mutual CHAP authentication or empty string if not configured. Default: ``
- `peeruser` (String) - Username for mutual CHAP authentication or empty string if not configured. Default: ``

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_auth.example <id>
```
