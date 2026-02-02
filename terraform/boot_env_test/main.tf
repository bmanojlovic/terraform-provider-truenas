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

# Clone the old boot environment
resource "truenas_action_boot_environment_clone" "backup" {
  boot_environment_clone = jsonencode({
    id     = "25.04.2.4"
    target = "25.04.2.4-terraform-clone"
  })
}

output "boot_envs" {
  value = data.truenas_boot_environments.all.items
}

output "clone_result" {
  value = {
    state  = truenas_action_boot_environment_clone.backup.state
    result = truenas_action_boot_environment_clone.backup.result
  }
}
