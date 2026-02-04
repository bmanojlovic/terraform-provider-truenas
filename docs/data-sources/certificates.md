---
page_title: "truenas_certificates Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query certificate
---

# truenas_certificates (Data Source)

Query certificate

## Example Usage

```terraform
data "truenas_certificates" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the certificates to retrieve.

### Read-Only

- None
