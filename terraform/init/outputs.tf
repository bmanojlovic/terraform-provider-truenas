output "vm_id" {
  description = "ID of the nested TrueNAS VM"
  value       = truenas_vm.test.id
}

output "vm_name" {
  description = "Name of the nested TrueNAS VM"
  value       = truenas_vm.test.name
}

output "next_steps" {
  description = "Instructions for completing setup"
  value       = <<-EOT
    Nested TrueNAS VM created and started!
    
    Next steps:
    1. VM is booting from ISO - access console from TrueNAS UI
    2. Install TrueNAS SCALE
    3. Configure network (bridged to ${var.bridge_interface})
    4. Access web UI at VM's IP address
    5. Create ZFS pool with 4x 128GB disks
    6. Generate API token: System Settings → API Keys
    7. Test provider against nested TrueNAS (safe to break!)
    
    When done testing: terraform destroy
  EOT
}
