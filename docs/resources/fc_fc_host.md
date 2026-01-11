---
page_title: "truenas_fc_fc_host Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Creates FC host (pairing).
---

# truenas_fc_fc_host (Resource)

Creates FC host (pairing).


## Example Usage

```terraform
resource "truenas_fc_fc_host" "example" {
  alias = "example-value"
}
```

## Schema

### Required

- `alias` (String) - Human-readable alias for the Fibre Channel host.

### Optional

- `npiv` (Int64) - Number of N_Port ID Virtualization (NPIV) virtual ports to create. Default: `0`
- `wwpn` (String) - World Wide Port Name for port A or `null` if not configured. Default: `None`
- `wwpn_b` (String) - World Wide Port Name for port B or `null` if not configured. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_fc_fc_host.example <id>
```
