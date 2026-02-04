# TrueNAS Terraform Provider

A comprehensive Terraform provider for TrueNAS SCALE using native JSON-RPC 2.0 over WebSocket.

## Features

- ✅ **Comprehensive API Coverage** - Auto-generated from TrueNAS JSON-RPC API schema
- ✅ **Native WebSocket JSON-RPC** - Direct protocol support for optimal performance
- ✅ **Complete CRUD Operations** for resources
- ✅ **Terraform Plugin Framework** implementation
- ✅ **Generated Documentation** for discovered resources and data sources

## Requirements

- **TrueNAS SCALE**: Version 25.10.1 (Goldeye) or later
- **Terraform**: 0.13+

## Quick Start

### Installation

```hcl
terraform {
  required_providers {
    truenas = {
      source = "bmanojlovic/truenas"
    }
  }
}
```

### Configuration

```hcl
provider "truenas" {
  host  = "192.168.1.100"
  token = "your-api-token"
}
```

### Basic Usage

```hcl
resource "truenas_vm" "example" {
  name        = "testvm"
  description = "Test VM created by Terraform"
  vcpus       = 2
  memory      = 2048  # 2GB in megabytes
  autostart   = false
}

# Attach devices while VM is stopped
resource "truenas_vm_device" "disk" {
  vm = truenas_vm.example.id
  attributes = jsonencode({
    dtype        = "DISK"
    create_zvol  = true
    zvol_name    = "pool/vm-disk"
    zvol_volsize = 10 * 1024 * 1024 * 1024  # 10GB in bytes
    type         = "VIRTIO"
  })
  order = 1000
}

# Start VM after all devices are attached
resource "truenas_action_vm_start" "example" {
  id = tonumber(truenas_vm.example.id)
  depends_on = [truenas_vm_device.disk]
}
```

**Note**: VMs don't auto-start by default (`start_on_create = false`). TrueNAS requires VMs to be stopped when attaching devices, so the recommended workflow is: create VM → attach devices → start with `truenas_action_vm_start`. Setting `start_on_create = true` will start the VM immediately, but devices cannot be attached until it's stopped.

## Available Resources

This provider covers:

- **Virtual Machines** (`truenas_vm`, `truenas_vm_device`)
- **Storage** (`truenas_pool`, `truenas_pool_dataset`, `truenas_pool_snapshot`)
- **Users & Groups** (`truenas_user`, `truenas_group`)
- **Sharing** (`truenas_sharing_nfs`, `truenas_sharing_smb`)
- **Network** (`truenas_interface`, `truenas_staticroute`)
- **Services** (`truenas_service`)
- **And many more...**

See [Resource Documentation](docs/resources/) for the complete list.

## Documentation

- [Provider Configuration](docs/index.md)
- [Resource Documentation](docs/resources/)
- [Integration Tests](terraform/) - Real-world usage examples

## Development

Built with:
- Go 1.21+
- Terraform Plugin Framework
- WebSocket JSON-RPC 2.0
- TrueNAS JSON-RPC schema-driven code generation

## Testing

This provider uses a two-tier testing strategy:

### Unit Tests (CI/CD)

**Purpose:** Validate the code generator logic, not the TrueNAS API itself.

Unit tests run automatically in CI and test core logic without requiring a TrueNAS instance:
- Schema validation (generator creates valid Terraform schemas)
- Parameter building (optional fields are omitted when null)
- Optional field handling (generator correctly checks IsNull())
- Business logic (start_on_create defaults, ID conversion)

**Important:** These tests validate the **generator's behavior**. If TrueNAS changes its API schema, resources are regenerated from the JSON-RPC spec. Tests ensure the generator produces correct code patterns.

```bash
# Run unit tests locally
go test -v ./internal/...
```

### Acceptance Tests (Local Only)

Acceptance tests require a real TrueNAS instance and are skipped in CI:
- Full CRUD lifecycle testing
- Real API integration
- Resource state management

**Prerequisites:**
- TrueNAS SCALE instance (accessible via network)
- API token with appropriate permissions

**Running acceptance tests:**

```bash
# Set environment variables
export TRUENAS_HOST=192.168.1.100
export TRUENAS_TOKEN=your-api-token

# Run acceptance tests
./test-local.sh
```

Or run directly:

```bash
TRUENAS_HOST=192.168.1.100 TRUENAS_TOKEN=your-token TF_ACC=1 go test ./internal/provider -v -run TestAcc
```

**Note:** Acceptance tests create and destroy real resources on your TrueNAS instance. Use a test environment.

## License

This provider is licensed under the Mozilla Public License 2.0.
