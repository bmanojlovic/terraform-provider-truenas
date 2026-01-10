---
page_title: "truenas_staticroute Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a Static Route.
---

# truenas_staticroute (Resource)

Create a Static Route.

## Example Usage

```terraform
resource "truenas_staticroute" "example" {
  destination = "example-value"
  gateway = "example-value"
}
```

## Schema

### Required

- `destination` (String) - Destination network or host for this static route.
- `gateway` (String) - Gateway IP address for this static route.

### Optional

- `description` (String) - Optional description for this static route. Default: ``

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_staticroute.example <id>
```
