---
page_title: "truenas_app Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  AppCreateArgs parameters.
---

# truenas_app (Resource)

AppCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_app" "example" {
  app_name = "example-app_name"
  start_on_create = true
  custom_app = false
  custom_compose_config_string = ""
  train = "stable"
  version = "latest"
}
```

## Schema

### Required

- `app_name` (Required) - Application name must have the following:

* Lowercase alphanumeric characters can be specified.
* Name must start with an alphabetic character and can end with alphanumeric character.
* Hyphen '-' is allowed but not as the first or last character.. Type: `string`

### Optional

- `start_on_create` (Optional) - Automatically start after creation. Default: `true`. Type: `boolean`
- `custom_app` (Optional) - Whether to create a custom application (`true`) or install from catalog (`false`). Default: `False`. Type: `boolean`
- `values` (Optional) - Configuration values for the application installation.. Type: `object`
- `custom_compose_config` (Optional) - Docker Compose configuration as a structured object for custom applications.. Type: `object`
- `custom_compose_config_string` (Optional) - Docker Compose configuration as a YAML string for custom applications. Default: ``. Type: `string`
- `catalog_app` (Optional) - Name of the catalog application to install. Required when `custom_app` is `false`.. Type: `string`
- `train` (Optional) - The catalog train to install from. Default: `stable`. Type: `string`
- `version` (Optional) - The version of the application to install. Default: `latest`. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_app.example <id>
```
