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

# NFS config
resource "truenas_nfs_config" "test" {
  servers    = 4
  mountd_log = true
}

# SMB config
resource "truenas_smb_config" "test" {
  netbiosname = "TRUENAS"
  workgroup   = "WORKGROUP"
}

# FTP config
resource "truenas_ftp_config" "test" {
  port    = 21
  clients = 32
}

# SSH config
resource "truenas_ssh_config" "test" {
  tcpport                  = 22
  passwordauth             = true
  kerberosauth             = false
}

# SNMP config
resource "truenas_snmp_config" "test" {
  location = "Server Room"
  contact  = "admin@example.com"
}

# UPS config - requires driver/port, skip for now
# resource "truenas_ups_config" "test" {
#   mode        = "MASTER"
#   shutdown    = "BATT"
#   shutdowntimer = 30
# }

# Mail config - has oauth object field issue, skip for now
# resource "truenas_mail_config" "test" {
#   fromemail = "truenas@example.com"
#   outgoingserver = "smtp.example.com"
#   port      = 587
# }

# iSCSI global config
resource "truenas_iscsi_global_config" "test" {
  basename = "iqn.2005-10.org.freenas.ctl"
}

# Network config - test by adding VLAN and committing
resource "truenas_interface" "test_vlan" {
  type        = "VLAN"
  name        = "vlan99"
  vlan_tag    = 99
  vlan_parent_interface = "eno1"
  aliases = [
    jsonencode({
      type    = "INET"
      address = "10.99.0.1"
      netmask = 24
    })
  ]
}

resource "truenas_interface_commit" "apply" {
  depends_on = [truenas_interface.test_vlan]
}

# Check interface is up after commit
data "truenas_interface" "verify" {
  id = truenas_interface.test_vlan.id
  depends_on = [truenas_interface_commit.apply]
}

output "results" {
  value = {
    nfs_servers   = truenas_nfs_config.test.servers
    smb_workgroup = truenas_smb_config.test.workgroup
    ftp_port      = truenas_ftp_config.test.port
    ssh_port      = truenas_ssh_config.test.tcpport
    snmp_location = truenas_snmp_config.test.location
    vlan_name     = data.truenas_interface.verify.name
    vlan_tag      = data.truenas_interface.verify.vlan_tag
    vlan_parent   = data.truenas_interface.verify.vlan_parent_interface
  }
}
