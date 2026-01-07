---
page_title: "truenas_pool_dataset Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS pool_dataset data
---

# truenas_pool_dataset (Data Source)

Retrieves TrueNAS pool_dataset data

## Example Usage

```terraform
data "truenas_pool_dataset" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_dataset to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
