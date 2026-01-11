---
page_title: "truenas_alertservice Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an Alert Service of specified `type`.
---

# truenas_alertservice (Resource)

Create an Alert Service of specified `type`.


## Example Usage

```terraform
resource "truenas_alertservice" "example" {
  attributes = "example-value"
  level = "example-value"
  name = "example-value"
}
```

## Schema

### Required

- `attributes` (String) - Service-specific configuration attributes (credentials, endpoints, etc.).
- `level` (String) - Minimum alert severity level that triggers notifications through this service. Valid values: `INFO`, `NOTICE`, `WARNING`, `ERROR`, `CRITICAL`, `ALERT`, `EMERGENCY`
- `name` (String) - Human-readable name for the alert service.

### Optional

- `enabled` (Bool) - Whether the alert service is active and will send notifications. Default: `True`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_alertservice.example <id>
```
