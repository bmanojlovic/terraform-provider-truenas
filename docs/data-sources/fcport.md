---
page_title: "truenas_fcport Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS fcport data
---

# truenas_fcport (Data Source)

Retrieves TrueNAS fcport data

## Example Usage

```terraform
data "truenas_fcport" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the fcport to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
