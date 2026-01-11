---
page_title: "truenas_keychaincredential Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a Keychain Credential.
---

# truenas_keychaincredential (Resource)

Create a Keychain Credential.


## Variants

This resource has **2 variants** controlled by the `type` field. Choose the appropriate variant for your use case:

### SSH_CREDENTIALS

```terraform
resource "truenas_keychaincredential" "example" {
  type = "SSH_CREDENTIALS"
  attributes = "value"
  name = "value"
}
```

**Required fields:** `attributes`, `name`, `type`

### SSH_KEY_PAIR

```terraform
resource "truenas_keychaincredential" "example" {
  type = "SSH_KEY_PAIR"
  attributes = "value"
  name = "value"
}
```

**Required fields:** `attributes`, `name`, `type`



## Schema

### Required

- `attributes` (String) - SSH connection attributes including host, authentication, and connection settings.
- `name` (String) - Distinguishes this Keychain Credential from others.
- `type` (String) - Keychain credential type identifier for SSH connection credentials. Valid values: `SSH_CREDENTIALS`, `SSH_KEY_PAIR`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_keychaincredential.example <id>
```
