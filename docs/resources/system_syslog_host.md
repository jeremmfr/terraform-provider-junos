---
page_title: "Junos: junos_system_syslog_host"
---

# junos_system_syslog_host

Configure a system syslog host.

## Example Usage

```hcl
# Add a system syslog host
resource "junos_system_syslog_host" "demo_syslog_host" {
  host = "192.0.2.1"
  port = 514
}
```

## Argument Reference

The following arguments are supported:

- **host** (Required, String, Forces new resource)  
  Host to be notified.
- **allow_duplicates** (Optional, Boolean)  
  Do not suppress the repeated message.
- **exclude_hostname** (Optional, Boolean)  
  Exclude hostname field in messages.
- **explicit_priority** (Optional, Boolean)  
  Include priority and facility in messages.
- **facility_override** (Optional, String)  
  Alternate facility for logging to remote host.
- **log_prefix** (Optional, String)  
  Prefix for all logging to this host.
- **match** (Optional, String)  
  Regular expression for lines to be logged.
- **match_strings** (Optional, List of String)  
  Matching string(s) for lines to be logged.
- **port** (Optional, Number)  
  Port number.
- **source_address** (Optional, String)  
  Use specified address as source address.
- **structured_data** (Optional, Block)  
  Log system message in structured format.
  - **brief** (Optional, Boolean)  
    Omit English-language text from end of logged message.
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
  An identifier for the resource with format `<host>`.

## Import

Junos system syslog host can be imported using an id made up of `<host>`, e.g.

```shell
$ terraform import junos_system_syslog_host.demo_syslog_host 192.0.2.1
```
