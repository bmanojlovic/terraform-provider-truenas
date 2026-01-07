---
page_title: "truenas_certificate Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  CertificateCreateArgs parameters.
---

# truenas_certificate (Resource)

CertificateCreateArgs parameters.

## Example Usage

```terraform
resource "truenas_certificate" "example" {
  name = "example-name"
  create_type = "CERTIFICATE_CREATE_IMPORTED"
  add_to_trusted_store = false
  key_type = "RSA"
  ec_curve = "SECP256R1"
  digest_algorithm = "SHA224"
  renew_days = 10
}
```

## Schema

### Required

- `name` (Required) - Certificate name.. Type: `string`
- `create_type` (Required) - Type of certificate creation operation. Valid values: `CERTIFICATE_CREATE_IMPORTED`, `CERTIFICATE_CREATE_CSR`, `CERTIFICATE_CREATE_IMPORTED_CSR`. Type: `string`

### Optional

- `add_to_trusted_store` (Optional) - Whether to add this certificate to the trusted certificate store. Default: `False`. Type: `boolean`
- `certificate` (Optional) - PEM-encoded certificate to import or `null`.. Type: `string`
- `privatekey` (Optional) - PEM-encoded private key to import or `null`.. Type: `string`
- `CSR` (Optional) - PEM-encoded certificate signing request to import or `null`.. Type: `string`
- `key_length` (Optional) - RSA key length in bits or `null`.. Type: `string`
- `key_type` (Optional) - Type of cryptographic key to generate. Valid values: `RSA`, `EC` Default: `RSA`. Type: `string`
- `ec_curve` (Optional) - Elliptic curve to use for EC keys. Valid values: `SECP256R1`, `SECP384R1`, `SECP521R1` Default: `SECP384R1`. Type: `string`
- `passphrase` (Optional) - Passphrase to protect the private key or `null`.. Type: `string`
- `city` (Optional) - City or locality name for certificate subject or `null`.. Type: `string`
- `common` (Optional) - Common name for certificate subject or `null`.. Type: `string`
- `country` (Optional) - Country name for certificate subject or `null`.. Type: `string`
- `email` (Optional) - Email address for certificate subject or `null`.. Type: `string`
- `organization` (Optional) - Organization name for certificate subject or `null`.. Type: `string`
- `organizational_unit` (Optional) - Organizational unit for certificate subject or `null`.. Type: `string`
- `state` (Optional) - State or province name for certificate subject or `null`.. Type: `string`
- `digest_algorithm` (Optional) - Hash algorithm for certificate signing. Valid values: `SHA224`, `SHA256`, `SHA384` Default: `SHA256`. Type: `string`
- `san` (Optional) - Subject alternative names for the certificate.. Type: `array`
- `cert_extensions` (Optional) - Certificate extensions configuration.. Type: `object`
- `acme_directory_uri` (Optional) - ACME directory URI to be used for ACME certificate creation.. Type: `string`
- `csr_id` (Optional) - CSR to be used for ACME certificate creation.. Type: `string`
- `tos` (Optional) - Set this when creating an ACME certificate to accept terms of service of the ACME service.. Type: `string`
- `dns_mapping` (Optional) - A mapping of domain to ACME DNS Authenticator ID for each domain listed in SAN or common name of the CSR.. Type: `object`
- `renew_days` (Optional) - Number of days before the certificate expiration date to attempt certificate renewal. If certificate renewal     fails, renewal will be reattempted every day until expiration. Default: `10`. Type: `integer`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_certificate.example <id>
```
