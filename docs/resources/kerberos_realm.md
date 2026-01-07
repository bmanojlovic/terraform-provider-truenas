---
page_title: "truenas_kerberos_realm Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Kerberos realm configuration data for creation.
---

# truenas_kerberos_realm (Resource)

Kerberos realm configuration data for creation.

## Example Usage

```terraform
resource "truenas_kerberos_realm" "example" {
  realm = "example-realm"
  kdc = []
  admin_server = []
  kpasswd_server = []
}
```

## Schema

### Required

- `realm` (Required) - Kerberos realm name. This is external to TrueNAS and is case-sensitive.     The general convention for kerberos realms is that they are upper-case.. Type: `string`

### Optional

- `primary_kdc` (Optional) - The master Kerberos domain controller for this realm. TrueNAS uses this as a fallback if it cannot get     credentials because of an invalid password. This can help in environments where the domain uses a hub-and-spoke     topology. Use this setting to reduce credential errors after TrueNAS automatically changes its machine password. . Type: `string`
- `kdc` (Optional) - List of kerberos domain controllers. If the list is empty then the kerberos     libraries will use DNS to look up KDCs. In some situations this is undesirable     as kerberos libraries are, for intance, not active directory site aware and so     may be suboptimal. Default: `[]`. Type: `array`
- `admin_server` (Optional) - List of kerberos admin servers. If the list is empty then the kerberos     libraries will use DNS to look them up. Default: `[]`. Type: `array`
- `kpasswd_server` (Optional) - List of kerberos kpasswd servers. If the list is empty then DNS will be used     to look them up if needed. Default: `[]`. Type: `array`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_kerberos_realm.example <id>
```
