---
page_title: "truenas_reporting_exporters Resource - terraform-provider-truenas"
subcategory: "Reporting"
description: |-
  Create a specific reporting exporter configuration containing required details for exporting reporting metrics.
---

# truenas_reporting_exporters (Resource)

Create a specific reporting exporter configuration containing required details for exporting reporting metrics.


## Example Usage

```terraform
resource "truenas_reporting_exporters" "example" {
  attributes = "example"
  enabled = true
  name = "example"
}
```

## Schema

### Required

- `attributes` (String) - Specific attributes for the exporter. **Note:** This is a JSON object. Use `jsonencode()` to pass structured data. Example: `jsonencode({exporter_type = "value", destination_ip = "value", destination_port = 0, ...})`
- `enabled` (Bool) - Whether this exporter is enabled and active.
- `name` (String) - User defined name of exporter configuration.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_reporting_exporters.example <id>
```
