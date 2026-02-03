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

# Configure NFS service - change servers
resource "truenas_nfs_config" "main" {
  servers      = 2
  mountd_log   = true
}

output "nfs_servers" {
  value = truenas_nfs_config.main.servers
}
