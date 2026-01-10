# Generate ISO locally
resource "null_resource" "cloud_init_iso" {
  depends_on = [local_file.user_data, local_file.meta_data]
  
  triggers = {
    user_data = local_file.user_data.content
    meta_data = local_file.meta_data.content
  }
  
  provisioner "local-exec" {
    command = <<-EOT
      set -e
      mkdir -p ${path.module}/.terraform
      
      # Use mkisofs or genisoimage (whichever is available)
      if command -v mkisofs >/dev/null 2>&1; then
        mkisofs -output ${path.module}/.terraform/cloud-init.iso \
          -volid cidata -rock ${path.module}/.terraform/cloud-init/ 2>&1 >/dev/null
      elif command -v genisoimage >/dev/null 2>&1; then
        genisoimage -output ${path.module}/.terraform/cloud-init.iso \
          -volid cidata -rock ${path.module}/.terraform/cloud-init/ 2>&1 >/dev/null
      else
        echo "ERROR: Neither mkisofs nor genisoimage found" >&2
        exit 1
      fi
      
      # Encode to base64
      base64 -w 0 ${path.module}/.terraform/cloud-init.iso > ${path.module}/.terraform/cloud-init.iso.b64
    EOT
  }
}

# Read the base64 encoded ISO
data "local_file" "cloud_init_iso_b64" {
  depends_on = [null_resource.cloud_init_iso]
  filename   = "${path.module}/.terraform/cloud-init.iso.b64"
}

# Upload ISO to TrueNAS
resource "truenas_filesystem_put" "cloud_init_iso" {
  path    = "/mnt/fast/iso-store/cloud-init-${var.vm_hostname}.iso"
  content = data.local_file.cloud_init_iso_b64.content
}
