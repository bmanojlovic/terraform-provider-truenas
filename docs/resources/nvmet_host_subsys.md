---
page_title: "truenas_nvmet_host_subsys Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an association between a `host` and a subsystem (`subsys`).
---

# truenas_nvmet_host_subsys (Resource)

Create an association between a `host` and a subsystem (`subsys`).

## Example Usage

```terraform
resource "truenas_nvmet_host_subsys" "example" {
  host_id = 1
  subsys_id = 1
}
```

## Schema

### Required

- `host_id` (Int64) - ID of the NVMe-oF host to authorize.
- `subsys_id` (Int64) - ID of the NVMe-oF subsystem to grant access to.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_host_subsys.example <id>
```
