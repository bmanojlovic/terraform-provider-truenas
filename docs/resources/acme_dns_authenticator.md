---
page_title: "truenas_acme_dns_authenticator Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a DNS Authenticator
---

# truenas_acme_dns_authenticator (Resource)

Create a DNS Authenticator


## Example Usage

```terraform
resource "truenas_acme_dns_authenticator" "example" {
  attributes = "example-value"
  name = "example-value"
}
```

## Schema

### Required

- `attributes` (String) - Authentication credentials and configuration for the DNS provider.
- `name` (String) - Human-readable name for the DNS authenticator.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_acme_dns_authenticator.example <id>
```
