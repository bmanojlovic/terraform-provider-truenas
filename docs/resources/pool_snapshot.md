---
page_title: "truenas_pool_snapshot Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a snapshot with either an explicit name or naming schema.
---

# truenas_pool_snapshot (Resource)

Configuration for creating a snapshot with either an explicit name or naming schema.

## Example Usage

```terraform
resource "truenas_pool_snapshot" "example" {
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
terraform import truenas_pool_snapshot.example <id>
```
