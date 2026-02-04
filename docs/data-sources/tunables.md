---
page_title: "truenas_tunables Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query tunable
---

# truenas_tunables (Data Source)

Query tunable

## Example Usage

```terraform
data "truenas_tunables" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the tunables to retrieve.

### Read-Only

- None
