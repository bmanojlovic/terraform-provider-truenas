---
page_title: "truenas_boot_environments Data Source - terraform-provider-truenas"
subcategory: "System - Boot"
description: |-
  Query boot.environment
---

# truenas_boot_environments (Data Source)

Query boot.environment

## Example Usage

```terraform
data "truenas_boot_environments" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the boot_environments to retrieve.

### Read-Only

- None
