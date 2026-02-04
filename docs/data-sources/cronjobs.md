---
page_title: "truenas_cronjobs Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query cronjob
---

# truenas_cronjobs (Data Source)

Query cronjob

## Example Usage

```terraform
data "truenas_cronjobs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the cronjobs to retrieve.

### Read-Only

- None
