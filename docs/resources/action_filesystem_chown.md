---
page_title: "truenas_action_filesystem_chown Resource - terraform-provider-truenas"
subcategory: "Actions - Filesystem"
description: |-
  Change owner or group of file at `path`.  `uid` and `gid` specify new owner of the file. If either key is absent or None, then existing value on the file is not changed.  `user` and `group` alternatively allow specifying a uid gid by user name or group name.  `recursive` performs action recursively, but does not traverse filesystem mount points.  If `traverse` and `recursive` are specified, then the chown operation will traverse filesystem mount points.
---

# truenas_action_filesystem_chown (Resource)

Change owner or group of file at `path`.  `uid` and `gid` specify new owner of the file. If either key is absent or None, then existing value on the file is not changed.  `user` and `group` alternatively allow specifying a uid gid by user name or group name.  `recursive` performs action recursively, but does not traverse filesystem mount points.  If `traverse` and `recursive` are specified, then the chown operation will traverse filesystem mount points.

This is an action resource that executes the `filesystem.chown` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_filesystem_chown" "example" {
  filesystem_chown = "value"
}
```

## Schema

### Input Parameters

- `filesystem_chown` (String, Required) FilesystemChownArgs parameters.

### Computed Outputs

- `action_id` (String) Unique identifier for this action execution
- `job_id` (Int64) Background job ID (if applicable)
- `state` (String) Job state: SUCCESS, FAILED, or RUNNING
- `progress` (Float64) Job progress percentage (0-100)
- `result` (String) Action result data
- `error` (String) Error message if action failed

## Notes

- Actions execute immediately when the resource is created
- Background jobs are monitored until completion
- Progress updates are logged during execution
- The resource cannot be updated - changes force recreation
- Destroying the resource does not undo the action
