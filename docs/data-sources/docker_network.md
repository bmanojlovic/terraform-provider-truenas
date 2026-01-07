---
page_title: "truenas_docker_network Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Retrieves TrueNAS docker_network data
---

# truenas_docker_network (Data Source)

Retrieves TrueNAS docker_network data

## Example Usage

```terraform
data "truenas_docker_network" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the docker_network to retrieve.

### Read-Only

- `id` (string) - id value
- `options` (string) - options value
