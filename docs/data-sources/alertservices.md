---
page_title: "truenas_alertservices Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query alertservice
---

# truenas_alertservices (Data Source)

Query alertservice

## Example Usage

```terraform
data "truenas_alertservices" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the alertservices to retrieve.

### Read-Only

- None
