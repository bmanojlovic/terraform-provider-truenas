---
page_title: "truenas_iscsi_targetextents Data Source - terraform-provider-truenas"
subcategory: "iSCSI"
description: |-
  Query iscsi.targetextent
---

# truenas_iscsi_targetextents (Data Source)

Query iscsi.targetextent

## Example Usage

```terraform
data "truenas_iscsi_targetextents" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_targetextents to retrieve.

### Read-Only

- None
