---
layout: "junos"
page_title: "Junos: junos_eventoptions_policy"
sidebar_current: "docs-junos-resource-eventoptions-policy"
description: |-
  Create an event-options policy
---

# junos_eventoptions_policy

Provides an event-options policy resource.

## Example Usage

```hcl
# Add an event-options policy
resource "junos_eventoptions_policy" "demo" {
  name   = "demo"
  events = ["ping_test_failed"]
  then {
    execute_commands {
      commands = ["cmd"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of policy.
* `events` - (Required)(`ListOfString`) List of events that trigger this policy.
* `then` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'then' configuration. See the [`then` arguments](#then-arguments) block.
* `attributes_match` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) List of attributes to compare for two events. Can be specified multiple times for each combination of block arguments.
  * `from` - (Required)(`String`) First attribute to compare.
  * `compare` - (Required)(`String`) Type to compare. Need to be 'equals', 'matches' or 'starts-with'.
  * `to` - (Required)(`String`) Second attribute or value to compare.
* `within` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) List of events correlated with triggering events. Can be specified multiple times for each time interval.
  * `time_interval` - (Required)(`Int`) Time within which correlated events must occur (or not) (1..604800 seconds).
  * `events` - (Optional)(`ListOfString`) List of events that must occur within time interval.
  * `not_events` - (Optional)(`ListOfString`) List of events must not occur within time interval.
  * `trigger_count` - (Optional)(`Int`) Number of occurrences of triggering event.
  * `trigger_when` - (Optional)(`String`) To compare with `trigger_count`. Need to be 'after' (for event > count), 'on' (for event = count) or 'until' (for event < count).

### then arguments

* `change_configuration` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'change-configuration' configuration.
  * `commands` - (Required)(`ListOfString`) List of configuration commands.
  * `commit_options_check` - (Optional)(`Bool`) Check correctness of syntax; do not apply changes.
  * `commit_options_check_synchronize` - (Optional)(`Bool`) Synchronize commit check on both Routing Engines.
  * `commit_options_force` - (Optional)(`Bool`) Force commit on other Routing Engine (ignore warnings).
  * `commit_options_log` - (Optional)(`String`) Message to write to commit log.
  * `commit_options_synchronize` - (Optional)(`Bool`) Synchronize commit on both Routing Engines.
  * `retry_count` - (Optional)(`Int`) Change configuration retry attempt count (0..10).
  * `retry_interval` - (Optional)(`Int`) Time interval between each retry (seconds).
  * `user_name` - (Optional)(`String`) User under whose privileges configuration should be changed.
* `event_script` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Invoke event scripts. Can be specified multiple times for each filename. See the [`event_script` arguments for then](#event_script-arguments-for-then) block.
* `execute_commands` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Issue one or more CLI commands. Can be specified only once to declare 'execute-commands' configuration. See the [`execute_commands` arguments for then](#execute_commands-arguments-for-then) block.
* `ignore` - (Optional)(`Bool`) Do not log event or perform any other action. Conflict with other `then` arguments.
* `priority_override_facility` - (Optional)(`String`) Change syslog priority facility value.
* `priority_override_severity` - (Optional)(`String`) Change syslog priority severity value.
* `raise_trap` - (Optional)(`Bool`) Raise SNMP trap.
* `upload` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Upload file to specified destination. Can be specified multiple times for each combination of `filename` and `destination` arguments.
  * `filename` - (Optional)(`String`) Name of file to upload.
  * `destination` - (Optional)(`String`) Location to which to output file.
  * `retry_count` - (Optional)(`Int`) Upload output-filename retry attempt count (0..10).
  * `retry_interval` - (Optional)(`Int`) Time interval between each retry (seconds).
  * `transfer_delay` - (Optional)(`Int`) Delay before uploading file to the destination (seconds).
  * `user_name` - (Optional)(`String`) User under whose privileges upload action will execute.

### event_script arguments for then

* `filename` - (Required)(`String`) Local filename of the script file.
* `arguments` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Command line argument to the script. Can be specified multiple times for each combination of block arguments.
  * `name` - (Required)(`String`) Name of the argument.
  * `value` - (Required)(`String`) Value of the argument.
* `destination` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Location to which to upload event script output. Can be specified only once to declare 'destination' configuration.
  * `name` - (Required)(`String`) Destination name.
  * `retry_count` - (Optional)(`Int`) Upload output-filename retry attempt count (0..10).
  * `retry_interval` - (Optional)(`Int`) Time interval between each retry (seconds).
  * `transfer_delay` - (Optional)(`Int`) Delay before uploading files (seconds).
* `output_filename` - (Optional)(`String`) Name of file in which to write event script output.
* `output_format` - (Optional)(`String`) Format of output from event-script. Need to be 'text' or 'xml'.
* `user_name` - (Optional)(`String`) User under whose privileges event script will execute.

### execute_commands arguments for then

* `commands` - (Optional)(`ListOfString`) List of CLI commands to issue.
* `destination` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Location to which to upload command output. Can be specified only once to declare 'destination' configuration.
  * `name` - (Required)(`String`) Destination name.
  * `retry_count` - (Optional)(`Int`) Upload output-filename retry attempt count (0..10).
  * `retry_interval` - (Optional)(`Int`) Time interval between each retry (seconds).
  * `transfer_delay` - (Optional)(`Int`) Delay before uploading file to the destination (seconds).
* `output_filename` - (Optional)(`String`) Name of file in which to write command output.
* `output_format` - (Optional)(`String`) Format of output from CLI commands. Need to be 'text' or 'xml'.
* `user_name` - (Optional)(`String`) User under whose privileges command will execute.

## Import

Junos event-options policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_eventoptions_policy.demo demo
```
