---
page_title: "truenas_iscsi_portal Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS iscsi_portal data
---

# truenas_iscsi_portal (Data Source)

Retrieves TrueNAS iscsi_portal data

## Example Usage

```terraform
data "truenas_iscsi_portal" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_portal to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
