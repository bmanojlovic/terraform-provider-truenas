---
page_title: "truenas_action_audit_download_report Resource - terraform-provider-truenas"
subcategory: "Actions - Audit"
description: |-
  Download the audit report with the specified name from the server. Note that users will only be able to download reports that they personally generated.
---

# truenas_action_audit_download_report (Resource)

Download the audit report with the specified name from the server. Note that users will only be able to download reports that they personally generated.

This is an action resource that executes the `audit.download_report` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_audit_download_report" "example" {
  data = "value"
}
```

## Schema

### Input Parameters

- `data` (String, Required) AuditDownloadReportArgs parameters.

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
