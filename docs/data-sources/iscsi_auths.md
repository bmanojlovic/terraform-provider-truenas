---
page_title: "truenas_iscsi_auths Data Source - terraform-provider-truenas"
subcategory: "iSCSI"
description: |-
  Query iscsi.auth
---

# truenas_iscsi_auths (Data Source)

Query iscsi.auth

## Example Usage

```terraform
data "truenas_iscsi_auths" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_auths to retrieve.

### Read-Only

- None
