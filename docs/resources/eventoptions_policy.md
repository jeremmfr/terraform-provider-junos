---
page_title: "Junos: junos_eventoptions_policy"
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

- **name** (Required, String, Forces new resource)  
  Name of policy.
- **events** (Required, Set of String)  
  List of events that trigger this policy.
- **then** (Required, Block)  
  Declare `then` configuration.  
  See [below for nested schema](#then-arguments).
- **attributes_match** (Optional, Block List)  
  For each combination of block arguments, attributes to compare for two events.
  - **from** (Required, String)  
    First attribute to compare.
  - **compare** (Required, String)  
    Type to compare.  
    Need to be `equals`, `matches` or `starts-with`.
  - **to** (Required, String)  
    Second attribute or value to compare.
- **within** (Optional, Block List)  
  For each time interval, list of events correlated with triggering events.
  - **time_interval** (Required, Number)  
    Time within which correlated events must occur (or not) (1..604800 seconds).
  - **events** (Optional, Set of String)  
    List of events that must occur within time interval.
  - **not_events** (Optional, Set of String)  
    List of events must not occur within time interval.
  - **trigger_count** (Optional, Number)  
    Number of occurrences of triggering event.
  - **trigger_when** (Optional, String)  
    To compare with `trigger_count`.  
    Need to be `after` (for event > count), `on` (for event = count) or `until` (for event < count).

### then arguments

- **change_configuration** (Optional, Block)  
  Declare `change-configuration` configuration.
  - **commands** (Required, List of String)  
    List of configuration commands.
  - **commit_options_check** (Optional, Boolean)  
    Check correctness of syntax; do not apply changes.
  - **commit_options_check_synchronize** (Optional, Boolean)  
    Synchronize commit check on both Routing Engines.
  - **commit_options_force** (Optional, Boolean)  
    Force commit on other Routing Engine (ignore warnings).
  - **commit_options_log** (Optional, String)  
    Message to write to commit log.
  - **commit_options_synchronize** (Optional, Boolean)  
    Synchronize commit on both Routing Engines.
  - **retry_count** (Optional, Number)  
    Change configuration retry attempt count (0..10).
  - **retry_interval** (Optional, Number)  
    Time interval between each retry (seconds).
  - **user_name** (Optional, String)  
    User under whose privileges configuration should be changed.
- **event_script** (Optional, Block List)  
  For each filename, invoke event scripts.  
  See [below for nested schema](#event_script-arguments-for-then).
- **execute_commands** (Optional, Block)  
  Declare `execute-commands` configuration.  
  Issue one or more CLI commands.  
  See [below for nested schema](#execute_commands-arguments-for-then).
- **ignore** (Optional, Boolean)  
  Do not log event or perform any other action.  
  Conflict with other `then` arguments.
- **priority_override_facility** (Optional, String)  
  Change syslog priority facility value.
- **priority_override_severity** (Optional, String)  
  Change syslog priority severity value.
- **raise_trap** (Optional, Boolean)  
  Raise SNMP trap.
- **upload** (Optional, Block List)  
  For each combination of `filename` and `destination` arguments, upload file to specified destination.
  - **filename** (Optional, String)  
    Name of file to upload.
  - **destination** (Optional, String)  
    Location to which to output file.
  - **retry_count** (Optional, Number)  
    Upload output-filename retry attempt count (0..10).
  - **retry_interval** (Optional, Number)  
    Time interval between each retry (seconds).
  - **transfer_delay** (Optional, Number)  
    Delay before uploading file to the destination (seconds).
  - **user_name** (Optional, String)  
    User under whose privileges upload action will execute.

### event_script arguments for then

- **filename** (Required, String)  
  Local filename of the script file.
- **arguments** (Optional, Block List)  
  For each name of arguments, command line argument to the script.
  - **name** (Required, String)  
    Name of the argument.
  - **value** (Required, String)  
    Value of the argument.
- **destination** (Optional, Block)  
  Declare `destination` configuration.  
  Location to which to upload event script output.
  - **name** (Required, String)  
    Destination name.
  - **retry_count** (Optional, Number)  
    Upload output-filename retry attempt count (0..10).
  - **retry_interval** (Optional, Number)  
    Time interval between each retry (seconds).
  - **transfer_delay** (Optional, Number)  
    Delay before uploading files (seconds).
- **output_filename** (Optional, String)  
  Name of file in which to write event script output.
- **output_format** (Optional, String)  
  Format of output from event-script.  
  Need to be `text` or `xml`.
- **user_name** (Optional, String)  
  User under whose privileges event script will execute.

### execute_commands arguments for then

- **commands** (Optional, List of String)  
  List of CLI commands to issue.
- **destination** (Optional, Block)  
  Declare `destination` configuration.  
  Location to which to upload command output.
  - **name** (Required, String)  
    Destination name.
  - **retry_count** (Optional, Number)  
    Upload output-filename retry attempt count (0..10).
  - **retry_interval** (Optional, Number)  
    Time interval between each retry (seconds).
  - **transfer_delay** (Optional, Number)  
    Delay before uploading file to the destination (seconds).
- **output_filename** (Optional, String)  
  Name of file in which to write command output.
- **output_format** (Optional, String)  
  Format of output from CLI commands.  
  Need to be `text` or `xml`.
- **user_name** (Optional, String)  
  User under whose privileges command will execute.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos event-options policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_eventoptions_policy.demo demo
```
