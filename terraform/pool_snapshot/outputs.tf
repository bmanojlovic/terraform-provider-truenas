output "explicit_id" {
  description = "Explicit snapshot ID"
  value       = truenas_pool_snapshot.explicit.id
}

output "schema_id" {
  description = "Schema-based snapshot ID"
  value       = truenas_pool_snapshot.schema.id
}

output "recursive_id" {
  description = "Recursive snapshot ID"
  value       = truenas_pool_snapshot.recursive.id
}
