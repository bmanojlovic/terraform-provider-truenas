---
page_title: "truenas_certificate Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS certificate data
---

# truenas_certificate (Data Source)

Retrieves TrueNAS certificate data

## Example Usage

```terraform
data "truenas_certificate" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the certificate to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
