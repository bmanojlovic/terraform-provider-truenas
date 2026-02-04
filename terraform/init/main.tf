# Nested TrueNAS VM for safe provider testing
resource "truenas_vm" "test" {
  name        = "testtruenasvm"
  description = "Nested TrueNAS for provider testing (self-bootstrap)"
  vcpus       = 4
  memory      = 10240  # 10GB
  autostart   = false
}

# Devices reference VM but shouldn't block VM destruction
# VM deletion with force=true will clean up all devices automatically

# Boot disk - 32GB (power of 2 for block alignment)
resource "truenas_vm_device" "boot_disk" {
  vm = truenas_vm.test.id
  attributes = jsonencode({
    dtype        = "DISK"
    create_zvol  = true
    zvol_name    = "${var.pool_name}/vm-test-boot"
    zvol_volsize = 32 * 1024 * 1024 * 1024  # 32GB in bytes
    type         = "VIRTIO"
  })
  order = 1000
  
  lifecycle {
    # Don't try to destroy if VM is already gone
    ignore_changes = []
  }
}

# Data disks for ZFS pool testing inside nested TrueNAS
resource "truenas_vm_device" "data_disk" {
  count = 4
  vm    = truenas_vm.test.id
  attributes = jsonencode({
    dtype        = "DISK"
    create_zvol  = true
    zvol_name    = "${var.pool_name}/vm-test-data${count.index + 1}"
    zvol_volsize = 128 * 1024 * 1024 * 1024  # 128GB in bytes
    type         = "VIRTIO"
  })
  order = 1001 + count.index
}

# CD-ROM for TrueNAS ISO
resource "truenas_vm_device" "cdrom" {
  vm = truenas_vm.test.id
  attributes = jsonencode({
    dtype = "CDROM"
    path  = local.truenas_iso_path
  })
  order = 1005
}

# Network interface (bridged to LAN)
resource "truenas_vm_device" "nic" {
  vm = truenas_vm.test.id
  attributes = jsonencode({
    dtype      = "NIC"
    type       = "VIRTIO"
    nic_attach = var.bridge_interface
  })
  order = 1006
}

# Display device for console access (SPICE)
resource "truenas_vm_device" "display" {
  vm = truenas_vm.test.id
  attributes = jsonencode({
    dtype    = "DISPLAY"
    type     = "SPICE"
    password = "letmein"
  })
  order = 1007
}

# Start VM after all devices are configured
# VMs don't auto-start by default (cloud-like behavior)
resource "truenas_action_vm_start" "test" {
  id = tonumber(truenas_vm.test.id)
  depends_on = [
    truenas_vm_device.boot_disk,
    truenas_vm_device.data_disk,
    truenas_vm_device.cdrom,
    truenas_vm_device.nic,
    truenas_vm_device.display
  ]
}
