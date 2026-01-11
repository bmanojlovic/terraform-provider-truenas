# TrueNAS Pool Snapshot Test

Tests the `truenas_pool_snapshot` resource with both explicit name and naming schema variants.

## Prerequisites

Requires existing datasets from pool_dataset test:
```bash
cd ../pool_dataset
terraform apply
```

## Setup

Secrets are already copied from pool_dataset directory.

```bash
terraform init -upgrade
terraform plan
terraform apply
```

## What This Tests

- **Explicit name** - snapshot with user-defined name
- **Naming schema** - snapshot with auto-generated name pattern
- **Recursive** - snapshot including child datasets

## Cleanup

```bash
terraform destroy
```

**Note:** Snapshots are created on datasets from the pool_dataset test.
