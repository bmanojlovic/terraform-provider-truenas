---
page_title: "truenas_iscsi_portals Data Source - terraform-provider-truenas"
subcategory: "iSCSI"
description: |-
  Query iscsi.portal
---

# truenas_iscsi_portals (Data Source)

Query iscsi.portal

## Example Usage

```terraform
data "truenas_iscsi_portals" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_portals to retrieve.

### Read-Only

- None
