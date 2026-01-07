---
page_title: "truenas_nvmet_port_subsys Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS nvmet_port_subsys data
---

# truenas_nvmet_port_subsys (Data Source)

Retrieves TrueNAS nvmet_port_subsys data

## Example Usage

```terraform
data "truenas_nvmet_port_subsys" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_port_subsys to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
