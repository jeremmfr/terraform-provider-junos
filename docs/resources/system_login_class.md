---
page_title: "Junos: junos_system_login_class"
---

# junos_system_login_class

Provides a system login class resource.

## Example Usage

```hcl
# Add a system login class
resource "junos_system_login_class" "engineering" {
  name         = "engineering"
  idle_timeout = 30
  login_alarms = true
  permissions  = ["all"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of system login class.
- **access_end** (Optional, String)  
  End time for remote access (HH:MM:SS).
- **access_start** (Optional, String)  
  Start time for remote access (HH:MM:SS).
- **allow_commands** (Optional, String)  
  Regular expression for commands to allow explicitly.  
  Conflict with `allow_commands_regexps`.
- **allow_commands_regexps** (Optional, List of String)  
  Object path regular expressions to allow commands.  
  Conflict with `allow_commands`.
- **allow_configuration** (Optional, String)  
  Regular expression for configure to allow explicitly.  
  Conflict with `allow_configuration_regexps`.
- **allow_configuration_regexps** (Optional, List of String)  
  Object path regular expressions to allow.  
  Conflict with `allow_configuration`.
- **allow_hidden_commands** (Optional, Boolean)  
  Allow all hidden commands to be executed.  
  Conflict with `no_hidden_commands_except`.
- **allowed_days** (Optional, List of String)  
  Day(s) of week when access is allowed.
- **cli_prompt** (Optional, String)  
  Cli prompt name for this class.
- **configuration_breadcrumbs** (Optional, Boolean)  
  Enable breadcrumbs during display of configuration.
- **confirm_commands** (Optional, List of String)  
  List of commands to be confirmed explicitly.
- **deny_commands** (Optional, String)  
  Regular expression for commands to deny explicitly.  
  Conflict with `deny_commands_regexps`.
- **deny_commands_regexps** (Optional, List of String)  
  Object path regular expressions to deny commands.  
  Conflict with `deny_commands`.
- **deny_configuration** (Optional, String)  
  Regular expression for configure to deny explicitly.  
  Conflict with `deny_configuration_regexps`.
- **deny_configuration_regexps** (Optional, List of String)  
  Object path regular expressions to deny.  
  Conflict with `deny_configuration`.
- **idle_timeout** (Optional, Number)  
  Maximum idle time before logout (minutes).
- **logical_system** (Optional, String)  
  Logical system associated with login.
- **login_alarms** (Optional, Boolean)  
  Display system alarms when logging in.
- **login_script** (Optional, String)  
  Execute this login-script when logging in.
- **login_tip** (Optional, Boolean)  
  Display tip when logging in.
- **no_hidden_commands_except** (Optional, List of String)  
  Deny all hidden commands with exemptions.
- **permissions** (Optional, Set of String)  
  Set of permitted operation categories.
- **security_role** (Optional, String)  
  Common Criteria security role.
- **tenant** (Optional, String)  
  Tenant associated with this login.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos system login class can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_system_login_class.engineering engineering
```
