output "nodes" {
  value = { for n in local.nodes : n.name => {
    role = n.role
    ip   = n.ip
  }}
}

output "master_ip" {
  value = var.master_ip
}

output "ssh_command" {
  value = "ssh root@${var.master_ip}"
}

output "display" {
  value = { for name, dev in truenas_vm_device.display : name => {
    port     = jsondecode(dev.attributes).port
    web_port = jsondecode(dev.attributes).web_port
  }}
}
