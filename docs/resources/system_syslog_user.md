---
page_title: "Junos: junos_system_syslog_user"
---

# junos_system_syslog_user

Configure a system syslog user.

## Example Usage

```hcl
# Add a system syslog user
resource "junos_system_syslog_user" "all" {
  username     = "*"
  any_severity = "emergency"
}
resource "junos_system_syslog_user" "demo_syslog_user" {
  username        = "admin"
  kernel_severity = "any"
}
```

## Argument Reference

The following arguments are supported:

- **username** (Required, String, Forces new resource)  
  Name of user to notify (or `*` for all).
- **allow_duplicates** (Optional, Boolean)  
  Do not suppress the repeated message.
- **match** (Optional, String)  
  Regular expression for lines to be logged.
- **match_strings** (Optional, List of String)  
  Matching string(s) for lines to be logged.
- **any_severity** (Optional, String)  
  All facilities severity.
- **authorization_severity** (Optional, String)  
  Authorization system severity.
- **changelog_severity** (Optional, String)  
  Configuration change log severity.
- **conflictlog_severity** (Optional, String)  
  Configuration conflict log severity.
- **daemon_severity** (Optional, String)  
  Various system processes severity.
- **dfc_severity** (Optional, String)  
  Dynamic flow capture severity.
- **external_severity** (Optional, String)  
  Local external applications severity.
- **firewall_severity** (Optional, String)  
  Firewall filtering system severity.
- **ftp_severity** (Optional, String)  
  FTP process severity.
- **interactivecommands_severity** (Optional, String)  
  Commands executed by the UI severity.
- **kernel_severity** (Optional, String)  
  Kernel severity.
- **ntp_severity** (Optional, String)  
  NTP process severity.
- **pfe_severity** (Optional, String)  
  Packet Forwarding Engine severity.
- **security_severity** (Optional, String)  
  Security related severity.
- **user_severity** (Optional, String)  
  User processes severity.

**WARNING** All severities need to be
`alert`, `any`, `critical`, `emergency`, `error`, `info`, `none`, `notice` or `warning`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<username>`.

## Import

Junos system syslog host can be imported using an id made up of `<username>`, e.g.

```shell
$ terraform import junos_system_syslog_user.demo_syslog_user admin
```
