---
page_title: "truenas_cloudsync_credentials Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Cloud credential configuration data for the new credential.
---

# truenas_cloudsync_credentials (Resource)

Cloud credential configuration data for the new credential.

## Example Usage

```terraform
resource "truenas_cloudsync_credentials" "example" {
  name = "example-name"
  provider = "example-provider"
}
```

## Schema

### Required

- `name` (Required) - Human-readable name for the cloud credential.. Type: `string`
- `provider` (Required) - Cloud provider configuration including type and authentication details.. Type: `string`

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cloudsync_credentials.example <id>
```
