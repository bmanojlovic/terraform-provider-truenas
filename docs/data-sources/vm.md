---
page_title: "truenas_vm Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS vm data
---

# truenas_vm (Data Source)

Retrieves TrueNAS vm data

## Example Usage

```terraform
data "truenas_vm" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the vm to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
