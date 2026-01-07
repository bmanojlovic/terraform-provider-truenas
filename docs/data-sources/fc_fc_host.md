---
page_title: "truenas_fc_fc_host Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS fc_fc_host data
---

# truenas_fc_fc_host (Data Source)

Retrieves TrueNAS fc_fc_host data

## Example Usage

```terraform
data "truenas_fc_fc_host" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the fc_fc_host to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
