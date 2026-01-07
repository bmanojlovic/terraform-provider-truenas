---
page_title: "truenas_system_ntpserver Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new NTP server.
---

# truenas_system_ntpserver (Resource)

Configuration for creating a new NTP server.

## Example Usage

```terraform
resource "truenas_system_ntpserver" "example" {
  address = "example-address"
  burst = false
  iburst = true
  prefer = false
  minpoll = 6
  maxpoll = 10
  force = false
}
```

## Schema

### Required

- `address` (Required) - Hostname or IP address of the NTP server.. Type: `string`

### Optional

- `burst` (Optional) - Send a burst of packets when the server is reachable. Default: `False`. Type: `boolean`
- `iburst` (Optional) - Send a burst of packets when the server is unreachable. Default: `True`. Type: `boolean`
- `prefer` (Optional) - Mark this server as preferred for time synchronization. Default: `False`. Type: `boolean`
- `minpoll` (Optional) - Minimum polling interval (log2 seconds). Default: `6`. Type: `integer`
- `maxpoll` (Optional) - Maximum polling interval (log2 seconds). Default: `10`. Type: `integer`
- `force` (Optional) - Force creation even if the server is unreachable. Default: `False`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_system_ntpserver.example <id>
```
