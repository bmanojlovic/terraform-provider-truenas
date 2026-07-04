# --- Generate cloud-init ISOs per node ---

resource "local_file" "user_data" {
  for_each = { for n in local.nodes : n.name => n }
  filename = "${path.module}/.terraform/cloud-init/${each.key}/user-data"
  content  = <<-EOF
    #cloud-config
    hostname: ${each.key}
    fqdn: ${each.key}.${var.domain}

    users:
      - name: root
        ssh_authorized_keys:
          - ${local.ssh_public_key}

    runcmd:
      - systemctl enable serial-getty@ttyS0.service
      - systemctl start serial-getty@ttyS0.service
  EOF
}

resource "local_file" "meta_data" {
  for_each = { for n in local.nodes : n.name => n }
  filename = "${path.module}/.terraform/cloud-init/${each.key}/meta-data"
  content  = <<-EOF
    instance-id: ${each.key}
    local-hostname: ${each.key}
  EOF
}

resource "local_file" "network_config" {
  for_each = { for n in local.nodes : n.name => n }
  filename = "${path.module}/.terraform/cloud-init/${each.key}/network-config"
  content  = <<-EOF
    version: 2
    ethernets:
      eth0:
        addresses:
          - ${each.value.ip}/24
        gateway4: ${var.gateway}
        nameservers:
          addresses:
            - ${var.dns_server}
          search:
            - ${var.domain}
  EOF
}

resource "null_resource" "cloud_init_iso" {
  for_each = { for n in local.nodes : n.name => n }

  depends_on = [local_file.user_data, local_file.meta_data, local_file.network_config]

  triggers = {
    user_data      = local_file.user_data[each.key].content
    meta_data      = local_file.meta_data[each.key].content
    network_config = local_file.network_config[each.key].content
  }

  provisioner "local-exec" {
    command = <<-EOT
      mkisofs -output ${path.module}/.terraform/cloud-init/${each.key}.iso \
        -volid cidata -rock ${path.module}/.terraform/cloud-init/${each.key}/ 2>/dev/null
      base64 -w 0 ${path.module}/.terraform/cloud-init/${each.key}.iso \
        > ${path.module}/.terraform/cloud-init/${each.key}.iso.b64
    EOT
  }
}

data "local_file" "cloud_init_iso_b64" {
  for_each   = { for n in local.nodes : n.name => n }
  depends_on = [null_resource.cloud_init_iso]
  filename   = "${path.module}/.terraform/cloud-init/${each.key}.iso.b64"
}

resource "truenas_filesystem_put" "cloud_init_iso" {
  for_each     = { for n in local.nodes : n.name => n }
  path         = "/mnt/${var.iso_pool}/iso-store/cloud-init-${each.key}.iso"
  file_content = data.local_file.cloud_init_iso_b64[each.key].content
}
