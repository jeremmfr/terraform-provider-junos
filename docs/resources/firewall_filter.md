---
page_title: "Junos: junos_firewall_filter"
---

# junos_firewall_filter

Provides a firewall filter resource.

## Example Usage

```hcl
# Configure a firewall filter
resource "junos_firewall_filter" "filterdemo" {
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
  Filter name.
- **family** (Required, String, Forces new resource)  
  Family where create this filter.  
  Need to be `inet`, `inet6`, `any`, `ccc`, `mpls`, `vpls` or `ethernet-switching`.
- **interface_specific** (Optional, Boolean)  
  Defined counters are interface specific.
- **term** (Required, Block List)  
  For each name of term.
  - **name** (Required, String)  
    Term name.
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

- **address** (Optional, Set of String)  
  Match IP source or destination address.
- **address_except** (Optional, Set of String)  
  Match IP source or destination address not in this list of prefix.
- **destination_address** (Optional, Set of String)  
  Match IP destination address.
- **destination_address_except** (Optional, Set of String)  
  Match IP destination address not in this prefix.
- **destination_mac_address** (Optional, Set of String)  
  Destination MAC address.
- **destination_mac_address_except** (Optional, Set of String)  
  Destination MAC address not in this range.
- **destination_port** (Optional, Set of String)  
  Match TCP/UDP destination port.  
  Conflict with `destination_port_except`.
- **destination_port_except** (Optional, Set of String)  
  Do not match TCP/UDP destination port.  
  Conflict with `destination_port`.
- **destination_prefix_list** (Optional, Set of String)  
  Match IP destination prefixes in named list.
- **destination_prefix_list_except** (Optional, Set of String)  
  Match addresses not in this prefix list.
- **forwarding_class** (Optional, Set of String)  
  Match forwarding class.  
  Conflict with `forwarding_class_except`.
- **forwarding_class_except** (Optional, Set of String)  
  Do not match forwarding class.  
  Conflict with `forwarding_class`.
- **icmp_code** (Optional, Set of String)  
  Match ICMP message code.  
  Conflict with `icmp_code_except`.
- **icmp_code_except** (Optional, Set of String)  
  Do not match ICMP message code.  
  Conflict with `icmp_code`.
- **icmp_type** (Optional, Set of String)  
  Match ICMP message type.  
  Conflict with `icmp_type_except`.
- **icmp_type_except** (Optional, Set of String)  
  Do not match ICMP message type.  
  Conflict with `icmp_type`.
- **interface** (Optional, Set of String)  
  Match interface name.
- **is_fragment** (Optional, Boolean)  
  Match if packet is a fragment.
- **loss_priority** (Optional, Set of String)  
  Match Loss Priority.  
  Conflict with `loss_priority_except`.
- **loss_priority_except** (Optional, Set of String)  
  Do not match Loss Priority.  
  Conflict with `loss_priority`.
- **next_header** (Optional, Set of String)  
  Match next header protocol type.  
  Conflict with `next_header_except`.
- **next_header_except** (Optional, Set of String)  
  Do not match next header protocol type.  
  Conflict with `next_header`.
- **packet_length** (Optional, Set of String)  
  Match packet length.  
  Conflict with `packet_length_except`.
- **packet_length_except** (Optional, Set of String)  
  Do not match packet length.  
  Conflict with `packet_length`.
- **policy_map** (Optional, Set of String)  
  Match policy map.  
  Conflict with `policy_map_except`.
- **policy_map_except** (Optional, Set of String)  
  Do not match policy map.  
  Conflict with `policy_map`.
- **port** (Optional, Set of String)  
  Match TCP/UDP source or destination port.  
  Conflict with `port_except`.
- **port_except** (Optional, Set of String)  
  Do not match TCP/UDP source or destination port.  
  Conflict with `port`.
- **prefix_list** (Optional, Set of String)  
  Match IP source or destination prefixes in named list.
- **prefix_list_except** (Optional, Set of String)  
  Match addresses not in this prefix list.
- **protocol** (Optional, Set of String)  
  Match IP protocol type.  
  Conflict with `protocol_except`.
- **protocol_except** (Optional, Set of String)  
  Do not match IP protocol type.  
  Conflict with `protocol`.
- **source_address** (Optional, Set of String)  
  Match IP source address.
- **source_address_except** (Optional, Set of String)  
  Match IP source address not in this prefix.
- **source_mac_address** (Optional, Set of String)  
  Source MAC address.
- **source_mac_address_except** (Optional, Set of String)  
  Source MAC address not in this range.
- **source_port** (Optional, Set of String)  
  Match TCP/UDP source port.  
  Conflict with `source_port_except`.
- **source_port_except** (Optional, Set of String)  
  Do not match TCP/UDP source port.  
  Conflict with `source_port`.
- **source_prefix_list** (Optional, Set of String)  
  Match IP source prefixes in named list.
- **source_prefix_list_except** (Optional, Set of String)  
  Match IP source prefixes not in this prefix list.
- **tcp_established** (Optional, Boolean)  
  Match packet of an established TCP connection.  
  Conflict with `tcp_flags`.
- **tcp_flags** (Optional, String)  
  Match TCP flags (in symbolic or hex formats).  
  Conflict with `tcp_established` and `tcp_initial`.  
- **tcp_initial** (Optional, Boolean)  
  Match initial packet of a TCP connection.  
  Conflict with `tcp_flags`.

---

### then arguments for term

- **action** (Optional, String)  
  Action for term if needed.  
  Need to be `accept`, `reject`, `discard` or `next term`.
- **count** (Optional, String)  
  Count the packet in the named counter.
- **forwarding_class** (Optional, String)  
  Classify packet to forwarding class.
- **log** (Optional, Boolean)  
  Log the packet.
- **loss_priority** (Optional, String)  
  Packet's loss priority.
- **packet_mode** (Optional, Boolean)  
  Bypass flow mode for the packet.
- **policer** (Optional, String)  
  Name of policer to use to rate-limit traffic.
- **port_mirror** (Optional, Boolean)  
  Port-mirror the packet.
- **routing_instance** (Optional, String)  
  Packets are directed to specified routing instance.
- **sample** (Optional, Boolean)  
  Sample the packet.
- **service_accounting** (Optional, Boolean)  
  Count the packets for service accounting.
- **syslog** (Optional, Boolean)  
  System log (syslog) information about the packet.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<family>`.

## Import

Junos firewall filter can be imported using an id made up of `<name>_-_<family>`, e.g.

```shell
$ terraform import junos_firewall_filter.filterdemo filterDemo_-_inet
```
