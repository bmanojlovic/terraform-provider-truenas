# TrueNAS VM Bootstrap - Self-Contained Test Environment

This creates a **nested TrueNAS VM** on your production TrueNAS for safe provider testing.

## Concept: Self-Bootstrap Testing

Instead of testing the provider against your production TrueNAS (risky!), this creates a disposable TrueNAS VM where you can:
- Test all provider features safely
- Break things without consequences
- Reset/recreate anytime
- Develop and validate changes

## VM Specifications

- **Name:** test-truenas-vm
- **CPU:** 4 vCPUs
- **Memory:** 10GB
- **Boot Disk:** 30GB
- **Data Disks:** 4x 128GB (for ZFS pool testing inside the VM)

## Setup Steps

### 1. Bootstrap the VM (on production TrueNAS)

Create `secrets.auto.tfvars.json`:
```json
{
  "truenas_host": "192.168.1.100",
  "truenas_token": "your-production-api-token"
}
```

Optional overrides in `terraform.tfvars`:
```hcl
pool_name        = "tank"
truenas_iso_path = "/mnt/tank/iso/TrueNAS-SCALE.iso"
bridge_interface = "br0"
```

Apply:
```bash
terraform init
terraform apply
```

### 2. Install TrueNAS in the VM

1. Start VM from TrueNAS UI
2. Boot from ISO, install TrueNAS
3. Configure network (bridge to your LAN)
4. Set root password
5. Access web UI at VM's IP

### 3. Configure Test Environment

Inside the nested TrueNAS:
1. Create ZFS pool with the 4x 128GB disks
2. Enable API access
3. Generate API token: System Settings → API Keys
4. Note the VM's IP address

### 4. Test Provider Against Nested TrueNAS

Now point your provider tests at the **nested** TrueNAS:

```bash
export TF_VAR_truenas_host="192.168.1.200"  # VM IP
export TF_VAR_truenas_token="nested-vm-token"

cd ../your_test/
terraform apply  # Safe to break things!
```

## Benefits

- **Isolation:** Production TrueNAS stays untouched
- **Reproducible:** Destroy and recreate test VM anytime
- **Realistic:** Full TrueNAS environment for testing
- **Safe:** Test destructive operations without fear

## Cleanup

When done testing:
```bash
terraform destroy  # Removes VM and all test data
```
