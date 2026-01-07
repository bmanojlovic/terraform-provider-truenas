---
page_title: "truenas_vmware Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS vmware data
---

# truenas_vmware (Data Source)

Retrieves TrueNAS vmware data

## Example Usage

```terraform
data "truenas_vmware" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vmware to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
