---
page_title: "truenas_system_ntpserver Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS system_ntpserver data
---

# truenas_system_ntpserver (Data Source)

Retrieves TrueNAS system_ntpserver data

## Example Usage

```terraform
data "truenas_system_ntpserver" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the system_ntpserver to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
