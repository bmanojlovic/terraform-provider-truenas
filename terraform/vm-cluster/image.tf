# --- Fetch cloud image to TrueNAS via container ---

resource "truenas_app" "fetch_image" {
  app_name   = "${var.cluster_name}fetch"
  custom_app = true
  custom_compose_config_string = <<-YAML
    services:
      fetch:
        image: alpine:latest
        command: ["sh", "-c", "test -f /target/${local.image_filename} || wget -O /target/${local.image_filename} ${var.image_url}"]
        volumes:
          - /mnt/${var.iso_pool}/iso-store:/target
  YAML
}

# --- Create template zvol and import image into it ---

resource "truenas_pool_dataset" "template_zvol" {
  type    = "VOLUME"
  name    = "${var.pool}/${var.cluster_name}template"
  volsize = var.disk_size
}

resource "truenas_vm_import_disk_image" "template" {
  depends_on = [truenas_pool_dataset.template_zvol, truenas_app.fetch_image]
  diskimg    = local.image_stor_path
  zvol       = "${var.pool}/${var.cluster_name}template"
}

# --- Snapshot template for cloning ---

resource "truenas_pool_snapshot" "template" {
  depends_on = [truenas_vm_import_disk_image.template]
  dataset    = "${var.pool}/${var.cluster_name}template"
  name       = "base"
  recursive  = false

  lifecycle {
    ignore_changes = [naming_schema]
  }
}

# --- Clone template zvol for each node ---

resource "truenas_pool_snapshot_clone" "node" {
  for_each    = { for n in local.nodes : n.name => n }
  depends_on  = [truenas_pool_snapshot.template]
  snapshot    = "${var.pool}/${var.cluster_name}template@base"
  dataset_dst = "${var.pool}/${each.key}"
}
