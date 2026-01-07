---
page_title: "truenas_alertservice Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Alert service configuration data for the new service.
---

# truenas_alertservice (Resource)

Alert service configuration data for the new service.

## Example Usage

```terraform
resource "truenas_alertservice" "example" {
  name = "example-name"
  attributes = "example-attributes"
  level = "INFO"
  enabled = true
}
```

## Schema

### Required

- `name` (Required) - Human-readable name for the alert service.. Type: `string`
- `attributes` (Required) - Service-specific configuration attributes (credentials, endpoints, etc.).. Type: `string`
- `level` (Required) - Minimum alert severity level that triggers notifications through this service. Valid values: `INFO`, `NOTICE`, `WARNING`. Type: `string`

### Optional

- `enabled` (Optional) - Whether the alert service is active and will send notifications. Default: `True`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_alertservice.example <id>
```
