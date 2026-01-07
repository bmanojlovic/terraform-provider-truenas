---
page_title: "truenas_vmware Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new VMware integration.
---

# truenas_vmware (Resource)

Configuration for creating a new VMware integration.

## Example Usage

```terraform
resource "truenas_vmware" "example" {
  datastore = "example-datastore"
  filesystem = "example-filesystem"
  hostname = "example-hostname"
  username = "example-username"
  password = "example-password"
}
```

## Schema

### Required

- `datastore` (Required) - Valid datastore name which exists on the VMWare host.. Type: `string`
- `filesystem` (Required) - ZFS filesystem or dataset to use for VMware storage.. Type: `string`
- `hostname` (Required) - Valid IP address / hostname of a VMWare host. When clustering, this is the vCenter server for the cluster.. Type: `string`
- `username` (Required) - Credentials used to authorize access to the VMWare host.. Type: `string`
- `password` (Required) - Password for VMware host authentication.. Type: `string`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_vmware.example <id>
```
