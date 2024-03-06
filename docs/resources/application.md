---
page_title: "Junos: junos_application"
---

# junos_application

Provides an application resource.

## Example Usage

```hcl
# Add an application
resource "junos_application" "mysql" {
  name             = "mysql"
  protocol         = "tcp"
  destination_port = "3306"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Application name.
- **application_protocol** (Optional, String)  
  Application protocol type.
- **description** (Optional, String)  
  Text description of application.
- **destination_port** (Optional, String)  
  Match TCP/UDP destination port.
- **ether_type** (Optional, String)  
  Match ether type.  
  Must be in hex (example: 0x8906).
- **inactivity_timeout** (Optional, Number)  
  Application-specific inactivity timeout (4..86400 seconds).  
  Conflict with `inactivity_timeout_never`.
- **inactivity_timeout_never** (Optional, Boolean)  
  Disables inactivity timeout.  
  Conflict with `inactivity_timeout`.
- **protocol** (Optional, String)  
  Match IP protocol type.
- **rpc_program_number** (Optional, String)  
  Match range of RPC program numbers.  
  Must be an integer or a range of integers.
- **source_port** (Optional, String)  
  Match TCP/UDP source port.
- **term** (Optional, Block List)  
  For each name of term to declare.  
  Conflict with `application_protocol`, `destination_port`, `inactivity_timeout`, `protocol`,
  `rpc_program_number`, `source_port` and `uuid`.  
  See [below for nested schema](#term-arguments).
- **uuid** (Optional, String)  
  Match universal unique identifier for DCE RPC objects.

### term arguments

- **name** (Required, String)  
  Term name.
- **protocol** (Required, String)  
  Match IP protocol type.
- **alg** (Optional, String)  
  Application Layer Gateway.
- **destination_port** (Optional, String)  
  Match TCP/UDP destination port.
- **icmp_code** (Optional, String)  
  Match ICMP message code.
- **icmp_type** (Optional, String)  
  Match ICMP message type.
- **icmp6_code** (Optional, String)  
  Match ICMP6 message code.
- **icmp6_type** (Optional, String)  
  Match ICMP6 message type.
- **inactivity_timeout** (Optional, Number)  
  Application-specific inactivity timeout (4..86400 seconds).  
  Conflict with `inactivity_timeout_never`.
- **inactivity_timeout_never** (Optional, Boolean)  
  Disables inactivity timeout.  
  Conflict with `inactivity_timeout`.
- **rpc_program_number** (Optional, String)  
  Match range of RPC program numbers.  
  Must be an integer or a range of integers.
- **source_port** (Optional, String)  
  Match TCP/UDP source port.
- **uuid** (Optional, String)  
  Match universal unique identifier for DCE RPC objects.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos application can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_application.mysql mysql
```
