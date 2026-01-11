# SSH Key Pair credential
resource "truenas_keychaincredential" "ssh_keypair" {
  name = "terraform-test-keypair"
  type = "SSH_KEY_PAIR"
  
  attributes = jsonencode({
    private_key = file("${path.module}/test_key")
    public_key  = file("${path.module}/test_key.pub")
  })
}

# Note: SSH_CREDENTIALS requires an existing SSH_KEY_PAIR credential ID
# Skipping SSH_CREDENTIALS test as it needs real host and key reference
