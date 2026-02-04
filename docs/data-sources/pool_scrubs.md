---
page_title: "truenas_pool_scrubs Data Source - terraform-provider-truenas"
subcategory: "Storage - Pools"
description: |-
  Query pool.scrub
---

# truenas_pool_scrubs (Data Source)

Query pool.scrub

## Example Usage

```terraform
data "truenas_pool_scrubs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_scrubs to retrieve.

### Read-Only

- None
