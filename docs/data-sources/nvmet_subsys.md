---
page_title: "truenas_nvmet_subsys Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS nvmet_subsys data
---

# truenas_nvmet_subsys (Data Source)

Retrieves TrueNAS nvmet_subsys data

## Example Usage

```terraform
data "truenas_nvmet_subsys" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_subsys to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
