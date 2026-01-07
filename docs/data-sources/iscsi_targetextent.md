---
page_title: "truenas_iscsi_targetextent Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS iscsi_targetextent data
---

# truenas_iscsi_targetextent (Data Source)

Retrieves TrueNAS iscsi_targetextent data

## Example Usage

```terraform
data "truenas_iscsi_targetextent" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the iscsi_targetextent to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
