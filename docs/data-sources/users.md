---
page_title: "truenas_users Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query users with `query-filters` and `query-options`.
---

# truenas_users (Data Source)

Query users with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_users" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the users to retrieve.

### Read-Only

- None
