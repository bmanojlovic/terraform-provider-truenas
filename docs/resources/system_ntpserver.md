---
page_title: "truenas_system_ntpserver Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Add an NTP Server.
---

# truenas_system_ntpserver (Resource)

Add an NTP Server.


## Example Usage

```terraform
resource "truenas_system_ntpserver" "example" {
  address = "example-value"
}
```

## Schema

### Required

- `address` (String) - Hostname or IP address of the NTP server.

### Optional

- `burst` (Bool) - Send a burst of packets when the server is reachable. Default: `False`
- `force` (Bool) - Force creation even if the server is unreachable. Default: `False`
- `iburst` (Bool) - Send a burst of packets when the server is unreachable. Default: `True`
- `maxpoll` (Int64) - Maximum polling interval (log2 seconds). Default: `10`
- `minpoll` (Int64) - Minimum polling interval (log2 seconds). Default: `6`
- `prefer` (Bool) - Mark this server as preferred for time synchronization. Default: `False`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_system_ntpserver.example <id>
```
