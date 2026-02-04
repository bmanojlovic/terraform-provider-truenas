---
page_title: "truenas_cloudsyncs Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query cloudsync
---

# truenas_cloudsyncs (Data Source)

Query cloudsync

## Example Usage

```terraform
data "truenas_cloudsyncs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cloudsyncs to retrieve.

### Read-Only

- None
