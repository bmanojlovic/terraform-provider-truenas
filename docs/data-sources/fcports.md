---
page_title: "truenas_fcports Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query fcport
---

# truenas_fcports (Data Source)

Query fcport

## Example Usage

```terraform
data "truenas_fcports" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the fcports to retrieve.

### Read-Only

- None
