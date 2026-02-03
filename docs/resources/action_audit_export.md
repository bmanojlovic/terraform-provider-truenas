---
page_title: "truenas_action_audit_export Resource - terraform-provider-truenas"
subcategory: "Actions - Audit"
description: |-
  Generate an audit report based on the specified `query-filters` and `query-options` for the specified `services` in the specified `export_format`.  Supported export_formats are CSV, JSON, and YAML. The endpoint returns a local filesystem path where the resulting audit report is located.
---

# truenas_action_audit_export (Resource)

Generate an audit report based on the specified `query-filters` and `query-options` for the specified `services` in the specified `export_format`.  Supported export_formats are CSV, JSON, and YAML. The endpoint returns a local filesystem path where the resulting audit report is located.

This is an action resource that executes the `audit.export` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_audit_export" "example" {
}
```

## Schema

### Input Parameters

- `data` (String, Optional) Audit export configuration specifying services, filters, and format.

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
