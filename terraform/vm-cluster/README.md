# VM Cluster from Cloud Image

Creates a multi-node VM cluster on TrueNAS SCALE from a cloud image (qcow2).

## What it does

1. Fetches a cloud image from URL using a container on TrueNAS
2. Imports it into a template zvol
3. Clones the template for each node
4. Generates per-node cloud-init ISOs (hostname, static IP, SSH key)
5. Creates VMs with VIRTIO disk, NIC, and SPICE display
6. Starts all VMs

## Usage

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
terraform init
terraform apply
```

## Requirements

- TrueNAS SCALE 25.04+
- `mkisofs` or `genisoimage` installed locally (for cloud-init ISO generation)
- SSH public key at `~/.ssh/id_ed25519.pub`

## Credentials

Set via environment variables:
```bash
export TF_VAR_truenas_host=192.168.1.100
export TF_VAR_truenas_token=your-api-token
```

## Customization

- Change `node_count` for more/fewer nodes
- Change `image_url` to use a different OS cloud image
- Adjust `master_*` and `worker_*` for different VM specs
