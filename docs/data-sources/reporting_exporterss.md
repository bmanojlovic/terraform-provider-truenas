---
page_title: "truenas_reporting_exporterss Data Source - terraform-provider-truenas"
subcategory: "Reporting"
description: |-
  Query reporting.exporters
---

# truenas_reporting_exporterss (Data Source)

Query reporting.exporters

## Example Usage

```terraform
data "truenas_reporting_exporterss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the reporting_exporterss to retrieve.

### Read-Only

- None
