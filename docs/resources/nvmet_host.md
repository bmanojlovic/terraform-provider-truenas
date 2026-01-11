---
page_title: "truenas_nvmet_host Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an NVMe target `host`.
---

# truenas_nvmet_host (Resource)

Create an NVMe target `host`.


## Example Usage

```terraform
resource "truenas_nvmet_host" "example" {
  hostnqn = "example-value"
}
```

## Schema

### Required

- `hostnqn` (String) - 

### Optional

- `dhchap_ctrl_key` (String) - If set, the secret that this TrueNAS will present to the host when the host is connecting (Bi-Directional     Authentication).  A suitable secret can be generated using `nvme gen-dhchap-key`, or by us Default: `None`
- `dhchap_dhgroup` (String) - If selected, the DH (Diffie-Hellman) key exchange built on top of CHAP to be used for authentication. Default: `None`
- `dhchap_hash` (String) - HMAC (Hashed Message Authentication Code) to be used in conjunction if a `dhchap_dhgroup` is selected. Default: `SHA-256` Valid values: `SHA-256`, `SHA-384`, `SHA-512`
- `dhchap_key` (String) - If set, the secret that the host must present when connecting.  A suitable secret can be generated using `nvme gen-dhchap-key`, or by using the `nvmet.host.generate_key` API. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_host.example <id>
```
