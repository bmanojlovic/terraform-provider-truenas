---
page_title: "truenas_vms Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query vm
---

# truenas_vms (Data Source)

Query vm

## Example Usage

```terraform
data "truenas_vms" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vms to retrieve.

### Read-Only

- None
