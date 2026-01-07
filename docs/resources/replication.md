---
page_title: "truenas_replication Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Configuration for creating a new replication task.
---

# truenas_replication (Resource)

Configuration for creating a new replication task.

## Example Usage

```terraform
resource "truenas_replication" "example" {
  name = "example-name"
  direction = "PUSH"
  transport = "SSH"
  source_datasets = []
  target_dataset = "example-target_dataset"
  recursive = false
  auto = false
  retention_policy = "SOURCE"
}
```

## Schema

### Required

- `name` (Required) - Name for replication task.. Type: `string`
- `direction` (Required) - Whether task will `PUSH` or `PULL` snapshots. Valid values: `PUSH`, `PULL`. Type: `string`
- `transport` (Required) - Method of snapshots transfer.

* `SSH` transfers snapshots via SSH connection. This method is supported everywhere but does not achieve       great performance.
* `SSH+NETCAT` uses unencrypted connection for data transfer. This can only be used in trusted networks       and requires a port (specified by range from `netcat_active_side_port_min` to `netcat_active_side_port_max`)       to be open on `netcat_active_side`.
* `LOCAL` replicates to or from localhost. Valid values: `SSH`, `SSH+NETCAT`, `LOCAL`. Type: `string`
- `source_datasets` (Required) - List of datasets to replicate snapshots from.. Type: `array`
- `target_dataset` (Required) - Dataset to put snapshots into.. Type: `string`
- `recursive` (Required) - Whether to recursively replicate child datasets.. Type: `boolean`
- `auto` (Required) - Allow replication to run automatically on schedule or after bound periodic snapshot task.. Type: `boolean`
- `retention_policy` (Required) - How to delete old snapshots on target side:

* `SOURCE`: Delete snapshots that are absent on source side.
* `CUSTOM`: Delete snapshots that are older than `lifetime_value` and `lifetime_unit`.
* `NONE`: Do not delete any snapshots. Valid values: `SOURCE`, `CUSTOM`, `NONE`. Type: `string`

### Optional

- `ssh_credentials` (Optional) - Keychain Credential ID of type `SSH_CREDENTIALS`.. Type: `string`
- `netcat_active_side` (Optional) - Which side actively establishes the netcat connection for `SSH+NETCAT` transport.

* `LOCAL`: Local system initiates the connection
* `REMOTE`: Remote system initiates the connection
* `null`: Not applicable for other transport types. Type: `string`
- `netcat_active_side_listen_address` (Optional) - IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable.. Type: `string`
- `netcat_active_side_port_min` (Optional) - Minimum port number in the range for netcat connections. `null` if not applicable.. Type: `string`
- `netcat_active_side_port_max` (Optional) - Maximum port number in the range for netcat connections. `null` if not applicable.. Type: `string`
- `netcat_passive_side_connect_address` (Optional) - IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable.. Type: `string`
- `sudo` (Optional) - `SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs`     command on the remote machine. Default: `False`. Type: `boolean`
- `exclude` (Optional) - Array of dataset patterns to exclude from replication. Default: `[]`. Type: `array`
- `properties` (Optional) - Send dataset properties along with snapshots. Default: `True`. Type: `boolean`
- `properties_exclude` (Optional) - Array of dataset property names to exclude from replication. Default: `[]`. Type: `array`
- `properties_override` (Optional) - Object mapping dataset property names to override values during replication. Default: `{}`. Type: `object`
- `replicate` (Optional) - Whether to use full ZFS replication. Default: `False`. Type: `boolean`
- `encryption` (Optional) - Whether to enable encryption for the replicated datasets. Default: `False`. Type: `boolean`
- `encryption_inherit` (Optional) - Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled.. Type: `string`
- `encryption_key` (Optional) - Encryption key for replicated datasets. `null` if not specified.. Type: `string`
- `encryption_key_format` (Optional) - Format of the encryption key.

* `HEX`: Hexadecimal-encoded key
* `PASSPHRASE`: Text passphrase
* `null`: Not applicable when encryption is disabled. Type: `string`
- `encryption_key_location` (Optional) - Filesystem path where encryption key is stored. `null` if not using key file.. Type: `string`
- `periodic_snapshot_tasks` (Optional) - List of periodic snapshot task IDs that are sources of snapshots for this replication task. Only push     replication tasks can be bound to periodic snapshot tasks. Default: `[]`. Type: `array`
- `naming_schema` (Optional) - List of naming schemas for pull replication. Default: `[]`. Type: `array`
- `also_include_naming_schema` (Optional) - List of naming schemas for push replication. Default: `[]`. Type: `array`
- `name_regex` (Optional) - Replicate all snapshots which names match specified regular expression.. Type: `string`
- `schedule` (Optional) - Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot tasks can have     a schedule.. Type: `string`
- `restrict_schedule` (Optional) - Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have periodic     snapshot tasks that run every 15 minutes, but only run replication task every hour.. Type: `string`
- `only_matching_schedule` (Optional) - Will only replicate snapshots that match `schedule` or `restrict_schedule`. Default: `False`. Type: `boolean`
- `allow_from_scratch` (Optional) - Will destroy all snapshots on target side and replicate everything from scratch if none of the snapshots on     target side matches source snapshots. Default: `False`. Type: `boolean`
- `readonly` (Optional) - Controls destination datasets readonly property.

* `SET`: Set all destination datasets to readonly=on after finishing the replication.
* `REQUIRE`: Require all existing destination datasets to have readonly=on property.
* `IGNORE`: Avoid this kind of behavior. Valid values: `SET`, `REQUIRE`, `IGNORE` Default: `SET`. Type: `string`
- `hold_pending_snapshots` (Optional) - Prevent source snapshots from being deleted by retention of replication fails for some reason. Default: `False`. Type: `boolean`
- `lifetime_value` (Optional) - Number of time units to retain snapshots for custom retention policy. Only applies when `retention_policy` is     CUSTOM.. Type: `string`
- `lifetime_unit` (Optional) - Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` is CUSTOM.. Type: `string`
- `lifetimes` (Optional) - Array of different retention schedules with their own cron schedules and lifetime settings. Default: `[]`. Type: `array`
- `compression` (Optional) - Compresses SSH stream. Available only for SSH transport.. Type: `string`
- `speed_limit` (Optional) - Limits speed of SSH stream. Available only for SSH transport.. Type: `string`
- `large_block` (Optional) - Enable large block support for ZFS send streams. Default: `True`. Type: `boolean`
- `embed` (Optional) - Enable embedded block support for ZFS send streams. Default: `False`. Type: `boolean`
- `compressed` (Optional) - Enable compressed ZFS send streams. Default: `True`. Type: `boolean`
- `retries` (Optional) - Number of retries before considering replication failed. Default: `5`. Type: `integer`
- `logging_level` (Optional) - Log level for replication task execution. Controls verbosity of replication logs.. Type: `string`
- `enabled` (Optional) - Whether this replication task is enabled. Default: `True`. Type: `boolean`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_replication.example <id>
```
