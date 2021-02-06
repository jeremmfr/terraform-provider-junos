---
layout: "junos"
page_title: "Junos: junos_security_log_stream"
sidebar_current: "docs-junos-resource-security-log-stream"
description: |-
  Create a security log stream (when Junos device supports it)
---

# junos_security_log_stream

Provides a security log stream resource.

## Example Usage

```hcl
# Add a security log stream
resource junos_security_log_stream "demo_logstream" {
  name     = "demo_logstream"
  category = ["idp", "screen"]
  format   = "sd-syslog"
  host {
    ip_address = "192.0.2.10"
    port       = 514
  }
  severity = "info"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of security log stream.
* `category` - (Optional)(`ListOfString`) Selects the type of events that may be logged. Conflict with `filter_threat_attack`.
* `file` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure log file options for logs in local file. Max of 1. Conflict with `host`.
  * `name` - (Required)(`String`) Name of local log file.
  * `allow_duplicates` - (Optional)(`Bool`) To disable log consolidation.
  * `rotation` - (Optional)(`Int`) Maximum number of rotate files (2..19).
  * `size` - (Optional)(`Int`) Maximum size of local log file in megabytes (1..3).
* `filter_threat_attack` - (Optional)(`Bool`) Threat-attack security events are logged. Conflict with `category`.
* `format` - (Optional)(`String`) Specify the log stream format. Need to be 'binary', 'sd-syslog', 'syslog' or 'welf'.
* `host` - Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure destination to send security logs to. Max of 1. Conflict with `file`.
  * `ip_address` - (Required)(`String`) IP address.
  * `port` - (Optional)(`Int`) Host port number.
  * `routing_instance` - (Optional)(`String`) Routing-instance name.
* `rate_limit` - (Optional)(`Int`) Rate-limit for security logs.
* `severity` - (Optional)(`String`) Severity threshold for security logs.

## Import

Junos security log stream can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_log_stream.demo_logstream "demo_logstream"
```
