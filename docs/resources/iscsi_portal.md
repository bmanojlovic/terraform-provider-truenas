---
page_title: "truenas_iscsi_portal Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  iSCSI portal configuration data for creation.
---

# truenas_iscsi_portal (Resource)

iSCSI portal configuration data for creation.

## Example Usage

```terraform
resource "truenas_iscsi_portal" "example" {
  listen = []
  comment = ""
}
```

## Schema

### Required

- `listen` (Required) - Array of IP addresses for the portal to listen on.. Type: `array`

### Optional

- `comment` (Optional) - Optional comment describing the portal. Default: ``. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_portal.example <id>
```
