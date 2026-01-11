variable "truenas_host" {
  description = "TrueNAS host address"
  type        = string
}

variable "truenas_token" {
  description = "TrueNAS API token"
  type        = string
  sensitive   = true
}

variable "pool_name" {
  description = "Pool name to create datasets in"
  type        = string
  default     = "tank"
}
