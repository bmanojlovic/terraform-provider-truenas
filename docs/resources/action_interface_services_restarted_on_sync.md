---
page_title: "truenas_action_interface_services_restarted_on_sync Resource - terraform-provider-truenas"
subcategory: "Actions - Network"
description: |-
  Returns which services will be set to listen on 0.0.0.0 (and, thus, restarted) on sync.  Example result: [     // Samba service will be set ot listen on 0.0.0.0 and restarted because it was set up to listen on     // 192.168.0.1 which is being removed.     {"type": "SYSTEM_SERVICE", "service": "cifs", "ips": ["192.168.0.1"]}, ]
---

# truenas_action_interface_services_restarted_on_sync (Resource)

Returns which services will be set to listen on 0.0.0.0 (and, thus, restarted) on sync.  Example result: [     // Samba service will be set ot listen on 0.0.0.0 and restarted because it was set up to listen on     // 192.168.0.1 which is being removed.     {"type": "SYSTEM_SERVICE", "service": "cifs", "ips": ["192.168.0.1"]}, ]

This is an action resource that executes the `interface.services_restarted_on_sync` operation. Actions are triggered on resource creation and cannot be undone on destroy.

## Example Usage

```terraform
resource "truenas_action_interface_services_restarted_on_sync" "example" {
}
```

## Schema

### Input Parameters

None

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
