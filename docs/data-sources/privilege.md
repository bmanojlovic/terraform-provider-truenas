---
page_title: "truenas_privilege Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS privilege data
---

# truenas_privilege (Data Source)

Retrieves TrueNAS privilege data

## Example Usage

```terraform
data "truenas_privilege" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the privilege to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
