---
page_title: "truenas_nvmet_port_subsys Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Port-subsystem association configuration data for creation.
---

# truenas_nvmet_port_subsys (Resource)

Port-subsystem association configuration data for creation.

## Example Usage

```terraform
resource "truenas_nvmet_port_subsys" "example" {
  port_id = 1
  subsys_id = 1
}
```

## Schema

### Required

- `port_id` (Required) - ID of the NVMe-oF port to associate.. Type: `integer`
- `subsys_id` (Required) - ID of the NVMe-oF subsystem to make accessible.. Type: `integer`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_port_subsys.example <id>
```
