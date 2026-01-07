---
page_title: "truenas_fc_fc_host Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Fibre Channel host configuration data for the new host.
---

# truenas_fc_fc_host (Resource)

Fibre Channel host configuration data for the new host.

## Example Usage

```terraform
resource "truenas_fc_fc_host" "example" {
  alias = "example-alias"
  npiv = 0
}
```

## Schema

### Required

- `alias` (Required) - Human-readable alias for the Fibre Channel host.. Type: `string`

### Optional

- `wwpn` (Optional) - World Wide Port Name for port A or `null` if not configured.. Type: `string`
- `wwpn_b` (Optional) - World Wide Port Name for port B or `null` if not configured.. Type: `string`
- `npiv` (Optional) - Number of N_Port ID Virtualization (NPIV) virtual ports to create. Default: `0`. Type: `integer`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_fc_fc_host.example <id>
```
