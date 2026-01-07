---
page_title: "truenas_nvmet_host Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS nvmet_host data
---

# truenas_nvmet_host (Data Source)

Retrieves TrueNAS nvmet_host data

## Example Usage

```terraform
data "truenas_nvmet_host" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_host to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
