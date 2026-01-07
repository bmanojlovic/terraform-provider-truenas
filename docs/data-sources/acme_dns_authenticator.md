---
page_title: "truenas_acme_dns_authenticator Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS acme_dns_authenticator data
---

# truenas_acme_dns_authenticator (Data Source)

Retrieves TrueNAS acme_dns_authenticator data

## Example Usage

```terraform
data "truenas_acme_dns_authenticator" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the acme_dns_authenticator to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
