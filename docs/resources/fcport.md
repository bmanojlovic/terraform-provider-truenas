---
page_title: "truenas_fcport Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Fibre Channel port configuration data for the new port.
---

# truenas_fcport (Resource)

Fibre Channel port configuration data for the new port.

## Example Usage

```terraform
resource "truenas_fcport" "example" {
  port = "example-port"
  target_id = 1
}
```

## Schema

### Required

- `port` (Required) - Alias name for the Fibre Channel port.. Type: `string`
- `target_id` (Required) - ID of the target to associate with this FC port.. Type: `integer`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_fcport.example <id>
```
