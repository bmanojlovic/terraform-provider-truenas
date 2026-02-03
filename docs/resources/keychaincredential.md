---
page_title: "truenas_keychaincredential Resource - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Create a Keychain Credential.
---

# truenas_keychaincredential (Resource)

Create a Keychain Credential.

## Variants

This resource has **2 variants** controlled by the `type` field.

### SSH_KEY_PAIR

```terraform
resource "truenas_keychaincredential" "example" {
  type = "SSH_KEY_PAIR"
  attributes = "value"
  name = "value"
}
```

**Required fields:** `attributes`, `name`, `type`

### SSH_CREDENTIALS

```terraform
resource "truenas_keychaincredential" "example" {
  type = "SSH_CREDENTIALS"
  attributes = "value"
  name = "value"
}
```

**Required fields:** `attributes`, `name`, `type`



## Schema

### Required

- `attributes` (String) - SSH connection attributes including host, authentication, and connection settings. **Note:** This is a JSON object. Use `jsonencode()` to pass structured data. Example: `jsonencode({host = "value", port = 0, username = "value", ...})`
- `name` (String) - Distinguishes this Keychain Credential from others.
- `type` (String) - Keychain credential type identifier for SSH connection credentials. Valid values: `SSH_KEY_PAIR`, `SSH_CREDENTIALS`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_keychaincredential.example <id>
```
