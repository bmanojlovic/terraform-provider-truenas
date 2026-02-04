---
page_title: "truenas_rsynctasks Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query rsynctask
---

# truenas_rsynctasks (Data Source)

Query rsynctask

## Example Usage

```terraform
data "truenas_rsynctasks" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the rsynctasks to retrieve.

### Read-Only

- None
