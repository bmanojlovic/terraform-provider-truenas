---
page_title: "truenas_acme_dns_authenticator Resource - terraform-provider-truenas"
subcategory: "Certificates"
description: |-
  Create a DNS Authenticator
---

# truenas_acme_dns_authenticator (Resource)

Create a DNS Authenticator


## Example Usage

```terraform
resource "truenas_acme_dns_authenticator" "example" {
  attributes = "example"
  name = "example"
}
```

## Schema

### Required

- `attributes` (String) - Authentication credentials and configuration for the DNS provider. **Note:** This is a JSON object. Use `jsonencode()` to pass structured data. Example: `jsonencode({authenticator = "value", cloudflare_email = "value", api_key = "value", ...})`
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
