---
page_title: "truenas_service Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS service data
---

# truenas_service (Data Source)

Retrieves TrueNAS service data

## Example Usage

```terraform
data "truenas_service" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the service to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
