---
page_title: "truenas_kerberos_keytab Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a kerberos keytab. Uploaded keytab files will be merged with the system
---

# truenas_kerberos_keytab (Resource)

Create a kerberos keytab. Uploaded keytab files will be merged with the system

## Example Usage

```terraform
resource "truenas_kerberos_keytab" "example" {
  file = "example-value"
  name = "example-value"
}
```

## Schema

### Required

- `file` (String) - Base64 encoded kerberos keytab entries to append to the system keytab.
- `name` (String) - Name of the kerberos keytab entry. This is an identifier for the keytab and not     the name of the keytab file. Some names are used for internal purposes such     as AD_MACHINE_ACCOUNT and IPA_MACHIN

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_kerberos_keytab.example <id>
```
