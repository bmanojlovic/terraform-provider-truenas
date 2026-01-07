---
page_title: "truenas_sharing_smb Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS sharing_smb data
---

# truenas_sharing_smb (Data Source)

Retrieves TrueNAS sharing_smb data

## Example Usage

```terraform
data "truenas_sharing_smb" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the sharing_smb to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
