output "filesystem_id" {
  description = "Filesystem dataset ID"
  value       = truenas_pool_dataset.test_filesystem.id
}

output "volume_id" {
  description = "Volume dataset ID"
  value       = truenas_pool_dataset.test_volume.id
}

output "parent_id" {
  description = "Parent dataset ID"
  value       = truenas_pool_dataset.test_parent.id
}

output "child_id" {
  description = "Child dataset ID"
  value       = truenas_pool_dataset.test_child.id
}
