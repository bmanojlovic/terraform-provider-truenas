---
page_title: "truenas_docker_networks Data Source - terraform-provider-truenas"
subcategory: "Applications"
description: |-
  Query all docker networks
---

# truenas_docker_networks (Data Source)

Query all docker networks

## Example Usage

```terraform
data "truenas_docker_networks" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the docker_networks to retrieve.

### Read-Only

- None
