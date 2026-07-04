# --- Provider credentials (from TF_VAR_* environment) ---

variable "truenas_host" {
  type = string
}

variable "truenas_token" {
  type      = string
  sensitive = true
}

# --- Cluster identity ---

variable "cluster_name" {
  description = "Name prefix for VMs and zvols (alphanumeric only)"
  type        = string
  default     = "cluster"
}

variable "domain" {
  description = "DNS domain for node FQDNs"
  type        = string
  default     = "local"
}

# --- Cloud image ---

variable "image_url" {
  description = "URL to download cloud image (qcow2)"
  type        = string
  default     = "https://download.opensuse.org/distribution/leap/16.0/appliances/Leap-16.0-Minimal-VM.x86_64-Cloud.qcow2"
}

variable "iso_pool" {
  description = "Pool for ISO/image storage (typically fast/SSD pool)"
  type        = string
  default     = "fast"
}

# --- Storage ---

variable "pool" {
  description = "ZFS pool for VM zvols"
  type        = string
  default     = "tank"
}

variable "disk_size" {
  description = "Template zvol size in bytes"
  type        = number
  default     = 25769803776 # 24GB
}

# --- Network ---

variable "nic_attach" {
  description = "Network interface to attach VMs to"
  type        = string
  default     = "eno1"
}

variable "master_ip" {
  description = "Static IP for master node"
  type        = string
}

variable "node_network" {
  description = "Network CIDR for computing worker IPs"
  type        = string
}

variable "worker_ip_offset" {
  description = "Host offset in node_network for first worker"
  type        = number
  default     = 101
}

variable "netmask" {
  type    = string
  default = "255.255.255.0"
}

variable "gateway" {
  type = string
}

variable "dns_server" {
  type = string
}

# --- SSH ---

variable "ssh_public_key_path" {
  description = "Path to SSH public key for cloud-init"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

# --- VM resources ---

variable "master_vcpus" {
  type    = number
  default = 4
}

variable "master_memory_mb" {
  type    = number
  default = 4096
}

variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
}

variable "worker_vcpus" {
  type    = number
  default = 2
}

variable "worker_memory_mb" {
  type    = number
  default = 3072
}
