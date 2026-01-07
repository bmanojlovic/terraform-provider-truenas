---
page_title: "truenas_rsynctask Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS rsynctask data
---

# truenas_rsynctask (Data Source)

Retrieves TrueNAS rsynctask data

## Example Usage

```terraform
data "truenas_rsynctask" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the rsynctask to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
