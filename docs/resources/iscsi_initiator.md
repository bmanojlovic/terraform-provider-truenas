---
page_title: "truenas_iscsi_initiator Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Authorized initiator group configuration data for creation.
---

# truenas_iscsi_initiator (Resource)

Authorized initiator group configuration data for creation.

## Example Usage

```terraform
resource "truenas_iscsi_initiator" "example" {
  initiators = []
  comment = ""
}
```

## Schema

### Required

- None

### Optional

- `initiators` (Optional) - Array of iSCSI Qualified Names (IQNs) or IP addresses of authorized initiators. Default: `[]`. Type: `array`
- `comment` (Optional) - Optional comment describing the authorized initiator group. Default: ``. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_initiator.example <id>
```
