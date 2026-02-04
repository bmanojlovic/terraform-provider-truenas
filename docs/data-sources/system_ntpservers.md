---
page_title: "truenas_system_ntpservers Data Source - terraform-provider-truenas"
subcategory: "System"
description: |-
  Query system.ntpserver
---

# truenas_system_ntpservers (Data Source)

Query system.ntpserver

## Example Usage

```terraform
data "truenas_system_ntpservers" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the system_ntpservers to retrieve.

### Read-Only

- None
