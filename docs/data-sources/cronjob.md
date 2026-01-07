---
page_title: "truenas_cronjob Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS cronjob data
---

# truenas_cronjob (Data Source)

Retrieves TrueNAS cronjob data

## Example Usage

```terraform
data "truenas_cronjob" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cronjob to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
