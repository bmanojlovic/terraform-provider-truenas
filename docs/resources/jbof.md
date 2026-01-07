---
page_title: "truenas_jbof Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  JBOF configuration data for creation.
---

# truenas_jbof (Resource)

JBOF configuration data for creation.

## Example Usage

```terraform
resource "truenas_jbof" "example" {
  mgmt_ip1 = "example-mgmt_ip1"
  mgmt_username = "example-mgmt_username"
  mgmt_password = "example-mgmt_password"
}
```

## Schema

### Required

- `mgmt_ip1` (Required) - IP of first Redfish management interface.. Type: `string`
- `mgmt_username` (Required) - Redfish administrative username.. Type: `string`
- `mgmt_password` (Required) - Redfish administrative password.. Type: `string`

### Optional

- `description` (Optional) - Optional description of the JBOF.. Type: `string`
- `mgmt_ip2` (Optional) - Optional IP of second Redfish management interface.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_jbof.example <id>
```
