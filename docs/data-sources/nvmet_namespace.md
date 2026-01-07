---
page_title: "truenas_nvmet_namespace Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS nvmet_namespace data
---

# truenas_nvmet_namespace (Data Source)

Retrieves TrueNAS nvmet_namespace data

## Example Usage

```terraform
data "truenas_nvmet_namespace" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_namespace to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
