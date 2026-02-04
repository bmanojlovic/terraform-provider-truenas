---
page_title: "truenas_acme_dns_authenticators Data Source - terraform-provider-truenas"
subcategory: "Certificates"
description: |-
  Query acme.dns.authenticator
---

# truenas_acme_dns_authenticators (Data Source)

Query acme.dns.authenticator

## Example Usage

```terraform
data "truenas_acme_dns_authenticators" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the acme_dns_authenticators to retrieve.

### Read-Only

- None
