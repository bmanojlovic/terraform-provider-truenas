---
page_title: "truenas_iscsi_target Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an iSCSI Target.
---

# truenas_iscsi_target (Resource)

Create an iSCSI Target.

## Example Usage

```terraform
resource "truenas_iscsi_target" "example" {
  name = "example-value"
}
```

## Schema

### Required

- `name` (String) - Name of the iSCSI target (maximum 120 characters).

### Optional

- `alias` (String) - Optional alias name for the iSCSI target. Default: `None`
- `auth_networks` (List) - Array of network addresses allowed to access this target. Default: `[]`
- `groups` (List) - Array of portal-initiator group associations for this target. Default: `[]`
- `iscsi_parameters` (String) - Optional iSCSI-specific parameters for this target. Default: `None`
- `mode` (String) - Protocol mode for the target.  * `ISCSI`: iSCSI protocol only * `FC`: Fibre Channel protocol only * `BOTH`: Both iSCSI and Fibre Channel protocols  Fibre Channel may only be selected on TrueNAS Enterp Default: `ISCSI` Valid values: `ISCSI`, `FC`, `BOTH`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_target.example <id>
```
