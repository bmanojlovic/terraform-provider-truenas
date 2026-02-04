---
page_title: "truenas_jbofs Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query jbof
---

# truenas_jbofs (Data Source)

Query jbof

## Example Usage

```terraform
data "truenas_jbofs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the jbofs to retrieve.

### Read-Only

- None
