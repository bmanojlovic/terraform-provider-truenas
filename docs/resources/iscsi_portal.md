---
page_title: "truenas_iscsi_portal Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new iSCSI Portal.
---

# truenas_iscsi_portal (Resource)

Create a new iSCSI Portal.

## Example Usage

```terraform
resource "truenas_iscsi_portal" "example" {
  listen = ["item1"]
}
```

## Schema

### Required

- `listen` (List) - Array of IP addresses for the portal to listen on.

### Optional

- `comment` (String) - Optional comment describing the portal. Default: ``

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_portal.example <id>
```
