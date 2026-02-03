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

resource "truenas_smb_config" "main" {
  netbiosname = "TRUENAS"
  workgroup   = "WORKGROUP"
}
