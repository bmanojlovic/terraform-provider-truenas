---
page_title: "truenas_nvmet_ports Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query nvmet.port
---

# truenas_nvmet_ports (Data Source)

Query nvmet.port

## Example Usage

```terraform
data "truenas_nvmet_ports" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_ports to retrieve.

### Read-Only

- None
