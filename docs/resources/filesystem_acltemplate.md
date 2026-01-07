---
page_title: "truenas_filesystem_acltemplate Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  ACL template configuration data for the new template.
---

# truenas_filesystem_acltemplate (Resource)

ACL template configuration data for the new template.

## Example Usage

```terraform
resource "truenas_filesystem_acltemplate" "example" {
  name = "example-name"
  acltype = "NFS4"
  acl = "example-acl"
  comment = ""
}
```

## Schema

### Required

- `name` (Required) - Human-readable name for the ACL template.. Type: `string`
- `acltype` (Required) - ACL type this template provides. Valid values: `NFS4`, `POSIX1E`. Type: `string`
- `acl` (Required) - Array of Access Control Entries defined by this template.. Type: `string`

### Optional

- `comment` (Optional) - Optional descriptive comment about the template's purpose. Default: ``. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_filesystem_acltemplate.example <id>
```
