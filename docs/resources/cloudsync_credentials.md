---
page_title: "truenas_cloudsync_credentials Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create Cloud Sync Credentials.
---

# truenas_cloudsync_credentials (Resource)

Create Cloud Sync Credentials.


## Example Usage

```terraform
resource "truenas_cloudsync_credentials" "example" {
  name = "example-value"
  provider = "example-value"
}
```

## Schema

### Required

- `name` (String) - Human-readable name for the cloud credential.

### Optional

- None

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_cloudsync_credentials.example <id>
```
