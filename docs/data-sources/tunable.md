---
page_title: "truenas_tunable Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS tunable data
---

# truenas_tunable (Data Source)

Retrieves TrueNAS tunable data

## Example Usage

```terraform
data "truenas_tunable" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the tunable to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
