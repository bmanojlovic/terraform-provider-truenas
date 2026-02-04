---
page_title: "truenas_kerberos_realms Data Source - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Query kerberos.realm
---

# truenas_kerberos_realms (Data Source)

Query kerberos.realm

## Example Usage

```terraform
data "truenas_kerberos_realms" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the kerberos_realms to retrieve.

### Read-Only

- None
