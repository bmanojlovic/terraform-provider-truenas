---
page_title: "truenas_sharing_nfss Data Source - terraform-provider-truenas"
subcategory: "Sharing"
description: |-
  Query sharing.nfs
---

# truenas_sharing_nfss (Data Source)

Query sharing.nfs

## Example Usage

```terraform
data "truenas_sharing_nfss" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the sharing_nfss to retrieve.

### Read-Only

- None
