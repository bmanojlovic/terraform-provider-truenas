---
page_title: "truenas_app_registry Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an app registry entry.
---

# truenas_app_registry (Resource)

Create an app registry entry.

## Example Usage

```terraform
resource "truenas_app_registry" "example" {
  name = "example-value"
  password = "example-value"
  username = "example-value"
}
```

## Schema

### Required

- `name` (String) - Human-readable name for the container registry.
- `password` (String) - Password or access token for registry authentication (masked for security).
- `username` (String) - Username for registry authentication (masked for security).

### Optional

- `description` (String) - Optional description of the container registry or `null`. Default: `None`
- `uri` (String) - Container registry URI endpoint (defaults to Docker Hub). Default: `https://index.docker.io/v1/`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_app_registry.example <id>
```
