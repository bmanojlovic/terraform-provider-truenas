---
page_title: "truenas_staticroute Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS staticroute data
---

# truenas_staticroute (Data Source)

Retrieves TrueNAS staticroute data

## Example Usage

```terraform
data "truenas_staticroute" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the staticroute to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
