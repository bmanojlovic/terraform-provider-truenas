---
page_title: "truenas_iscsi_targetextent Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an Associated Target.
---

# truenas_iscsi_targetextent (Resource)

Create an Associated Target.

## Example Usage

```terraform
resource "truenas_iscsi_targetextent" "example" {
  extent = 1
  target = 1
}
```

## Schema

### Required

- `extent` (Int64) - ID of the iSCSI extent to associate with the target.
- `target` (Int64) - ID of the iSCSI target to associate with the extent.

### Optional

- `lunid` (Int64) - LUN ID to assign or `null` to auto-assign the next available LUN. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_targetextent.example <id>
```
