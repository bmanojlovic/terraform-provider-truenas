variable "truenas_host" {
  description = "TrueNAS host IP or hostname"
  type        = string
}

variable "truenas_token" {
  description = "TrueNAS API token"
  type        = string
  sensitive   = true
}

variable "pool_name" {
  description = "Storage pool name for VM disks"
  type        = string
  default     = "tank"
}

variable "truenas_version" {
  description = "TrueNAS SCALE version to install"
  type        = string
  default     = "25.10.1"
}

variable "truenas_iso_path" {
  description = "Path to TrueNAS ISO file on host (leave empty to use default based on version)"
  type        = string
  default     = ""
}

variable "truenas_iso_url" {
  description = "URL to download TrueNAS ISO from (leave empty to use default based on version)"
  type        = string
  default     = ""
}

variable "bridge_interface" {
  description = "Bridge interface for VM network"
  type        = string
  default     = "eno1"
}
