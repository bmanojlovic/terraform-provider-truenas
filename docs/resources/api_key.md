---
page_title: "truenas_api_key Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  API key configuration data for the new key.
---

# truenas_api_key (Resource)

API key configuration data for the new key.

## Example Usage

```terraform
resource "truenas_api_key" "example" {
  username = "example-username"
  name = "nobody"
}
```

## Schema

### Required

- `username` (Required) - username configuration. Type: `string`

### Optional

- `name` (Optional) - Human-readable name for the API key. Default: `nobody`. Type: `string`
- `expires_at` (Optional) - Expiration timestamp for the API key or `null` for no expiration.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_api_key.example <id>
```
