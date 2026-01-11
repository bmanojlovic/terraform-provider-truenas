---
page_title: "truenas_api_key Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Creates API Key.
---

# truenas_api_key (Resource)

Creates API Key.


## Example Usage

```terraform
resource "truenas_api_key" "example" {
  username = "example-value"
}
```

## Schema

### Required

- `username` (String) - 

### Optional

- `expires_at` (String) - Expiration timestamp for the API key or `null` for no expiration. Default: `None`
- `name` (String) - Human-readable name for the API key. Default: `nobody`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_api_key.example <id>
```
