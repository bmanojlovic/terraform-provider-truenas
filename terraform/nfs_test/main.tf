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

variable "truenas_host" {
  type = string
}

variable "truenas_token" {
  type      = string
  sensitive = true
}

# List all NFS shares
data "truenas_sharing_nfss" "all" {}

output "nfs_shares" {
  value = data.truenas_sharing_nfss.all.items
}
