---
page_title: "truenas_filesystem_acltemplates Data Source - terraform-provider-truenas"
subcategory: "Filesystem"
description: |-
  Query filesystem.acltemplate
---

# truenas_filesystem_acltemplates (Data Source)

Query filesystem.acltemplate

## Example Usage

```terraform
data "truenas_filesystem_acltemplates" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the filesystem_acltemplates to retrieve.

### Read-Only

- None
