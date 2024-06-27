---
page_title: "Junos: junos_security_log_stream"
---

# junos_security_log_stream

Provides a security log stream resource.

## Example Usage

```hcl
# Add a security log stream
resource "junos_security_log_stream" "demo_logstream" {
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

- **name** (Required, String, Forces new resource)  
  The name of security log stream.
- **category** (Optional, List of String)  
  Selects the type of events that may be logged.  
  Conflict with `filter_threat_attack`.
- **file** (Optional, Block)  
  Configure log file options for logs in local file.  
  Conflict with `host`.
  - **name** (Required, String)  
    Name of local log file.
  - **allow_duplicates** (Optional, Boolean)  
    To disable log consolidation.
  - **rotation** (Optional, Number)  
    Maximum number of rotate files (2..19).
  - **size** (Optional, Number)  
    Maximum size of local log file in megabytes (1..3).
- **filter_threat_attack** (Optional, Boolean)  
  Threat-attack security events are logged.  
  Conflict with `category`.
- **format** (Optional, String)  
  Specify the log stream format.  
  Need to be `binary`, `sd-syslog`, `syslog` or `welf`.
- **host** (Optional, Block)  
  Configure destination to send security logs to.  
  Conflict with `file`.
  - **ip_address** (Required, String)  
    IP address.
  - **port** (Optional, Number)  
    Host port number.
  - **routing_instance** (Optional, String)  
    Routing instance name.
- **rate_limit** (Optional, Number)  
  Rate-limit for security logs.
- **severity** (Optional, String)  
  Severity threshold for security logs.
- **transport** (Optional, Block)  
  Set security log transport settings.
  - **protocol** (Optional, String)  
    Set security log transport protocol for the device.  
    Need to be `tcp`, `tls` or `udp`.
  - **tcp_connections** (Optional, Number)  
    Set tcp connection number per-stream (1..5).
  - **tls_profile** (Optional, String)  
    TLS profile.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security log stream can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_log_stream.demo_logstream "demo_logstream"
```
