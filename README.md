# TrueNAS Terraform Provider

A comprehensive Terraform provider for TrueNAS SCALE using native JSON-RPC 2.0 over WebSocket.

## Features

- ✅ **274 Generated Resources** covering the entire TrueNAS API surface
- ✅ **WebSocket JSON-RPC** for superior performance vs REST
- ✅ **Complete CRUD Operations** for all resources
- ✅ **Terraform Plugin Framework** implementation
- ✅ **Comprehensive Documentation** for all resources

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
  name        = "test-vm"
  description = "Test VM created by Terraform"
  vcpus       = 2
  memory      = 2147483648  # 2GB in bytes
  autostart   = false
}
```

## Available Resources

This provider includes **274 resources** covering:

- **Virtual Machines** (`truenas_vm`, `truenas_vm_device`)
- **Storage** (`truenas_pool`, `truenas_pool_dataset`, `truenas_pool_snapshot`)
- **Users & Groups** (`truenas_user`, `truenas_group`)
- **Sharing** (`truenas_sharing_nfs`, `truenas_sharing_smb`)
- **Network** (`truenas_interface`, `truenas_staticroute`)
- **Services** (`truenas_service`)
- **And many more...**

## Documentation

- [Provider Configuration](docs/index.md)
- [Resource Documentation](docs/resources/)
- [Examples](examples/)

## Development

Built with:
- Go 1.21+
- Terraform Plugin Framework
- WebSocket JSON-RPC 2.0
- OpenAPI-driven code generation

## License

This provider is licensed under the Mozilla Public License 2.0.
