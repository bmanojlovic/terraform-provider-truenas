---
page_title: "truenas_kerberos_keytab Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS kerberos_keytab data
---

# truenas_kerberos_keytab (Data Source)

Retrieves TrueNAS kerberos_keytab data

## Example Usage

```terraform
data "truenas_kerberos_keytab" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the kerberos_keytab to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
