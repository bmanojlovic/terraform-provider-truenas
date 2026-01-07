---
page_title: "truenas_initshutdownscript Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS initshutdownscript data
---

# truenas_initshutdownscript (Data Source)

Retrieves TrueNAS initshutdownscript data

## Example Usage

```terraform
data "truenas_initshutdownscript" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the initshutdownscript to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
