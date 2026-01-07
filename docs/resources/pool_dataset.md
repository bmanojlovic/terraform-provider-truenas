---
page_title: "truenas_pool_dataset Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration data for creating a new ZFS dataset.
---

# truenas_pool_dataset (Resource)

Configuration data for creating a new ZFS dataset.

## Example Usage

```terraform
resource "truenas_pool_dataset" "example" {
  # Configuration here
}
```

## Schema

### Required

- None

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_pool_dataset.example <id>
```
