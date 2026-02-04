---
page_title: "truenas_replications Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query replication
---

# truenas_replications (Data Source)

Query replication

## Example Usage

```terraform
data "truenas_replications" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the replications to retrieve.

### Read-Only

- None
