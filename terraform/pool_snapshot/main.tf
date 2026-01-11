# Explicit name snapshot
resource "truenas_pool_snapshot" "explicit" {
  dataset   = "${var.pool_name}/terraform-test-fs"
  name      = "terraform-test-snapshot"
  recursive = false
}

# Naming schema snapshot
resource "truenas_pool_snapshot" "schema" {
  dataset        = "${var.pool_name}/terraform-test-fs"
  naming_schema  = "auto-%Y%m%d-%H%M%S"
  recursive      = false
}

# Recursive snapshot
resource "truenas_pool_snapshot" "recursive" {
  dataset   = "${var.pool_name}/terraform-test-parent"
  name      = "recursive-snapshot"
  recursive = true
  
  depends_on = [truenas_pool_snapshot.explicit]
}
