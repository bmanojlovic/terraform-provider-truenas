# Test flat action without prefix (no conflict with CRUD resource)
resource "truenas_filesystem_mkdir" "flat_test" {
  path         = "/mnt/fast/terraform-flat-mkdir-test"
  options_mode = "750"
}

output "flat_mkdir_result" {
  value = truenas_filesystem_mkdir.flat_test.result
}
