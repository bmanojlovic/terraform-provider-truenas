---
page_title: "truenas_initshutdownscript Resource - terraform-provider-truenas"
subcategory: ""
description: |-
  Init/shutdown script configuration data for creation.
---

# truenas_initshutdownscript (Resource)

Init/shutdown script configuration data for creation.

## Example Usage

```terraform
resource "truenas_initshutdownscript" "example" {
  type = "COMMAND"
  when = "PREINIT"
  command = ""
  script = ""
  enabled = true
  timeout = 10
  comment = ""
}
```

## Schema

### Required

- `type` (Required) - Type of init/shutdown script to execute.

* `COMMAND`: Execute a single command
* `SCRIPT`: Execute a script file Valid values: `COMMAND`, `SCRIPT`. Type: `string`
- `when` (Required) - * "PREINIT": Early in the boot process before all services have started.
* "POSTINIT": Late in the boot process when most services have started.
* "SHUTDOWN": On shutdown. Valid values: `PREINIT`, `POSTINIT`, `SHUTDOWN`. Type: `string`

### Optional

- `command` (Optional) - Must be given if `type="COMMAND"`. Default: ``. Type: `string`
- `script` (Optional) - Must be given if `type="SCRIPT"`. Default: ``. Type: `string`
- `enabled` (Optional) - Whether the init/shutdown script is enabled to execute. Default: `True`. Type: `boolean`
- `timeout` (Optional) - An integer time in seconds that the system should wait for the execution of the script/command.

A hard limit for a timeout is configured by the base OS, so when a script/command is set to execute on SHUTDOWN,     the hard limit configured by the base OS is changed adding the timeout specified by script/command so it can be     ensured that it executes as desired and is not interrupted by the base OS's limit. Default: `10`. Type: `integer`
- `comment` (Optional) - Optional comment describing the purpose of this script. Default: ``. Type: `string`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import truenas_initshutdownscript.example <id>
```
