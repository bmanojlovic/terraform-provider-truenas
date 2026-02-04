---
page_title: "truenas_disks Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query disks.
---

# truenas_disks (Data Source)

Query disks.

## Example Usage

```terraform
data "truenas_disks" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the disks to retrieve.

### Read-Only

- None
