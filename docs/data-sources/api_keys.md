---
page_title: "truenas_api_keys Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query api_key
---

# truenas_api_keys (Data Source)

Query api_key

## Example Usage

```terraform
data "truenas_api_keys" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the api_keys to retrieve.

### Read-Only

- None
