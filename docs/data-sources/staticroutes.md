---
page_title: "truenas_staticroutes Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query staticroute
---

# truenas_staticroutes (Data Source)

Query staticroute

## Example Usage

```terraform
data "truenas_staticroutes" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the staticroutes to retrieve.

### Read-Only

- None
