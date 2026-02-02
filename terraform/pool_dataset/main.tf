locals {
  pool = coalesce(var.pool_name, var.truenas_pool)
}

# FILESYSTEM dataset
resource "truenas_pool_dataset" "test_filesystem" {
  name        = "${local.pool}/terraform-test-fs"
  type        = "FILESYSTEM"
  comments    = "Test filesystem dataset created by Terraform"
  compression = "LZ4"
  atime       = "OFF"
  quota       = 1073741824 # 1GB
}

# VOLUME dataset (zvol)
resource "truenas_pool_dataset" "test_volume" {
  name         = "${local.pool}/terraform-test-vol"
  type         = "VOLUME"
  comments     = "Test volume dataset created by Terraform"
  compression  = "LZ4"
  volsize      = 536870912 # 512MB
  volblocksize = "16K"
  sparse       = true
}

# Nested dataset
resource "truenas_pool_dataset" "test_parent" {
  name     = "${local.pool}/terraform-test-parent"
  type     = "FILESYSTEM"
  comments = "Parent dataset"
}

resource "truenas_pool_dataset" "test_child" {
  name       = "${local.pool}/terraform-test-parent/child"
  type       = "FILESYSTEM"
  comments   = "Child dataset"
  depends_on = [truenas_pool_dataset.test_parent]
}
