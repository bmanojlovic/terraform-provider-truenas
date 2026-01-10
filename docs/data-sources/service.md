---
page_title: "truenas_service Data Source - terraform-provider-truenas"
subcategory: ""
description: |-
  Returns instance matching `id`. If `id` is not found, Validation error is raised.
---

# truenas_service (Data Source)

Returns instance matching `id`. If `id` is not found, Validation error is raised.

## Example Usage

```terraform
data "truenas_service" "example" {
  id = "1"
}
```

## Schema

### Required

- `id` (String) The ID of the service to retrieve.

### Read-Only

- `enable` (Bool) - Whether the service is enabled to start on boot.
- `pids` (List) - Array of process IDs associated with this service.
- `service` (String) - Name of the system service.
- `state` (String) - Current state of the service (e.g., 'RUNNING', 'STOPPED').
