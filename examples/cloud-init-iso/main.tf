# Cloud-init ISO upload example

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

variable "truenas_host" {
  type = string
}

variable "truenas_token" {
  type      = string
  sensitive = true
}

variable "pool_name" {
  type    = string
  default = "tank"
}

variable "ssh_public_key" {
  type = string
}

# Generate cloud-init user-data locally
resource "local_file" "user_data" {
  filename = "${path.module}/.terraform/cloud-init/user-data"
  content  = <<-EOF
    #cloud-config
    hostname: testtruenas
    users:
      - name: admin
        ssh_authorized_keys:
          - ${var.ssh_public_key}
        sudo: ALL=(ALL) NOPASSWD:ALL
    network:
      version: 2
      ethernets:
        eth0:
          addresses: [192.168.1.50/24]
          gateway4: 192.168.1.1
          nameservers:
            addresses: [8.8.8.8]
  EOF
}

resource "local_file" "meta_data" {
  filename = "${path.module}/.terraform/cloud-init/meta-data"
  content  = "instance-id: testtruenas-001\nlocal-hostname: testtruenas\n"
}

# Generate ISO locally using external data source
data "external" "cloud_init_iso" {
  depends_on = [local_file.user_data, local_file.meta_data]
  
  program = ["bash", "-c", <<-EOT
    mkdir -p ${path.module}/.terraform
    genisoimage -output ${path.module}/.terraform/cloud-init.iso \
      -volid cidata -joliet -rock ${path.module}/.terraform/cloud-init/ 2>&1 >/dev/null
    base64 -w 0 ${path.module}/.terraform/cloud-init.iso > ${path.module}/.terraform/cloud-init.iso.b64
    echo "{\"content\": \"$(cat ${path.module}/.terraform/cloud-init.iso.b64)\"}"
  EOT
  ]
}

# Upload ISO to TrueNAS
resource "truenas_filesystem_put" "cloud_init_iso" {
  path    = "/mnt/${var.pool_name}/isos/cloud-init.iso"
  content = data.external.cloud_init_iso.result.content
  mode    = 0644
}

output "iso_path" {
  value = truenas_filesystem_put.cloud_init_iso.path
}
