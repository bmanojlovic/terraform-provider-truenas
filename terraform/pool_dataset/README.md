# TrueNAS Pool Dataset Test

Tests the `truenas_pool_dataset` resource with both FILESYSTEM and VOLUME types.

## Setup

1. Copy secrets from init directory:
```bash
cp ../init/secrets.auto.tfvars.json .
```

2. Initialize Terraform:
```bash
terraform init -upgrade
```

3. Review the plan:
```bash
terraform plan
```

4. Apply:
```bash
terraform apply
```

## What This Tests

- **FILESYSTEM dataset** - with quota, compression, atime settings
- **VOLUME dataset (zvol)** - with volsize, volblocksize, sparse
- **Nested datasets** - parent/child relationship

## Cleanup

```bash
terraform destroy
```

**Note:** This will create datasets in the pool specified by `pool_name` variable (default: "tank").
