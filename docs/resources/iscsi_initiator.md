---
page_title: "truenas_iscsi_initiator Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an iSCSI Initiator.
---

# truenas_iscsi_initiator (Resource)

Create an iSCSI Initiator.

## Example Usage

```terraform
resource "truenas_iscsi_initiator" "example" {
  # Configure required attributes
}
```

## Schema

### Required

- None

### Optional

- `comment` (String) - Optional comment describing the authorized initiator group. Default: ``
- `initiators` (List) - Array of iSCSI Qualified Names (IQNs) or IP addresses of authorized initiators. Default: `[]`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_iscsi_initiator.example <id>
```
