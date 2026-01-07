resource "truenas_vm" "example" {
  name        = "test-vm"
  description = "Test VM created by Terraform"
  vcpus       = 2
  memory      = 2147483648  # 2GB in bytes
  autostart   = false
}

resource "truenas_pool_dataset" "example" {
  name = "tank/terraform-test"
  type = "FILESYSTEM"
}

resource "truenas_user" "example" {
  username = "terraform-user"
  full_name = "Terraform Test User"
  group_create = true
  home_create = true
  shell = "/bin/bash"
}
