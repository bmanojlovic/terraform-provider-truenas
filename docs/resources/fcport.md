---
page_title: "truenas_fcport Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Creates mapping between a FC port and a target.
---

# truenas_fcport (Resource)

Creates mapping between a FC port and a target.


## Example Usage

```terraform
resource "truenas_fcport" "example" {
  port = "example-value"
  target_id = 1
}
```

## Schema

### Required

- `port` (String) - Alias name for the Fibre Channel port.
- `target_id` (Int64) - ID of the target to associate with this FC port.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_fcport.example <id>
```
