---
page_title: "truenas_filesystem_acltemplate Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new filesystem ACL template.
---

# truenas_filesystem_acltemplate (Resource)

Create a new filesystem ACL template.

## Example Usage

```terraform
resource "truenas_filesystem_acltemplate" "example" {
  acl = "example-value"
  acltype = "example-value"
  name = "example-value"
}
```

## Schema

### Required

- `acl` (String) - Array of Access Control Entries defined by this template.
- `acltype` (String) - ACL type this template provides. Valid values: `NFS4`, `POSIX1E`
- `name` (String) - Human-readable name for the ACL template.

### Optional

- `comment` (String) - Optional descriptive comment about the template's purpose. Default: ``

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_filesystem_acltemplate.example <id>
```
