output "vm_id" {
  description = "TrueNAS test VM ID"
  value       = truenas_vm.test_truenas.id
}

output "vm_name" {
  description = "TrueNAS test VM name"
  value       = truenas_vm.test_truenas.name
}

output "device_ids" {
  description = "VM device IDs"
  value = {
    boot_disk  = truenas_vm_device.boot_disk.id
    data_disk_1 = truenas_vm_device.data_disk_1.id
    data_disk_2 = truenas_vm_device.data_disk_2.id
    data_disk_3 = truenas_vm_device.data_disk_3.id
    data_disk_4 = truenas_vm_device.data_disk_4.id
    cdrom      = truenas_vm_device.cdrom.id
    nic        = truenas_vm_device.nic.id
  }
}




# Data source outputs
output "vm_lookup_name" {
  description = "VM name from data source"
  value       = data.truenas_vm.test_lookup.name
}

output "vm_lookup_memory" {
  description = "VM memory from data source"
  value       = data.truenas_vm.test_lookup.memory
}

# Query data source outputs
output "all_vms_count" {
  description = "Total number of VMs"
  value       = length(data.truenas_vms.all.items)
}

output "all_vm_names" {
  description = "Names of all VMs"
  value       = [for vm in data.truenas_vms.all.items : vm.name]
}
