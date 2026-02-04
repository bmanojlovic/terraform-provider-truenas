---
page_title: "truenas_initshutdownscripts Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query initshutdownscript
---

# truenas_initshutdownscripts (Data Source)

Query initshutdownscript

## Example Usage

```terraform
data "truenas_initshutdownscripts" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the initshutdownscripts to retrieve.

### Read-Only

- None
