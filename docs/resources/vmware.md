---
page_title: "truenas_vmware Resource - terraform-provider-truenas"
subcategory: "Other"
description: |-
  Create VMWare snapshot.
---

# truenas_vmware (Resource)

Create VMWare snapshot.


## Example Usage

```terraform
resource "truenas_vmware" "example" {
  datastore = "example"
  filesystem = "example"
  hostname = "example"
  password = "example"
  username = "example"
}
```

## Schema

### Required

- `datastore` (String) - Valid datastore name which exists on the VMWare host.
- `filesystem` (String) - ZFS filesystem or dataset to use for VMware storage.
- `hostname` (String) - Valid IP address / hostname of a VMWare host. When clustering, this is the vCenter server for the cluster.
- `password` (String) - Password for VMware host authentication.
- `username` (String) - Credentials used to authorize access to the VMWare host.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_vmware.example <id>
```
