---
page_title: "truenas_sharing_smbs Data Source - terraform-provider-truenas"
subcategory: "Sharing"
description: |-
  Query sharing.smb
---

# truenas_sharing_smbs (Data Source)

Query sharing.smb

## Example Usage

```terraform
data "truenas_sharing_smbs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the sharing_smbs to retrieve.

### Read-Only

- None
