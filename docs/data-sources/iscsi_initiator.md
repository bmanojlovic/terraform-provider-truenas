---
page_title: "truenas_iscsi_initiator Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS iscsi_initiator data
---

# truenas_iscsi_initiator (Data Source)

Retrieves TrueNAS iscsi_initiator data

## Example Usage

```terraform
data "truenas_iscsi_initiator" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_initiator to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
