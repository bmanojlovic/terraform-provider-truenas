---
page_title: "truenas_kerberos_keytabs Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query kerberos.keytab
---

# truenas_kerberos_keytabs (Data Source)

Query kerberos.keytab

## Example Usage

```terraform
data "truenas_kerberos_keytabs" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the kerberos_keytabs to retrieve.

### Read-Only

- None
