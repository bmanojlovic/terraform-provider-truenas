---
page_title: "truenas_app_registry Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Container registry configuration data for the new registry.
---

# truenas_app_registry (Resource)

Container registry configuration data for the new registry.

## Example Usage

```terraform
resource "truenas_app_registry" "example" {
  name = "example-name"
  username = "example-username"
  password = "example-password"
  uri = "https://index.docker.io/v1/"
}
```

## Schema

### Required

- `name` (Required) - Human-readable name for the container registry.. Type: `string`
- `username` (Required) - Username for registry authentication (masked for security).. Type: `string`
- `password` (Required) - Password or access token for registry authentication (masked for security).. Type: `string`

### Optional

- `description` (Optional) - Optional description of the container registry or `null`.. Type: `string`
- `uri` (Optional) - Container registry URI endpoint (defaults to Docker Hub). Default: `https://index.docker.io/v1/`. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_app_registry.example <id>
```
