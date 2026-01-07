---
page_title: "truenas_nvmet_port Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS nvmet_port data
---

# truenas_nvmet_port (Data Source)

Retrieves TrueNAS nvmet_port data

## Example Usage

```terraform
data "truenas_nvmet_port" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_port to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
