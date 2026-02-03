terraform {
  required_providers {
    truenas = {
      source  = "bmanojlovic/truenas"
      version = "0.1.0"
    }
  }
}

provider "truenas" {
  host  = var.truenas_host
  token = var.truenas_token
}

variable "truenas_host" { type = string }
variable "truenas_token" { 
  type = string
  sensitive = true 
}

# Change path - should force recreation
resource "truenas_filesystem_mkdir" "test" {
  path         = "/mnt/fast/terraform-update-test3"
  options_mode = "700"
}

output "result" {
  value = truenas_filesystem_mkdir.test.result
}
