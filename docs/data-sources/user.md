---
page_title: "truenas_user Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS user data
---

# truenas_user (Data Source)

Retrieves TrueNAS user data

## Example Usage

```terraform
data "truenas_user" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the user to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
