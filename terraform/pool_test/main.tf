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

variable "truenas_pool" {
  type = string
}

# Get all pools and find ours
data "truenas_pools" "all" {}

locals {
  pool = [for p in data.truenas_pools.all.items : p if p.name == var.truenas_pool][0]
}

# Test pool.upgrade action (flat structure)
resource "truenas_pool_upgrade" "test" {
  id = local.pool.id
}

output "upgrade_result" {
  value = {
    pool_name = local.pool.name
    pool_id   = local.pool.id
    state     = truenas_pool_upgrade.test.state
    result    = truenas_pool_upgrade.test.result
  }
}
