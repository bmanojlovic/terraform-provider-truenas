# File Upload Support

## Overview

Added multipart file upload support to the TrueNAS Terraform provider to enable uploading files via the `filesystem.put` API endpoint.

## Changes Made

### 1. Client Enhancement (`internal/client/client.go`)

- Added HTTP client alongside WebSocket client
- Implemented `UploadFile()` method for multipart form-data uploads
- Supports JSON metadata + binary file content

### 2. New Resource (`internal/provider/resource_filesystem_put.go`)

- `truenas_filesystem_put` resource for file uploads
- Accepts base64-encoded content
- Optional file mode/permissions

### 3. Generator Update (`generate.py`)

- Added `ACTION_FILE_UPLOAD_TEMPLATE` template
- Registered `filesystem_put` in resources list
- Auto-updates provider.go with new resource

## Usage Example

### Upload Cloud-Init ISO

```hcl
# Generate ISO locally
data "external" "cloud_init_iso" {
  program = ["bash", "-c", <<-EOT
    genisoimage -o /tmp/cloud-init.iso -V cidata -J -r /path/to/cloud-init/
    base64 -w 0 /tmp/cloud-init.iso > /tmp/cloud-init.iso.b64
    echo "{\"content\": \"$(cat /tmp/cloud-init.iso.b64)\"}"
  EOT
  ]
}

# Upload to TrueNAS
resource "truenas_filesystem_put" "cloud_init_iso" {
  path    = "/mnt/tank/isos/cloud-init.iso"
  content = data.external.cloud_init_iso.result.content
  mode    = 0644
}
```

### Upload Text File

```hcl
resource "truenas_filesystem_put" "config" {
  path    = "/mnt/tank/configs/app.conf"
  content = base64encode("key=value\nfoo=bar")
  mode    = 0600
}
```

## Technical Details

### Multipart Upload Format

The `filesystem.put` API expects:
- **Part 1 (`data`)**: JSON with `{"path": "/destination/path", "mode": 0644}`
- **Part 2 (`file`)**: Binary file content

### Base64 Encoding

Files must be base64-encoded in Terraform because:
- Terraform strings are UTF-8 text
- Binary data needs encoding for safe transport
- Use `base64encode()` function or external data source

## Complete Example

See `examples/cloud-init-iso/main.tf` for a full working example that:
1. Generates cloud-init user-data and meta-data
2. Creates ISO using `genisoimage`
3. Uploads to TrueNAS
4. Attaches to VM as CDROM device
