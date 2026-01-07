---
page_title: "truenas_alertservice Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS alertservice data
---

# truenas_alertservice (Data Source)

Retrieves TrueNAS alertservice data

## Example Usage

```terraform
data "truenas_alertservice" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the alertservice to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
