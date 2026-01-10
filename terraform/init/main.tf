# TrueNAS test VM
resource "truenas_vm" "test_truenas" {
  name        = "testvm5"
  description = "TrueNAS VM for provider testing"
  vcpus       = 4
  memory      = 10240  # 10GB
  autostart   = false
}

# Data source to verify VM was created
data "truenas_vm" "test_lookup" {
  id = truenas_vm.test_truenas.id
}

# Query data source to list all VMs
data "truenas_vms" "all" {}

# Boot disk - 30GB
resource "truenas_vm_device" "boot_disk" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "DISK"
    create_zvol = true
    zvol_name = "${var.pool_name}/vm-test5-boot"
    zvol_volsize = 32212254720  # 30GB in bytes
    type = "VIRTIO"
  })
  order = 1000
}

# Data disk 1 - 128GB
resource "truenas_vm_device" "data_disk_1" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "DISK"
    create_zvol = true
    zvol_name = "${var.pool_name}/vm-test5-data1"
    zvol_volsize = 137438953472  # 128GB in bytes
    type = "VIRTIO"
  })
  order = 1001
}

# Data disk 2 - 128GB
resource "truenas_vm_device" "data_disk_2" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "DISK"
    create_zvol = true
    zvol_name = "${var.pool_name}/vm-test5-data2"
    zvol_volsize = 137438953472
    type = "VIRTIO"
  })
  order = 1002
}

# Data disk 3 - 128GB
resource "truenas_vm_device" "data_disk_3" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "DISK"
    create_zvol = true
    zvol_name = "${var.pool_name}/vm-test5-data3"
    zvol_volsize = 137438953472
    type = "VIRTIO"
  })
  order = 1003
}

# Data disk 4 - 128GB
resource "truenas_vm_device" "data_disk_4" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "DISK"
    create_zvol = true
    zvol_name = "${var.pool_name}/vm-test5-data4"
    zvol_volsize = 137438953472
    type = "VIRTIO"
  })
  order = 1004
}

# CD-ROM for TrueNAS ISO
resource "truenas_vm_device" "cdrom" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "CDROM"
    path = local.truenas_iso_path
  })
  order = 1005
}

# Network interface
resource "truenas_vm_device" "nic" {
  vm = truenas_vm.test_truenas.id
  attributes = jsonencode({
    dtype = "NIC"
    type = "VIRTIO"
    nic_attach = var.bridge_interface
  })
  order = 1006
}
