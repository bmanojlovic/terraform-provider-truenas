---
page_title: "truenas_action_filesystem_setperm Resource - terraform-provider-truenas"
subcategory: "Actions - Filesystem"
description: |-
  Set unix permissions on given `path`.  If `mode` is specified then the mode will be applied to the path and files and subdirectories depending on which `options` are selected. Mode should be formatted as string representation of octal permissions bits.  `uid` the desired UID of the file user. If set to None (the default), then user is not changed.  `gid` the desired GID of the file group. If set to None (the default), then group is not changed.  `user` and `group` alternatively allow specifying the owner by name.  WARNING: `uid`, `gid, `user`, and `group` _should_ remain unset _unless_ the administrator wishes to change the owner or group of files.  `stripacl` setperm will fail if an extended ACL is present on `path`, unless `stripacl` is set to True.  `recursive` remove ACLs recursively, but do not traverse dataset boundaries.  `traverse` remove ACLs from child datasets.  If no `mode` is set, and `stripacl` is True, then non-trivial ACLs will be converted to trivial ACLs. An ACL is trivial if it can be expressed as a file mode without losing any access rules.
---

# truenas_action_filesystem_setperm (Resource)

Set unix permissions on given `path`.  If `mode` is specified then the mode will be applied to the path and files and subdirectories depending on which `options` are selected. Mode should be formatted as string representation of octal permissions bits.  `uid` the desired UID of the file user. If set to None (the default), then user is not changed.  `gid` the desired GID of the file group. If set to None (the default), then group is not changed.  `user` and `group` alternatively allow specifying the owner by name.  WARNING: `uid`, `gid, `user`, and `group` _should_ remain unset _unless_ the administrator wishes to change the owner or group of files.  `stripacl` setperm will fail if an extended ACL is present on `path`, unless `stripacl` is set to True.  `recursive` remove ACLs recursively, but do not traverse dataset boundaries.  `traverse` remove ACLs from child datasets.  If no `mode` is set, and `stripacl` is True, then non-trivial ACLs will be converted to trivial ACLs. An ACL is trivial if it can be expressed as a file mode without losing any access rules.

This is an action resource that executes the `filesystem.setperm` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_filesystem_setperm" "example" {
  filesystem_setperm = "value"
}
```

## Schema

### Input Parameters

- `filesystem_setperm` (String, Required) FilesystemSetpermArgs parameters.

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
