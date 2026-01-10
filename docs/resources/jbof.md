---
page_title: "truenas_jbof Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new JBOF.
---

# truenas_jbof (Resource)

Create a new JBOF.

## Example Usage

```terraform
resource "truenas_jbof" "example" {
  mgmt_ip1 = "example-value"
  mgmt_password = "example-value"
  mgmt_username = "example-value"
}
```

## Schema

### Required

- `mgmt_ip1` (String) - IP of first Redfish management interface.
- `mgmt_password` (String) - Redfish administrative password.
- `mgmt_username` (String) - Redfish administrative username.

### Optional

- `description` (String) - Optional description of the JBOF.
- `mgmt_ip2` (String) - Optional IP of second Redfish management interface.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_jbof.example <id>
```
