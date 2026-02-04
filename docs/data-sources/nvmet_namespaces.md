---
page_title: "truenas_nvmet_namespaces Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query nvmet.namespace
---

# truenas_nvmet_namespaces (Data Source)

Query nvmet.namespace

## Example Usage

```terraform
data "truenas_nvmet_namespaces" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the nvmet_namespaces to retrieve.

### Read-Only

- None
