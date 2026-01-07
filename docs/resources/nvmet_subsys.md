---
page_title: "truenas_nvmet_subsys Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration data for creating a new NVMe-oF subsystem.
---

# truenas_nvmet_subsys (Resource)

Configuration data for creating a new NVMe-oF subsystem.

## Example Usage

```terraform
resource "truenas_nvmet_subsys" "example" {
  name = "example-name"
  allow_any_host = false
}
```

## Schema

### Required

- `name` (Required) - Human readable name for the subsystem.

If `subnqn` is not provided on creation, then this name will be appended to the `basenqn` from     `nvmet.global.config` to generate a subnqn.. Type: `string`

### Optional

- `subnqn` (Optional) - NVMe Qualified Name (NQN) for the subsystem.

Must be a valid NQN format if provided.. Type: `string`
- `allow_any_host` (Optional) - Any host can access the storage associated with this subsystem (i.e. no access control). Default: `False`. Type: `boolean`
- `pi_enable` (Optional) - Enable Protection Information (PI) for data integrity checking.. Type: `string`
- `qid_max` (Optional) - Maximum number of queue IDs allowed for this subsystem.. Type: `string`
- `ieee_oui` (Optional) - IEEE Organizationally Unique Identifier for the subsystem.. Type: `string`
- `ana` (Optional) - If set to either `True` or `False`, then *override* the global `ana` setting from `nvmet.global.config` for this     subsystem only.

If `null`, then the global `ana` setting will take effect.. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_subsys.example <id>
```
