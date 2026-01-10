---
page_title: "truenas_initshutdownscript Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Create an initshutdown script task.
---

# truenas_initshutdownscript (Resource)

Create an initshutdown script task.

## Example Usage

```terraform
resource "truenas_initshutdownscript" "example" {
  type = "example-value"
  when = "example-value"
}
```

## Schema

### Required

- `type` (String) - Type of init/shutdown script to execute.  * `COMMAND`: Execute a single command * `SCRIPT`: Execute a script file Valid values: `COMMAND`, `SCRIPT`
- `when` (String) - * "PREINIT": Early in the boot process before all services have started. * "POSTINIT": Late in the boot process when most services have started. * "SHUTDOWN": On shutdown. Valid values: `PREINIT`, `POSTINIT`, `SHUTDOWN`

### Optional

- `command` (String) - Must be given if `type="COMMAND"`. Default: ``
- `comment` (String) - Optional comment describing the purpose of this script. Default: ``
- `enabled` (Bool) - Whether the init/shutdown script is enabled to execute. Default: `True`
- `script` (String) - Must be given if `type="SCRIPT"`. Default: ``
- `timeout` (Int64) - An integer time in seconds that the system should wait for the execution of the script/command.  A hard limit for a timeout is configured by the base OS, so when a script/command is set to execute on  Default: `10`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_initshutdownscript.example <id>
```
