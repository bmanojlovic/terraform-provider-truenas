# TrueNAS Keychain Credential Test

Tests the `truenas_keychaincredential` resource with SSH_KEY_PAIR type.

## Setup

Secrets are already copied from pool_dataset directory.

```bash
terraform init -upgrade
terraform plan
terraform apply
```

## What This Tests

- **SSH_KEY_PAIR** - Store SSH key pair in TrueNAS keychain
- **anyOf variant** - SSH_KEY_PAIR vs SSH_CREDENTIALS types
- **JSON attributes** - Complex object field handling

## Cleanup

```bash
terraform destroy
```

**Note:** The test uses a dummy SSH key. SSH_CREDENTIALS type requires an existing key ID and is not tested here.
