---
page_title: "truenas_nvmet_hosts Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query nvmet.host
---

# truenas_nvmet_hosts (Data Source)

Query nvmet.host

## Example Usage

```terraform
data "truenas_nvmet_hosts" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_hosts to retrieve.

### Read-Only

- None
