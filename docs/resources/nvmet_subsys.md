---
page_title: "truenas_nvmet_subsys Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a NVMe target subsystem (`subsys`).
---

# truenas_nvmet_subsys (Resource)

Create a NVMe target subsystem (`subsys`).


## Example Usage

```terraform
resource "truenas_nvmet_subsys" "example" {
  name = "example-value"
}
```

## Schema

### Required

- `name` (String) - Human readable name for the subsystem.  If `subnqn` is not provided on creation, then this name will be appended to the `basenqn` from     `nvmet.global.config` to generate a subnqn.

### Optional

- `allow_any_host` (Bool) - Any host can access the storage associated with this subsystem (i.e. no access control). Default: `False`
- `ana` (String) - If set to either `True` or `False`, then *override* the global `ana` setting from `nvmet.global.config` for this     subsystem only.  If `null`, then the global `ana` setting will take effect. Default: `None`
- `ieee_oui` (String) - IEEE Organizationally Unique Identifier for the subsystem. Default: `None`
- `pi_enable` (String) - Enable Protection Information (PI) for data integrity checking. Default: `None`
- `qid_max` (Int64) - Maximum number of queue IDs allowed for this subsystem. Default: `None`
- `subnqn` (String) - NVMe Qualified Name (NQN) for the subsystem.  Must be a valid NQN format if provided. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_subsys.example <id>
```
