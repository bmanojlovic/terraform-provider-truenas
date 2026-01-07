---
page_title: "truenas_replication Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS replication data
---

# truenas_replication (Data Source)

Retrieves TrueNAS replication data

## Example Usage

```terraform
data "truenas_replication" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the replication to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
