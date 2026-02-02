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

variable "truenas_host" {
  type = string
}

variable "truenas_token" {
  type      = string
  sensitive = true
}

variable "truenas_pool" {
  type = string
}

resource "truenas_virt_volume_import_iso" "test_iso" {
  virt_volume_import_iso = jsonencode({
    name         = "test-netboot-iso"
    upload_iso   = true
    storage_pool = var.truenas_pool
  })
  file_content = filebase64("/tmp/tiny.iso")
}

output "import_result" {
  value = {
    id = truenas_virt_volume_import_iso.test_iso.id
  }
}
