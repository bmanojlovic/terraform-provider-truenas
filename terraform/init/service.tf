# Enable SSH service
resource "truenas_service" "ssh" {
  enable = true
}
