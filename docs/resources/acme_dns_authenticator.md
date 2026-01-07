---
page_title: "truenas_acme_dns_authenticator Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  DNSAuthenticatorCreateArgs parameters.
---

# truenas_acme_dns_authenticator (Resource)

DNSAuthenticatorCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_acme_dns_authenticator" "example" {
  attributes = "example-attributes"
  name = "example-name"
}
```

## Schema

### Required

- `attributes` (Required) - Authentication credentials and configuration for the DNS provider.. Type: `string`
- `name` (Required) - Human-readable name for the DNS authenticator.. Type: `string`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_acme_dns_authenticator.example <id>
```
