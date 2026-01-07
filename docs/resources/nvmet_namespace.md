---
page_title: "truenas_nvmet_namespace Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  NVMe-oF namespace configuration data for creation.
---

# truenas_nvmet_namespace (Resource)

NVMe-oF namespace configuration data for creation.

## Example Usage

```terraform
resource "truenas_nvmet_namespace" "example" {
  device_type = "ZVOL"
  device_path = "example-device_path"
  subsys_id = 1
  enabled = true
}
```

## Schema

### Required

- `device_type` (Required) - Type of device (or file) used to implement the namespace.  Valid values: `ZVOL`, `FILE`. Type: `string`
- `device_path` (Required) - Normalized path to the device or file for the namespace.. Type: `string`
- `subsys_id` (Required) - ID of the NVMe-oF subsystem to contain this namespace.. Type: `integer`

### Optional

- `nsid` (Optional) - Namespace ID (NSID).

Each namespace within a subsystem has an associated NSID, unique within that subsystem.

If not supplied during `namespace` creation then the next available NSID will be used.. Type: `string`
- `filesize` (Optional) - When `device_type` is "FILE" then this will be the size of the file in bytes.. Type: `string`
- `enabled` (Optional) - If `enabled` is `False` then the namespace will not be accessible.

Some namespace configuration changes are blocked when that namespace is enabled. Default: `True`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_namespace.example <id>
```
