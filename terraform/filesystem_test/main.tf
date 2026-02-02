terraform {
  required_providers {
    truenas = {
      source = "bmanojlovic/truenas"
    }
  }
}

provider "truenas" {
  host  = var.truenas_host
  token = var.truenas_token
}

# Test filesystem.put - upload cloud-init config
resource "truenas_filesystem_put" "cloud_init" {
  path         = "/mnt/${var.pool_name}/test-upload.txt"
  file_content = base64encode("Hello from Terraform!\nThis is a test file.\n")
}

output "file_path" {
  value = truenas_filesystem_put.cloud_init.id
}
