---
page_title: "truenas_pool_snapshottask Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS pool_snapshottask data
---

# truenas_pool_snapshottask (Data Source)

Retrieves TrueNAS pool_snapshottask data

## Example Usage

```terraform
data "truenas_pool_snapshottask" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the pool_snapshottask to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
