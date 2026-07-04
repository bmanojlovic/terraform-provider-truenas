terraform {
  required_providers {
    truenas = {
      source = "bmanojlovic/truenas"
    }
  }
}

provider "truenas" {
  host  = var.truenas_host
  token = var.truenas_token
}

locals {
  nodes = concat(
    [{ name = "${var.cluster_name}master", role = "master", ip = var.master_ip, idx = 0 }],
    [for i in range(var.worker_count) : {
      name = format("${var.cluster_name}node%02d", i + 1)
      role = "worker"
      ip   = cidrhost(var.node_network, var.worker_ip_offset + i)
      idx  = i + 1
    }]
  )
  ssh_public_key  = file(var.ssh_public_key_path)
  image_filename  = basename(var.image_url)
  image_stor_path = "/mnt/${var.iso_pool}/iso-store/${local.image_filename}"
}
