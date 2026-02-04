---
page_title: "truenas_apps Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query all apps with `query-filters` and `query-options`.
---

# truenas_apps (Data Source)

Query all apps with `query-filters` and `query-options`.

## Example Usage

```terraform
data "truenas_apps" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the apps to retrieve.

### Read-Only

- None
