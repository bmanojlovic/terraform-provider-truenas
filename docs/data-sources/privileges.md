---
page_title: "truenas_privileges Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query privilege
---

# truenas_privileges (Data Source)

Query privilege

## Example Usage

```terraform
data "truenas_privileges" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the privileges to retrieve.

### Read-Only

- None
