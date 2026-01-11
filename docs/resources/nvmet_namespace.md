---
page_title: "truenas_nvmet_namespace Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a NVMe target namespace in a subsystem (`subsys`).
---

# truenas_nvmet_namespace (Resource)

Create a NVMe target namespace in a subsystem (`subsys`).


## Example Usage

```terraform
resource "truenas_nvmet_namespace" "example" {
  device_path = "example-value"
  device_type = "example-value"
  subsys_id = 1
}
```

## Schema

### Required

- `device_path` (String) - Normalized path to the device or file for the namespace.
- `device_type` (String) - Type of device (or file) used to implement the namespace. Valid values: `ZVOL`, `FILE`
- `subsys_id` (Int64) - ID of the NVMe-oF subsystem to contain this namespace.

### Optional

- `enabled` (Bool) - If `enabled` is `False` then the namespace will not be accessible.  Some namespace configuration changes are blocked when that namespace is enabled. Default: `True`
- `filesize` (Int64) - When `device_type` is "FILE" then this will be the size of the file in bytes. Default: `None`
- `nsid` (Int64) - Namespace ID (NSID).  Each namespace within a subsystem has an associated NSID, unique within that subsystem.  If not supplied during `namespace` creation then the next available NSID will be used. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_namespace.example <id>
```
