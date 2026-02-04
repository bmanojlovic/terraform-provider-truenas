---
page_title: "truenas_iscsi_initiators Data Source - terraform-provider-truenas"
subcategory: "iSCSI"
description: |-
  Query iscsi.initiator
---

# truenas_iscsi_initiators (Data Source)

Query iscsi.initiator

## Example Usage

```terraform
data "truenas_iscsi_initiators" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_initiators to retrieve.

### Read-Only

- None
