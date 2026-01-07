---
page_title: "truenas_filesystem_acltemplate Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS filesystem_acltemplate data
---

# truenas_filesystem_acltemplate (Data Source)

Retrieves TrueNAS filesystem_acltemplate data

## Example Usage

```terraform
data "truenas_filesystem_acltemplate" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the filesystem_acltemplate to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
