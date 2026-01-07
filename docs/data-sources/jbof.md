---
page_title: "truenas_jbof Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS jbof data
---

# truenas_jbof (Data Source)

Retrieves TrueNAS jbof data

## Example Usage

```terraform
data "truenas_jbof" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the jbof to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
