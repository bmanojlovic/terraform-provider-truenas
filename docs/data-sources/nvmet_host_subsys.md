---
page_title: "truenas_nvmet_host_subsys Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS nvmet_host_subsys data
---

# truenas_nvmet_host_subsys (Data Source)

Retrieves TrueNAS nvmet_host_subsys data

## Example Usage

```terraform
data "truenas_nvmet_host_subsys" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_host_subsys to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
