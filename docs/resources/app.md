---
page_title: "truenas_app Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an app with `app_name` using `catalog_app` with `train` and `version`.
---

# truenas_app (Resource)

Create an app with `app_name` using `catalog_app` with `train` and `version`.

## Example Usage

```terraform
resource "truenas_app" "example" {
  app_name = "example-value"
  start_on_create = true
}
```

## Schema

### Required

- `app_name` (String) - Application name must have the following:  * Lowercase alphanumeric characters can be specified. * Name must start with an alphabetic character and can end with alphanumeric character. * Hyphen '-' is

### Optional

- `start_on_create` (Bool) - Start the resource immediately after creation. Default: `true`
- `catalog_app` (String) - Name of the catalog application to install. Required when `custom_app` is `false`. Default: `None`
- `custom_app` (Bool) - Whether to create a custom application (`true`) or install from catalog (`false`). Default: `False`
- `custom_compose_config` (String) - Docker Compose configuration as a structured object for custom applications.
- `custom_compose_config_string` (String) - Docker Compose configuration as a YAML string for custom applications. Default: ``
- `train` (String) - The catalog train to install from. Default: `stable`
- `values` (String) - Configuration values for the application installation.
- `version` (String) - The version of the application to install. Default: `latest`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_app.example <id>
```
