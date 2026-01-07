---
page_title: "truenas_interface Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS interface data
---

# truenas_interface (Data Source)

Retrieves TrueNAS interface data

## Example Usage

```terraform
data "truenas_interface" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the interface to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
