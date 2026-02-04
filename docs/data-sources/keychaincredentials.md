---
page_title: "truenas_keychaincredentials Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query keychaincredential
---

# truenas_keychaincredentials (Data Source)

Query keychaincredential

## Example Usage

```terraform
data "truenas_keychaincredentials" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the keychaincredentials to retrieve.

### Read-Only

- None
