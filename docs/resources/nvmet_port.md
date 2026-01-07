---
page_title: "truenas_nvmet_port Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  NVMe-oF port configuration data for creation (TCP/RDMA or Fibre Channel).
---

# truenas_nvmet_port (Resource)

NVMe-oF port configuration data for creation (TCP/RDMA or Fibre Channel).

## Example Usage

```terraform
resource "truenas_nvmet_port" "example" {
  # Configuration here
}
```

## Schema

### Required

- None

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_port.example <id>
```
