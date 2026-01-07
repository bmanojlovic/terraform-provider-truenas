---
page_title: "truenas_iscsi_targetextent Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Target-to-extent association configuration data for creation.
---

# truenas_iscsi_targetextent (Resource)

Target-to-extent association configuration data for creation.

## Example Usage

```terraform
resource "truenas_iscsi_targetextent" "example" {
  target = 1
  extent = 1
}
```

## Schema

### Required

- `target` (Required) - ID of the iSCSI target to associate with the extent.. Type: `integer`
- `extent` (Required) - ID of the iSCSI extent to associate with the target.. Type: `integer`

### Optional

- `lunid` (Optional) - LUN ID to assign or `null` to auto-assign the next available LUN.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_targetextent.example <id>
```
