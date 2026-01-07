---
page_title: "truenas_keychaincredential Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Credential configuration data for the new keychain entry.
---

# truenas_keychaincredential (Resource)

Credential configuration data for the new keychain entry.

## Example Usage

```terraform
resource "truenas_keychaincredential" "example" {
  # Configuration here
}
```

## Schema

### Required

- None

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_keychaincredential.example <id>
```
