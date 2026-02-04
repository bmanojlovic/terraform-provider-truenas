---
page_title: "truenas_dnss Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query Name Servers with `query-filters` and `query-options`.
---

# truenas_dnss (Data Source)

Query Name Servers with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_dnss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the dnss to retrieve.

### Read-Only

- None
