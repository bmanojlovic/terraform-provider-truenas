output "keypair_id" {
  description = "SSH keypair credential ID"
  value       = truenas_keychaincredential.ssh_keypair.id
}

output "keypair_name" {
  description = "SSH keypair credential name"
  value       = truenas_keychaincredential.ssh_keypair.name
}
