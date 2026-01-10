---
page_title: "truenas_privilege Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Creates a privilege.
---

# truenas_privilege (Resource)

Creates a privilege.

## Example Usage

```terraform
resource "truenas_privilege" "example" {
  name = "example-value"
  web_shell = true
}
```

## Schema

### Required

- `name` (String) - Display name of the privilege.
- `web_shell` (Bool) - Whether this privilege grants access to the web shell.

### Optional

- `ds_groups` (List) - Array of directory service group IDs or SIDs to assign to this privilege. Default: `[]`
- `local_groups` (List) - Array of local group IDs to assign to this privilege. Default: `[]`
- `roles` (List) - Array of role names included in this privilege. Default: `[]`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_privilege.example <id>
```
