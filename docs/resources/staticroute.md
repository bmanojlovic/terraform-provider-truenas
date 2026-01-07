---
page_title: "truenas_staticroute Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for the new static route.
---

# truenas_staticroute (Resource)

Configuration for the new static route.

## Example Usage

```terraform
resource "truenas_staticroute" "example" {
  destination = "example-destination"
  gateway = "example-gateway"
  description = ""
}
```

## Schema

### Required

- `destination` (Required) - Destination network or host for this static route.. Type: `string`
- `gateway` (Required) - Gateway IP address for this static route.. Type: `string`

### Optional

- `description` (Optional) - Optional description for this static route. Default: ``. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_staticroute.example <id>
```
