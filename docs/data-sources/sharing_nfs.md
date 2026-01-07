---
page_title: "truenas_sharing_nfs Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS sharing_nfs data
---

# truenas_sharing_nfs (Data Source)

Retrieves TrueNAS sharing_nfs data

## Example Usage

```terraform
data "truenas_sharing_nfs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the sharing_nfs to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
