---
page_title: "truenas_fc_fc_hosts Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query fc.fc_host
---

# truenas_fc_fc_hosts (Data Source)

Query fc.fc_host

## Example Usage

```terraform
data "truenas_fc_fc_hosts" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the fc_fc_hosts to retrieve.

### Read-Only

- None
