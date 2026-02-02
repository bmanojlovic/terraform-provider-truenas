variable "truenas_host" {
  description = "TrueNAS host address"
  type        = string
}

variable "truenas_token" {
  description = "TrueNAS API token"
  type        = string
  sensitive   = true
}

variable "truenas_pool" {
  description = "Pool name to create datasets in"
  type        = string
}

variable "pool_name" {
  description = "Pool name (alias for truenas_pool)"
  type        = string
  default     = null
}
