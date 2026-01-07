---
page_title: "truenas_keychaincredential Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS keychaincredential data
---

# truenas_keychaincredential (Data Source)

Retrieves TrueNAS keychaincredential data

## Example Usage

```terraform
data "truenas_keychaincredential" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the keychaincredential to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
