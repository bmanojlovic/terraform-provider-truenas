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

# --- Test 1: SMB user + share + ACL with ae_who_id ---

# Create SMB-enabled user
resource "truenas_user" "smbtest" {
  username     = "smbtest"
  full_name    = "SMB Test User"
  password     = "Test1234!"
  smb          = true
  group_create = true
}

# Create a dataset for the test share
resource "truenas_pool_dataset" "smbtest" {
  name = "fast/smbtest"
  type = "FILESYSTEM"
}

# Create SMB share
resource "truenas_sharing_smb" "testshare" {
  name = "smbtest"
  path = "/mnt/fast/smbtest"

  depends_on = [truenas_pool_dataset.smbtest]
}

# Set ACL using ae_who_id (user UID)
resource "truenas_action_sharing_smb_setacl" "test_by_id" {
  smb_setacl = jsonencode({
    share_name = truenas_sharing_smb.testshare.name
    share_acl = [
      {
        ae_perm   = "FULL"
        ae_type   = "ALLOWED"
        ae_who_id = { id_type = "USER", id = truenas_user.smbtest.uid }
      }
    ]
  })

  depends_on = [truenas_user.smbtest, truenas_sharing_smb.testshare]
}

# --- Test 2: ACL with ae_who_sid (Everyone = S-1-1-0) ---

resource "truenas_pool_dataset" "smbtest2" {
  name = "fast/smbtest2"
  type = "FILESYSTEM"
}

resource "truenas_sharing_smb" "testshare2" {
  name = "smbtest2"
  path = "/mnt/fast/smbtest2"

  depends_on = [truenas_pool_dataset.smbtest2]
}

# Set ACL using SID - restrict to read-only for everyone
resource "truenas_action_sharing_smb_setacl" "test_by_sid" {
  smb_setacl = jsonencode({
    share_name = truenas_sharing_smb.testshare2.name
    share_acl = [
      {
        ae_perm    = "READ"
        ae_type    = "ALLOWED"
        ae_who_sid = "S-1-1-0"
      }
    ]
  })

  depends_on = [truenas_sharing_smb.testshare2]
}
