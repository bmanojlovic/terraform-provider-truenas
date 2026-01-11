---
page_title: "truenas_nvmet_port_subsys Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an association between a `port` and a subsystem (`subsys`).
---

# truenas_nvmet_port_subsys (Resource)

Create an association between a `port` and a subsystem (`subsys`).


## Example Usage

```terraform
resource "truenas_nvmet_port_subsys" "example" {
  port_id = 1
  subsys_id = 1
}
```

## Schema

### Required

- `port_id` (Int64) - ID of the NVMe-oF port to associate.
- `subsys_id` (Int64) - ID of the NVMe-oF subsystem to make accessible.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_port_subsys.example <id>
```
