---
page_title: "truenas_nvmet_port_subsyss Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query nvmet.port_subsys
---

# truenas_nvmet_port_subsyss (Data Source)

Query nvmet.port_subsys

## Example Usage

```terraform
data "truenas_nvmet_port_subsyss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_port_subsyss to retrieve.

### Read-Only

- None
