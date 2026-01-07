# Action Handling Strategy

This document explains how the TrueNAS Terraform provider handles different types of actions and operations.

## Overview

The provider implements three patterns for managing TrueNAS resources and their operations:

1. **CRUD Resources** - Standard create/read/update/delete lifecycle
2. **Lifecycle Actions** - Built-in start/stop behavior for stateful resources
3. **Action Resources** - Separate resources for operational actions

## 1. CRUD Resources (51 resources)

Standard Terraform resources that manage infrastructure state:

```terraform
resource "truenas_pool_dataset" "data" {
  name = "tank/data"
  type = "FILESYSTEM"
}

resource "truenas_user" "admin" {
  username  = "admin"
  full_name = "Administrator"
  password  = "secure-password"
}

resource "truenas_sharing_nfs" "export" {
  path = "/mnt/tank/data"
}
```

**Behavior:**
- `terraform apply` creates the resource
- `terraform apply` (again) updates if configuration changes
- `terraform destroy` deletes the resource
- State is tracked and managed by Terraform

**Available CRUD Resources:**
- Storage: `pool_dataset`, `pool_snapshot`, `pool_scrub`
- Users & Groups: `user`, `group`, `privilege`
- Sharing: `sharing_nfs`, `sharing_smb`
- VMs: `vm`, `vm_device`, `virt_instance`
- Apps: `app`, `app_registry`
- Network: `interface`, `staticroute`
- Services: `cloudsync`, `replication`, `cronjob`, `rsynctask`
- And 30+ more...

## 2. Lifecycle Actions (3 resources)

Resources with built-in lifecycle management via `start_on_create` field:

### Virtual Machines

```terraform
resource "truenas_vm" "server" {
  name            = "ubuntu-server"
  vcpus           = 4
  memory          = 8192
  start_on_create = true  # Default: true
}
```

**Behavior:**
- VM is created in stopped state
- If `start_on_create = true` (default), provider calls `vm.start` after creation
- If `start_on_create = false`, VM remains stopped
- Matches AWS EC2 behavior where instances start automatically

### Applications

```terraform
resource "truenas_app" "plex" {
  app_name        = "plex"
  train           = "stable"
  start_on_create = true  # Default: true
}
```

**Behavior:**
- App is deployed
- If `start_on_create = true` (default), provider calls `app.start` after deployment
- App is ready to use immediately after `terraform apply`

### Virtual Instances (Containers)

```terraform
resource "truenas_virt_instance" "nginx" {
  name            = "nginx-proxy"
  image           = "nginx:latest"
  start_on_create = true  # Default: true
}
```

**Behavior:**
- Container is created
- If `start_on_create = true` (default), provider calls `virt/instance.start`
- Container is running after creation

## 3. Action Resources (19 resources)

Separate resources for operational actions that don't fit CRUD model:

### Pattern

```terraform
# 1. Create the base resource
resource "truenas_<resource>" "example" {
  # ... configuration
}

# 2. Execute action on the resource
resource "truenas_<resource>_<action>_action" "example" {
  resource_id = truenas_<resource>.example.id
  # ... action-specific parameters
}
```

### Examples

#### Pool Scrub

```terraform
# Schedule weekly scrubs
resource "truenas_pool_scrub" "weekly" {
  pool     = truenas_pool.tank.id
  schedule = {
    dow  = "0"  # Sunday
    hour = "2"  # 2 AM
  }
}

# Run scrub immediately
resource "truenas_pool_scrub_run_action" "now" {
  resource_id = truenas_pool_scrub.weekly.id
}
```

#### CloudSync

```terraform
# Configure sync schedule
resource "truenas_cloudsync" "backup" {
  path        = "/mnt/tank/data"
  credentials = truenas_cloudsync_credentials.s3.id
  schedule    = { hour = "3" }
}

# Trigger immediate sync
resource "truenas_cloudsync_sync_action" "now" {
  resource_id = truenas_cloudsync.backup.id
}
```

#### App Upgrade

```terraform
resource "truenas_app" "plex" {
  app_name = "plex"
  train    = "stable"
}

# Upgrade to latest version
resource "truenas_app_upgrade_action" "upgrade_plex" {
  resource_id = truenas_app.plex.id
}
```

#### Replication

```terraform
resource "truenas_replication" "offsite" {
  name      = "offsite-backup"
  transport = "SSH"
  schedule  = { hour = "1" }
}

# Run replication immediately
resource "truenas_replication_run_action" "now" {
  resource_id = truenas_replication.offsite.id
}
```

#### CronJob

```terraform
resource "truenas_cronjob" "cleanup" {
  command  = "/scripts/cleanup.sh"
  schedule = { hour = "4" }
}

# Execute job immediately
resource "truenas_cronjob_run_action" "now" {
  resource_id = truenas_cronjob.cleanup.id
}
```

### Action Resource Behavior

**On Create:**
- Executes the action immediately
- Stores timestamp-based ID in state
- Returns success/failure

**On Update:**
- Re-executes the action with new parameters
- Updates timestamp ID

**On Destroy:**
- Removes from Terraform state
- **Cannot undo the action** (actions are irreversible)

**Important Notes:**
- Actions execute every time `terraform apply` runs if the resource changes
- Use with caution - actions are immediate and cannot be rolled back
- Consider using `lifecycle { ignore_changes = [...] }` to prevent accidental re-execution
- Best for one-time operations or triggered workflows

### Available Action Resources

| Base Resource | Actions | Use Case |
|--------------|---------|----------|
| `app` | `upgrade`, `rollback`, `redeploy` | App lifecycle management |
| `cloudsync` | `sync`, `sync_onetime`, `restore` | Trigger cloud sync operations |
| `cloud_backup` | `sync`, `restore` | Backup operations |
| `pool_scrub` | `run`, `scrub` | Manual pool scrubbing |
| `pool_snapshot` | `rollback` | Restore from snapshot |
| `pool_snapshottask` | `run` | Execute snapshot task |
| `replication` | `run`, `run_onetime`, `restore` | Trigger replication |
| `rsynctask` | `run` | Execute rsync task |
| `cronjob` | `run` | Execute cron job immediately |

## When to Use Each Pattern

### Use CRUD Resources When:
- Managing infrastructure state (pools, datasets, users, shares)
- Configuration should persist and be managed by Terraform
- Changes should be tracked and reversible

### Use Lifecycle Actions When:
- Resource needs to be started after creation (VMs, apps, containers)
- Default behavior should match cloud providers (AWS, Azure, GCP)
- State management is automatic

### Use Action Resources When:
- Triggering one-time operations (upgrade, sync, scrub)
- Executing scheduled tasks immediately
- Performing maintenance operations
- Testing or validating configurations

### Use External Tools When:
- Complex orchestration across multiple systems
- Operations requiring approval workflows
- Integration with existing automation (Ansible, scripts)
- Operations that need retry logic or error handling

## External Tool Integration

For operations not covered by action resources, use Terraform provisioners or external tools:

### Local-Exec Provisioner

```terraform
resource "truenas_vm" "server" {
  name   = "server"
  vcpus  = 2
  memory = 4096

  provisioner "local-exec" {
    command = "curl -X POST https://${var.truenas_host}/api/v2.0/vm/id/${self.id}/start"
  }
}
```

### Null Resource with Triggers

```terraform
resource "null_resource" "scrub_pool" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "truenas-cli pool scrub run --pool tank"
  }
}
```

### Ansible Integration

```terraform
resource "truenas_pool_dataset" "data" {
  name = "tank/data"
}

resource "null_resource" "configure_dataset" {
  depends_on = [truenas_pool_dataset.data]

  provisioner "local-exec" {
    command = "ansible-playbook configure-dataset.yml -e dataset_id=${truenas_pool_dataset.data.id}"
  }
}
```

## Best Practices

1. **Use CRUD resources for state management** - Let Terraform track and manage infrastructure
2. **Use lifecycle actions for automatic startup** - VMs and apps should start by default
3. **Use action resources sparingly** - Only for operations that must be triggered via Terraform
4. **Document action resource usage** - Add comments explaining why actions are needed
5. **Consider external tools for complex workflows** - Ansible, scripts, or CI/CD pipelines
6. **Test action resources in non-production** - Actions are immediate and irreversible
7. **Use `lifecycle` blocks to prevent accidental execution**:

```terraform
resource "truenas_app_upgrade_action" "upgrade" {
  resource_id = truenas_app.plex.id

  lifecycle {
    ignore_changes = [resource_id]  # Only upgrade when explicitly changed
  }
}
```

## Migration from Other Providers

If migrating from REST-based providers or manual management:

1. **Import existing resources**: `terraform import truenas_vm.server 123`
2. **Set `start_on_create = false`** for already-running VMs to prevent restart
3. **Avoid action resources during initial import** - They will execute immediately
4. **Use `terraform plan`** to verify changes before applying

## Summary

| Pattern | Count | Use For | Reversible |
|---------|-------|---------|------------|
| CRUD Resources | 51 | Infrastructure state | Yes |
| Lifecycle Actions | 3 | Auto-start VMs/apps | Yes |
| Action Resources | 19 | Operational tasks | No |
| Data Sources | 56 | Read-only queries | N/A |

**Total: 70 resources + 56 data sources = 126 Terraform resources**
