---
page_title: "truenas_kerberos_keytab Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Kerberos keytab configuration data for creation.
---

# truenas_kerberos_keytab (Resource)

Kerberos keytab configuration data for creation.

## Example Usage

```terraform
resource "truenas_kerberos_keytab" "example" {
  name = "example-name"
  file = "example-file"
}
```

## Schema

### Required

- `name` (Required) - Name of the kerberos keytab entry. This is an identifier for the keytab and not     the name of the keytab file. Some names are used for internal purposes such     as AD_MACHINE_ACCOUNT and IPA_MACHINE_ACCOUNT.. Type: `string`
- `file` (Required) - Base64 encoded kerberos keytab entries to append to the system keytab. . Type: `string`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_kerberos_keytab.example <id>
```
