---
page_title: "truenas_iscsi_target Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  iSCSI target configuration data for creation.
---

# truenas_iscsi_target (Resource)

iSCSI target configuration data for creation.

## Example Usage

```terraform
resource "truenas_iscsi_target" "example" {
  name = "example-name"
  mode = "ISCSI"
  groups = []
  auth_networks = []
}
```

## Schema

### Required

- `name` (Required) - Name of the iSCSI target (maximum 120 characters).. Type: `string`

### Optional

- `alias` (Optional) - Optional alias name for the iSCSI target.. Type: `string`
- `mode` (Optional) - Protocol mode for the target.

* `ISCSI`: iSCSI protocol only
* `FC`: Fibre Channel protocol only
* `BOTH`: Both iSCSI and Fibre Channel protocols

Fibre Channel may only be selected on TrueNAS Enterprise-licensed systems with a suitable Fibre Channel HBA. Valid values: `ISCSI`, `FC`, `BOTH` Default: `ISCSI`. Type: `string`
- `groups` (Optional) - Array of portal-initiator group associations for this target. Default: `[]`. Type: `array`
- `auth_networks` (Optional) - Array of network addresses allowed to access this target. Default: `[]`. Type: `array`
- `iscsi_parameters` (Optional) - Optional iSCSI-specific parameters for this target.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_target.example <id>
```
