---
page_title: "truenas_replication Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create a Replication Task that will push or pull ZFS snapshots to or from remote host.
---

# truenas_replication (Resource)

Create a Replication Task that will push or pull ZFS snapshots to or from remote host.


## Example Usage

```terraform
resource "truenas_replication" "example" {
  auto = true
  direction = "example-value"
  name = "example-value"
  recursive = true
  retention_policy = "example-value"
  source_datasets = ["item1"]
  target_dataset = "example-value"
  transport = "example-value"
}
```

## Schema

### Required

- `auto` (Bool) - Allow replication to run automatically on schedule or after bound periodic snapshot task.
- `direction` (String) - Whether task will `PUSH` or `PULL` snapshots. Valid values: `PUSH`, `PULL`
- `name` (String) - Name for replication task.
- `recursive` (Bool) - Whether to recursively replicate child datasets.
- `retention_policy` (String) - How to delete old snapshots on target side:  * `SOURCE`: Delete snapshots that are absent on source side. * `CUSTOM`: Delete snapshots that are older than `lifetime_value` and `lifetime_unit`. * `NONE Valid values: `SOURCE`, `CUSTOM`, `NONE`
- `source_datasets` (List) - List of datasets to replicate snapshots from.
- `target_dataset` (String) - Dataset to put snapshots into.
- `transport` (String) - Method of snapshots transfer.  * `SSH` transfers snapshots via SSH connection. This method is supported everywhere but does not achieve       great performance. * `SSH+NETCAT` uses unencrypted connect Valid values: `SSH`, `SSH+NETCAT`, `LOCAL`

### Optional

- `allow_from_scratch` (Bool) - Will destroy all snapshots on target side and replicate everything from scratch if none of the snapshots on     target side matches source snapshots. Default: `False`
- `also_include_naming_schema` (List) - List of naming schemas for push replication. Default: `[]`
- `compressed` (Bool) - Enable compressed ZFS send streams. Default: `True`
- `compression` (String) - Compresses SSH stream. Available only for SSH transport. Default: `None`
- `embed` (Bool) - Enable embedded block support for ZFS send streams. Default: `False`
- `enabled` (Bool) - Whether this replication task is enabled. Default: `True`
- `encryption` (Bool) - Whether to enable encryption for the replicated datasets. Default: `False`
- `encryption_inherit` (String) - Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled. Default: `None`
- `encryption_key` (String) - Encryption key for replicated datasets. `null` if not specified. Default: `None`
- `encryption_key_format` (String) - Format of the encryption key.  * `HEX`: Hexadecimal-encoded key * `PASSPHRASE`: Text passphrase * `null`: Not applicable when encryption is disabled Default: `None`
- `encryption_key_location` (String) - Filesystem path where encryption key is stored. `null` if not using key file. Default: `None`
- `exclude` (List) - Array of dataset patterns to exclude from replication. Default: `[]`
- `hold_pending_snapshots` (Bool) - Prevent source snapshots from being deleted by retention of replication fails for some reason. Default: `False`
- `large_block` (Bool) - Enable large block support for ZFS send streams. Default: `True`
- `lifetime_unit` (String) - Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` is CUSTOM. Default: `None`
- `lifetime_value` (Int64) - Number of time units to retain snapshots for custom retention policy. Only applies when `retention_policy` is     CUSTOM. Default: `None`
- `lifetimes` (List) - Array of different retention schedules with their own cron schedules and lifetime settings. Default: `[]`
- `logging_level` (String) - Log level for replication task execution. Controls verbosity of replication logs. Default: `None`
- `name_regex` (String) - Replicate all snapshots which names match specified regular expression. Default: `None`
- `naming_schema` (List) - List of naming schemas for pull replication. Default: `[]`
- `netcat_active_side` (String) - Which side actively establishes the netcat connection for `SSH+NETCAT` transport.  * `LOCAL`: Local system initiates the connection * `REMOTE`: Remote system initiates the connection * `null`: Not app Default: `None`
- `netcat_active_side_listen_address` (String) - IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable. Default: `None`
- `netcat_active_side_port_max` (Int64) - Maximum port number in the range for netcat connections. `null` if not applicable. Default: `None`
- `netcat_active_side_port_min` (Int64) - Minimum port number in the range for netcat connections. `null` if not applicable. Default: `None`
- `netcat_passive_side_connect_address` (String) - IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable. Default: `None`
- `only_matching_schedule` (Bool) - Will only replicate snapshots that match `schedule` or `restrict_schedule`. Default: `False`
- `periodic_snapshot_tasks` (List) - List of periodic snapshot task IDs that are sources of snapshots for this replication task. Only push     replication tasks can be bound to periodic snapshot tasks. Default: `[]`
- `properties` (Bool) - Send dataset properties along with snapshots. Default: `True`
- `properties_exclude` (List) - Array of dataset property names to exclude from replication. Default: `[]`
- `properties_override` (String) - Object mapping dataset property names to override values during replication. Default: `{}`
- `readonly` (String) - Controls destination datasets readonly property.  * `SET`: Set all destination datasets to readonly=on after finishing the replication. * `REQUIRE`: Require all existing destination datasets to have r Default: `SET` Valid values: `SET`, `REQUIRE`, `IGNORE`
- `replicate` (Bool) - Whether to use full ZFS replication. Default: `False`
- `restrict_schedule` (String) - Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have periodic     snapshot tasks that run every 15 minutes, but only run replication task every hour. Default: `None`
- `retries` (Int64) - Number of retries before considering replication failed. Default: `5`
- `schedule` (String) - Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot tasks can have     a schedule. Default: `None`
- `speed_limit` (Int64) - Limits speed of SSH stream. Available only for SSH transport. Default: `None`
- `ssh_credentials` (Int64) - Keychain Credential ID of type `SSH_CREDENTIALS`. Default: `None`
- `sudo` (Bool) - `SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs`     command on the remote machine. Default: `False`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_replication.example <id>
```
