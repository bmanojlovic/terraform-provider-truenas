---
page_title: "truenas_iscsi_target Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS iscsi_target data
---

# truenas_iscsi_target (Data Source)

Retrieves TrueNAS iscsi_target data

## Example Usage

```terraform
data "truenas_iscsi_target" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_target to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
