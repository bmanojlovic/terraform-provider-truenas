---
page_title: "truenas_kerberos_realm Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS kerberos_realm data
---

# truenas_kerberos_realm (Data Source)

Retrieves TrueNAS kerberos_realm data

## Example Usage

```terraform
data "truenas_kerberos_realm" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the kerberos_realm to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
