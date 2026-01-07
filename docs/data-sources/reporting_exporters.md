---
page_title: "truenas_reporting_exporters Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS reporting_exporters data
---

# truenas_reporting_exporters (Data Source)

Retrieves TrueNAS reporting_exporters data

## Example Usage

```terraform
data "truenas_reporting_exporters" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the reporting_exporters to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
