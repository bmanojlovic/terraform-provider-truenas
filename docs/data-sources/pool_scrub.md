---
page_title: "truenas_pool_scrub Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS pool_scrub data
---

# truenas_pool_scrub (Data Source)

Retrieves TrueNAS pool_scrub data

## Example Usage

```terraform
data "truenas_pool_scrub" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_scrub to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
