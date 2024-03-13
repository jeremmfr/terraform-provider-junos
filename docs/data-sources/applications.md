---
page_title: "Junos: junos_applications"
---

# junos_applications

Get list of filtered applications on the Junos device  
(in `applications` and `group junos-defaults applications` level).

## Example Usage

```hcl
# Find default application junos-ssh 
data "junos_applications" "default_ssh" {
  match_name = "^junos-"
  match_options {
    protocol         = "tcp"
    destination_port = 22
  }
}
```

## Argument Reference

The following arguments are supported:

- **match_name** (Optional, String)  
  A regexp to apply a filter on applications name.  
  Need to be a valid regexp.
- **match_options** (Optional, Block Set)  
  List of options to apply a filter on applications.  
  Use multiple blocks to match applications with multiple terms.
  - **alg** (Optional, String)  
    Application Layer Gateway.
  - **application_protocol** (Optional, String)  
    Application protocol type.
  - **destination_port** (Optional, String)  
    Match TCP/UDP destination port.
  - **ether_type** (Optional, String)  
    Match ether type.
  - **icmp_code** (Optional, String)  
    Match ICMP message code.
  - **icmp_type** (Optional, String)  
    Match ICMP message type.
  - **icmp6_code** (Optional, String)  
    Match ICMP6 message code.
  - **icmp6_type** (Optional, String)  
    Match ICMP6 message type.
  - **inactivity_timeout** (Optional, Number)  
    Application-specific inactivity timeout.
  - **inactivity_timeout_never** (Optional, Boolean)  
    Disables inactivity timeout.
  - **protocol** (Optional, String)  
    Match IP protocol type.
  - **rpc_program_number** (Optional, String)  
    Match range of RPC program numbers.
  - **source_port** (Optional, String)  
    Match TCP/UDP source port.
  - **uuid** (Optional, String)  
    Match universal unique identifier for DCE RPC objects.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source.
- **applications** (Block List)  
  For each application found.
  - **name** (String)  
    Name of application.
  - **application_protocol** (String)  
    Application protocol type.
  - **description** (String)  
    Text description of application.
  - **destination_port** (String)  
    Port(s) destination used by application.
  - **ether_type** (String)  
    Match ether type.
  - **inactivity_timeout** (Number)  
    Application-specific inactivity timeout (4..86400 seconds).
  - **inactivity_timeout_never** (Boolean)  
    Disables inactivity timeout.
  - **protocol** (String)  
    Protocol used by application.
  - **rpc_program_number** (String)  
    Match range of RPC program numbers.
  - **source_port** (String)  
    Port(s) source used by application.
  - **term** (Block List)  
    For each name of term.  
    See [below for nested schema](#term-attributes).
  - **uuid** (String)  
    Match universal unique identifier for DCE RPC objects.

### term attributes

- **name** (String)  
  Term name.
- **alg** (String)  
  Application Layer Gateway.
- **destination_port** (String)  
  Match TCP/UDP destination port.
- **icmp_code** (String)  
  Match ICMP message code.
- **icmp_type** (String)  
  Match ICMP message type.
- **icmp6_code** (String)  
  Match ICMP6 message code.
- **icmp6_type** (String)  
  Match ICMP6 message type.
- **inactivity_timeout** (Number)  
  Application-specific inactivity timeout (4..86400 seconds).
- **inactivity_timeout_never** (Boolean)  
  Disables inactivity timeout.
- **protocol** (String)  
  Match IP protocol type.
- **rpc_program_number** (String)  
  Match range of RPC program numbers.
- **source_port** (String)  
  Match TCP/UDP source port.
- **uuid** (String)  
  Match universal unique identifier for DCE RPC objects.
