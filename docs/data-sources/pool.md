---
page_title: "truenas_pool Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS pool data
---

# truenas_pool (Data Source)

Retrieves TrueNAS pool data

## Example Usage

```terraform
data "truenas_pool" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
