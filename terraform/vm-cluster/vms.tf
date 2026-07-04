# --- Create VMs ---

resource "truenas_vm" "node" {
  for_each   = { for n in local.nodes : n.name => n }
  depends_on = [truenas_pool_snapshot_clone.node]

  name        = each.key
  description = "${each.value.role} - ${each.key}.${var.domain}"
  vcpus       = each.value.role == "master" ? var.master_vcpus : var.worker_vcpus
  memory      = each.value.role == "master" ? var.master_memory_mb : var.worker_memory_mb
  cpu_mode    = "HOST-PASSTHROUGH"
  autostart   = true
}

resource "truenas_vm_device" "disk" {
  for_each = { for n in local.nodes : n.name => n }
  vm       = truenas_vm.node[each.key].id
  attributes = jsonencode({
    dtype = "DISK"
    path  = "/dev/zvol/${var.pool}/${each.key}"
    type  = "VIRTIO"
  })
  order = 1000
}

resource "truenas_vm_device" "cdrom" {
  for_each = { for n in local.nodes : n.name => n }
  vm       = truenas_vm.node[each.key].id
  attributes = jsonencode({
    dtype = "CDROM"
    path  = "/mnt/${var.iso_pool}/iso-store/cloud-init-${each.key}.iso"
  })
  order = 1005
}

resource "truenas_vm_device" "nic" {
  for_each = { for n in local.nodes : n.name => n }
  vm       = truenas_vm.node[each.key].id
  attributes = jsonencode({
    dtype      = "NIC"
    type       = "VIRTIO"
    nic_attach = var.nic_attach
  })
  order = 1006
}

resource "truenas_vm_device" "display" {
  for_each = { for n in local.nodes : n.name => n }
  vm       = truenas_vm.node[each.key].id
  attributes = jsonencode({
    dtype    = "DISPLAY"
    type     = "SPICE"
    web      = true
    password = "letmein"
    port     = 6900 + each.value.idx
    web_port = 7000 + each.value.idx
  })
  order = 1007
}

# --- Start VMs ---

resource "truenas_action_vm_start" "node" {
  for_each   = { for n in local.nodes : n.name => n }
  depends_on = [truenas_vm_device.disk, truenas_vm_device.cdrom, truenas_vm_device.nic, truenas_vm_device.display]
  id         = truenas_vm.node[each.key].id
}
