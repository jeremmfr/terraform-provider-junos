---
layout: "junos"
page_title: "Junos: junos_system_syslog_file"
sidebar_current: "docs-junos-resource-system-syslog-file"
description: |-
  Configure a system syslog file
---

# junos_system_syslog_file

Configure a system syslog file.

## Example Usage

```hcl
# Add a system syslog file
resource junos_system_syslog_file "demo_syslog_file" {
  filename     = "demo"
  any_severity = "emergency"
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required, Forces new resource)(`String`) Name of file in which to log data.
* `allow_duplicates` - (Optional)(`Bool`) Do not suppress the repeated message.
* `explicit_priority`- (Optional)(`Bool`) Include priority and facility in messages.
* `match` - (Optional)(`String`) Regular expression for lines to be logged.
* `match_strings` - (Optional)(`ListOfString`) Matching string(s) for lines to be logged.
* `structured_data` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Log system message in structured format. Max of 1.
  * `brief` - (Optional)(`Bool`) Omit English-language text from end of logged message.
* `archive` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Define parameters for archiving log messages. Max of 1. See the [`archive` arguments] (#archive-arguments) block.
* `any_severity` - (Optional)(`String`) All facilities severity.
* `authorization_severity` - (Optional)(`String`) Authorization system severity.
* `changelog_severity` - (Optional)(`String`) Configuration change log severity.
* `conflictlog_severity` - (Optional)(`String`) Configuration conflict log severity.
* `daemon_severity` - (Optional)(`String`) Various system processes severity.
* `dfc_severity` - (Optional)(`String`) Dynamic flow capture severity.
* `external_severity` - (Optional)(`String`) Local external applications severity.
* `firewall_severity` - (Optional)(`String`) Firewall filtering system severity.
* `ftp_severity` - (Optional)(`String`) FTP process severity.
* `interactivecommands_severity` - (Optional)(`String`) Commands executed by the UI severity.
* `kernel_severity` - (Optional)(`String`) Kernel severity.
* `ntp_severity` - (Optional)(`String`) NTP process severity.
* `pfe_severity` - (Optional)(`String`) Packet Forwarding Engine severity.
* `security_severity` - (Optional)(`String`) Security related severity.
* `user_severity` - (Optional)(`String`) User processes severity.

#### archive arguments
* `binary_data` - (Optional)(`Bool`) Mark file as if it contains binary data. Conflict with `no_binary_data`.
* `no_binary_data` - (Optional)(`Bool`) Don't mark file as if it contains binary data. Conflict with `binary_data`.
* `world_readable` - (Optional)(`Bool`) Allow any user to read the log file. Conflict with `no_world_readable`.
* `no_world_readable` - (Optional)(`Bool`) Don't allow any user to read the log file. Conflict with `world_readable`.
* `files` - (Optional)(`Int`) Number of files to be archived (1..1000)
* `size` - (Optional)(`Int`) Size of files to be archived (65536..1073741824 bytes)
* `start_time` - (Optional)(`String`) Start time for file transmission (yyyy-mm-dd.hh:mm)
* `transfer_interval` - (Optional)(`Int`) Frequency at which to transfer files to archive sites (5..2880 minutes)
* `sites` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure an archive site. Can be specified multiple times for each url (first declaration is primary URL, failover for others). 
  * `url` - (Required)(`String`) Primary or failover URLs to receive archive files.
  * `password` - (Optional)(`String`) Password for login into the archive site.
  **WARNING** Clear in tfstate.
  * `routing_instance` - (Optional)(`String`) Routing instance.

**WARNING** All severities need to be 'alert', 'any', 'critical', 'emergency', 'error', 'info', 'none', 'notice' or 'warning'.

## Import

Junos system syslog file can be imported using an id made up of `<filename>`, e.g.

```
$ terraform import junos_system_syslog_file.demo_syslog_file demo
```
