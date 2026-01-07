---
page_title: "truenas_nvmet_host Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  NVMe-oF host configuration data for creation.
---

# truenas_nvmet_host (Resource)

NVMe-oF host configuration data for creation.

## Example Usage

```terraform
resource "truenas_nvmet_host" "example" {
  hostnqn = "example-hostnqn"
  dhchap_hash = "SHA-256"
}
```

## Schema

### Required

- `hostnqn` (Required) - hostnqn configuration. Type: `string`

### Optional

- `dhchap_key` (Optional) - If set, the secret that the host must present when connecting.

A suitable secret can be generated using `nvme gen-dhchap-key`, or by using the `nvmet.host.generate_key` API.. Type: `string`
- `dhchap_ctrl_key` (Optional) - If set, the secret that this TrueNAS will present to the host when the host is connecting (Bi-Directional     Authentication).

A suitable secret can be generated using `nvme gen-dhchap-key`, or by using the `nvmet.host.generate_key` API.. Type: `string`
- `dhchap_dhgroup` (Optional) - If selected, the DH (Diffie-Hellman) key exchange built on top of CHAP to be used for authentication.. Type: `string`
- `dhchap_hash` (Optional) - HMAC (Hashed Message Authentication Code) to be used in conjunction if a `dhchap_dhgroup` is selected. Valid values: `SHA-256`, `SHA-384`, `SHA-512` Default: `SHA-256`. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_nvmet_host.example <id>
```
