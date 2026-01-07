---
page_title: "truenas_privilege Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new privilege.
---

# truenas_privilege (Resource)

Configuration for creating a new privilege.

## Example Usage

```terraform
resource "truenas_privilege" "example" {
  name = "example-name"
  web_shell = false
  local_groups = []
  ds_groups = []
  roles = []
}
```

## Schema

### Required

- `name` (Required) - Display name of the privilege.. Type: `string`
- `web_shell` (Required) - Whether this privilege grants access to the web shell.. Type: `boolean`

### Optional

- `local_groups` (Optional) - Array of local group IDs to assign to this privilege. Default: `[]`. Type: `array`
- `ds_groups` (Optional) - Array of directory service group IDs or SIDs to assign to this privilege. Default: `[]`. Type: `array`
- `roles` (Optional) - Array of role names included in this privilege. Default: `[]`. Type: `array`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_privilege.example <id>
```
