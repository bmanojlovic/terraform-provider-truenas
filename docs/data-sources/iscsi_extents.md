---
page_title: "truenas_iscsi_extents Data Source - terraform-provider-truenas"
subcategory: "iSCSI"
description: |-
  Query iscsi.extent
---

# truenas_iscsi_extents (Data Source)

Query iscsi.extent

## Example Usage

```terraform
data "truenas_iscsi_extents" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_extents to retrieve.

### Read-Only

- None
