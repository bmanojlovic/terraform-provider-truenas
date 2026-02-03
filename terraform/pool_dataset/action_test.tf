# Test flat action with prefix (exact match: pool.scrub is both CRUD and action)
resource "truenas_action_pool_scrub" "test" {
  id     = 1  # fast pool
  action = "START"
}

output "scrub_state" {
  value = truenas_action_pool_scrub.test.state
}

# Test flat action without prefix (vm.restart doesn't conflict with vm)
# resource "truenas_vm_restart" "test" {
#   id = 1
# }
