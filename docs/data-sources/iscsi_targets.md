---
page_title: "truenas_iscsi_targets Data Source - terraform-provider-truenas"
subcategory: "iSCSI"
description: |-
  Query iscsi.target
---

# truenas_iscsi_targets (Data Source)

Query iscsi.target

## Example Usage

```terraform
data "truenas_iscsi_targets" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_targets to retrieve.

### Read-Only

- None
