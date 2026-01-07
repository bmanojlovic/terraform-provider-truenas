---
page_title: "truenas_reporting_exporters Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  ReportingExportsCreateArgs parameters.
---

# truenas_reporting_exporters (Resource)

ReportingExportsCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_reporting_exporters" "example" {
  enabled = false
  attributes = "example-attributes"
  name = "example-name"
}
```

## Schema

### Required

- `enabled` (Required) - Whether this exporter is enabled and active.. Type: `boolean`
- `attributes` (Required) - Specific attributes for the exporter.. Type: `string`
- `name` (Required) - User defined name of exporter configuration.. Type: `string`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_reporting_exporters.example <id>
```
