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

data "truenas_boot_environments" "all" {}

# Clone the boot environment (flat structure - no jsonencode needed)
resource "truenas_boot_environment_clone" "backup" {
  id     = "25.04.2.4"
  target = "25.04.2.4-terraform-clone"
}

output "boot_envs" {
  value = data.truenas_boot_environments.all.items
}

output "clone_result" {
  value = {
    state  = truenas_boot_environment_clone.backup.state
    result = truenas_boot_environment_clone.backup.result
  }
}
