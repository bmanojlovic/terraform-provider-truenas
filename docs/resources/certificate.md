---
page_title: "truenas_certificate Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a new Certificate
---

# truenas_certificate (Resource)

Create a new Certificate

## Example Usage

```terraform
resource "truenas_certificate" "example" {
  create_type = "example-value"
  name = "example-value"
}
```

## Schema

### Required

- `create_type` (String) - Type of certificate creation operation. Valid values: `CERTIFICATE_CREATE_IMPORTED`, `CERTIFICATE_CREATE_CSR`, `CERTIFICATE_CREATE_IMPORTED_CSR`, `CERTIFICATE_CREATE_ACME`
- `name` (String) - Certificate name.

### Optional

- `CSR` (String) - PEM-encoded certificate signing request to import or `null`. Default: `None`
- `acme_directory_uri` (String) - ACME directory URI to be used for ACME certificate creation. Default: `None`
- `add_to_trusted_store` (Bool) - Whether to add this certificate to the trusted certificate store. Default: `False`
- `cert_extensions` (String) - Certificate extensions configuration.
- `certificate` (String) - PEM-encoded certificate to import or `null`. Default: `None`
- `city` (String) - City or locality name for certificate subject or `null`. Default: `None`
- `common` (String) - Common name for certificate subject or `null`. Default: `None`
- `country` (String) - Country name for certificate subject or `null`. Default: `None`
- `csr_id` (Int64) - CSR to be used for ACME certificate creation. Default: `None`
- `digest_algorithm` (String) - Hash algorithm for certificate signing. Default: `SHA256` Valid values: `SHA224`, `SHA256`, `SHA384`, `SHA512`
- `dns_mapping` (String) - A mapping of domain to ACME DNS Authenticator ID for each domain listed in SAN or common name of the CSR.
- `ec_curve` (String) - Elliptic curve to use for EC keys. Default: `SECP384R1` Valid values: `SECP256R1`, `SECP384R1`, `SECP521R1`, `ed25519`
- `email` (String) - Email address for certificate subject or `null`. Default: `None`
- `key_length` (Int64) - RSA key length in bits or `null`. Default: `None`
- `key_type` (String) - Type of cryptographic key to generate. Default: `RSA` Valid values: `RSA`, `EC`
- `organization` (String) - Organization name for certificate subject or `null`. Default: `None`
- `organizational_unit` (String) - Organizational unit for certificate subject or `null`. Default: `None`
- `passphrase` (String) - Passphrase to protect the private key or `null`. Default: `None`
- `privatekey` (String) - PEM-encoded private key to import or `null`. Default: `None`
- `renew_days` (Int64) - Number of days before the certificate expiration date to attempt certificate renewal. If certificate renewal     fails, renewal will be reattempted every day until expiration. Default: `10`
- `san` (List) - Subject alternative names for the certificate.
- `state` (String) - State or province name for certificate subject or `null`. Default: `None`
- `tos` (String) - Set this when creating an ACME certificate to accept terms of service of the ACME service. Default: `None`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_certificate.example <id>
```
