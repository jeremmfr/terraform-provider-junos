---
layout: "junos"
page_title: "Junos: junos_firewall_filter"
sidebar_current: "docs-junos-resource-firewall-filter"
description: |-
  Create firewall filter
---

# junos_firewall_filter

Provides a firewall filter resource.

## Example Usage

```hcl
# Configure a firewall filter
resource junos_firewall_filter "filterdemo" {
  name   = "filterDemo"
  family = "inet"
  term {
    name = "filterDemo_term1"
    from {
      port        = ["22"]
      prefix_list = ["prefixList1"]
      protocol    = ["tcp"]
    }
    then {
      action = "accept"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of filter.
- **family** (Required, String, Forces new resource)  
  Family where create this filter.  
  Need to be `inet`, `inet6`, `any`, `ccc`, `mpls`, `vpls` or `ethernet-switching`.
- **interface_specific** (Optional, Boolean)  
  Defined counters are interface specific
- **term** (Required, Block List)  
  For each name of term.
  - **name** (Required, String)  
    Name of term.
  - **filter** (Optional, String)  
    Filter to include.
  - **from** (Required, Block)  
    Define match criteria.  
    See [below for nested schema](#from-arguments-for-term).
  - **then** (Required, Block)  
    Define action to take if the `from` condition is matched.  
    See [below for nested schema](#then-arguments-for-term).

---

### from arguments for term

- **address** (Optional, List of String)  
  Match IP source or destination address.
- **address_except** (Optional, List of String)  
  Match address not in this list of prefix.
- **destination_address** (Optional, List of String)  
  Match IP destination address.
- **destination_address_except** (Optional, List of String)  
  Match address not in this prefix.
- **destination_port** (Optional, List of String)  
  Match TCP/UDP destination port.
- **destination_port_except** (Optional, List of String)  
  Do not match TCP/UDP destination port.
- **destination_prefix_list** (Optional, List of String)  
  Match IP destination prefixes in named list.
- **destination_prefix_list_except** (Optional, List of String)  
  Match addresses not in this prefix list.
- **icmp_code** (Optional, List of String)  
  Match ICMP message code.
- **icmp_code_except** (Optional, List of String)  
  Do not match ICMP message code.
- **icmp_type** (Optional, List of String)  
  Match ICMP message type.
- **icmp_type_except** (Optional, List of String)  
  Do not match ICMP message type.
- **is_fragment** (Optional, Boolean)  
  Match if packet is a fragment.
- **next_header** (Optional, List of String)  
  Match next header protocol type.  
  Conflict with `next_header_except`.
- **next_header_except** (Optional, List of String)  
  Do not match next header protocol type.  
  Conflict with `next_header`.
- **port** (Optional, List of String)  
  Match TCP/UDP source or destination port.
- **port_except** (Optional, List of String)  
  Do not match TCP/UDP source or destination port.
- **prefix_list** (Optional, List of String)  
  Match IP source or destination prefixes in named list.
- **prefix_list_except** (Optional, List of String)  
  Match addresses not in this prefix list.
- **protocol** (Optional, List of String)  
  Match IP protocol type.
- **protocol_except** (Optional, List of String)  
  Do not match IP protocol type.
- **source_address** (Optional, List of String)  
  Match IP source address.
- **source_address_except** (Optional, List of String)  
  Match address not in this prefix.
- **source_port** (Optional, List of String)  
  Match TCP/UDP source port.
- **source_port_except** (Optional, List of String)  
  Do not match TCP/UDP source port.
- **source_prefix_list** (Optional, List of String)  
  Match IP source prefixes in named list.
- **source_prefix_list_except** (Optional, List of String)  
  Match addresses not in this prefix list.
- **tcp_established** (Optional, Boolean)  
  Match packet of an established TCP connection.
- **tcp_flags** (Optional, String)  
  Match TCP flags (in symbolic or hex formats).
- **tcp_initial** (Optional, Boolean)  
  Match initial packet of a TCP connection.

---

### then arguments for term

- **action** (Optional, String)  
  Action for term if needed.  
  Need to be `accept`, `reject`, `discard` or `next term`.
- **count** (Optional, String)  
  Count the packet in the named counter.
- **log** (Optional, Boolean)  
  Log the packet.
- **policer** (Optional, String)  
  Name of policer to use to rate-limit traffic.
- **port_mirror** (Optional, Boolean)  
  Port-mirror the packet.
- **routing_instance** (Optional, String)  
  Packets are directed to specified routing stance.
- **sample** (Optional, Boolean)  
  Sample the packet.
- **service_accounting** (Optional, Boolean)  
  Count the packets for service accounting.
- **syslog** (Optional, Boolean)  
  System log (syslog) information about the packet.

## Import

Junos firewall filter can be imported using an id made up of `<name>_-_<family>`, e.g.

```shell
$ terraform import junos_firewall_filter.filterdemo filterDemo_-_inet
```
