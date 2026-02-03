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

# Test interface.commit action (flat structure)
resource "truenas_interface_commit" "test" {
  checkin_timeout = 30
  rollback        = true
}

# Test interface.checkin action (depends on commit)
resource "truenas_interface_checkin" "test" {
  depends_on = [truenas_interface_commit.test]
}
