---
page_title: "truenas_api_key Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS api_key data
---

# truenas_api_key (Data Source)

Retrieves TrueNAS api_key data

## Example Usage

```terraform
data "truenas_api_key" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the api_key to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
