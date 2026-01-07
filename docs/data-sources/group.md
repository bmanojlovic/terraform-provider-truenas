---
page_title: "truenas_group Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS group data
---

# truenas_group (Data Source)

Retrieves TrueNAS group data

## Example Usage

```terraform
data "truenas_group" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the group to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
