---
page_title: "truenas_nvmet_host_subsyss Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query nvmet.host_subsys
---

# truenas_nvmet_host_subsyss (Data Source)

Query nvmet.host_subsys

## Example Usage

```terraform
data "truenas_nvmet_host_subsyss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_host_subsyss to retrieve.

### Read-Only

- None
