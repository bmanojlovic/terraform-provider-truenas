---
page_title: "truenas_iscsi_auth Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS iscsi_auth data
---

# truenas_iscsi_auth (Data Source)

Retrieves TrueNAS iscsi_auth data

## Example Usage

```terraform
data "truenas_iscsi_auth" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_auth to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
