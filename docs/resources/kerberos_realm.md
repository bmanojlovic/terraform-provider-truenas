---
page_title: "truenas_kerberos_realm Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new kerberos realm. This will be automatically populated during the
---

# truenas_kerberos_realm (Resource)

Create a new kerberos realm. This will be automatically populated during the


## Example Usage

```terraform
resource "truenas_kerberos_realm" "example" {
  realm = "example-value"
}
```

## Schema

### Required

- `realm` (String) - Kerberos realm name. This is external to TrueNAS and is case-sensitive.     The general convention for kerberos realms is that they are upper-case.

### Optional

- `admin_server` (List) - List of kerberos admin servers. If the list is empty then the kerberos     libraries will use DNS to look them up. Default: `[]`
- `kdc` (List) - List of kerberos domain controllers. If the list is empty then the kerberos     libraries will use DNS to look up KDCs. In some situations this is undesirable     as kerberos libraries are, for intanc Default: `[]`
- `kpasswd_server` (List) - List of kerberos kpasswd servers. If the list is empty then DNS will be used     to look them up if needed. Default: `[]`
- `primary_kdc` (String) - The master Kerberos domain controller for this realm. TrueNAS uses this as a fallback if it cannot get     credentials because of an invalid password. This can help in environments where the domain us Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_kerberos_realm.example <id>
```
