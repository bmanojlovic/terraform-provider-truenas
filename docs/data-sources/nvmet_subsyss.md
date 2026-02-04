---
page_title: "truenas_nvmet_subsyss Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query nvmet.subsys
---

# truenas_nvmet_subsyss (Data Source)

Query nvmet.subsys

## Example Usage

```terraform
data "truenas_nvmet_subsyss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_subsyss to retrieve.

### Read-Only

- None
