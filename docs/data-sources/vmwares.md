---
page_title: "truenas_vmwares Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query vmware
---

# truenas_vmwares (Data Source)

Query vmware

## Example Usage

```terraform
data "truenas_vmwares" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vmwares to retrieve.

### Read-Only

- None
