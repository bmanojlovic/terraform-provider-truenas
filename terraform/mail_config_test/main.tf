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

resource "truenas_mail_config" "test" {
  fromemail      = "truenas@example.com"
  outgoingserver = "smtp.example.com"
  port           = 587
}

output "mail_from" {
  value = truenas_mail_config.test.fromemail
}
