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

* `name` - (Required, Forces new resource)(`String`) Name of filter.
* `family` - (Required, Forces new resource)(`String`) Family where create this filter. </br>Need to be 'inet', 'inet6', 'any', 'ccc', 'mpls', 'vpls' or 'ethernet-switching'.
* `interface_specific` - (Optional)(`Bool`) Defined counters are interface specific
* `term` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each term.
  * `name` - (Required)(`String`) Name of term.
  * `filter` - (Optional)(`String`) Filter to include.
  * `from` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Define match criteria.
  See the [`from` arguments](#from-arguments) block. Max of 1.
  * `then` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Define action to take if the `from` condition is matched. See the [`then` arguments](#then-arguments) block. Max of 1.

#### `from` arguments
  * `address` - (Optional)(`ListOfString`) Match IP source or destination address.
  * `address_except` - (Optional)(`ListOfString`) Match address not in this list of prefix.
  * `port` - (Optional)(`ListOfString`) Match TCP/UDP source or destination port
  * `port_except` - (Optional)(`ListOfString`) Do not match TCP/UDP source or destination port.
  * `prefix_list` - (Optional)(`ListOfString`) Match IP source or destination prefixes in named list.
  * `prefix_list_except` - (Optional)(`ListOfString`) Match addresses not in this prefix list
  * `destination_address` - (Optional)(`ListOfString`) Match IP destination address
  * `destination_address_except` - (Optional)(`ListOfString`) Match address not in this prefix
  * `destination_port` - (Optional)(`ListOfString`) Match TCP/UDP destination port.
  * `destination_port_except` - (Optional)(`ListOfString`) Do not match TCP/UDP destination port.
  * `destination_prefix_list` - (Optional)(`ListOfString`) Match IP destination prefixes in named list.
  * `destination_prefix_list_except` - (Optional)(`ListOfString`) Match addresses not in this prefix list.
  * `source_address` - (Optional)(`ListOfString`) Match IP source address
  * `source_address_except` - (Optional)(`ListOfString`) Match address not in this prefix
  * `source_port` - (Optional)(`ListOfString`) Match TCP/UDP source port.
  * `source_port_except` - (Optional)(`ListOfString`) Do not match TCP/UDP source port.
  * `source_prefix_list` - (Optional)(`ListOfString`) Match IP source prefixes in named list.
  * `source_prefix_list_except` - (Optional)(`ListOfString`) Match addresses not in this prefix list.
  * `protocol` - (Optional)(`ListOfString`) Match IP protocol type.
  * `protocol_except` - (Optional)(`ListOfString`) Do not match IP protocol type.
  * `tcp_flags` - (Optional)(`String`) Match TCP flags (in symbolic or hex formats).
  * `tcp_initial` - (Optional)(`Bool`) Match initial packet of a TCP connection.
  * `tcp_established` - (Optional)(`Bool`) Match packet of an established TCP connection.
  * `icmp_type` - (Optional)(`ListOfString`) Match ICMP message type.
  * `icmp_type_except` - (Optional)(`ListOfString`) Do not match ICMP message type.

#### `then` arguments
  * `action` - (Optional)(`String`) Action for term if needed. Need to be 'accept', 'reject', 'discard' or 'next term'.
  * `count` - (Optional)(`String`) Count the packet in the named counter.
  * `routing_instance` - (Optional)(`String`) Packets are directed to specified routing stance.
  * `policer` - (Optional)(`String`) Name of policer to use to rate-limit traffic.
  * `log` - (Optional)(`Bool`) Log the packet.
  * `syslog` - (Optional)(`Bool`) System log (syslog) information about the packet.
  * `port_mirror` - (Optional)(`Bool`) Port-mirror the packet.
  * `sample` - (Optional)(`Bool`) Sample the packet.
  * `service_accounting` - (Optional)(`Bool`) Count the packets for service accounting.

## Import

Junos firewall filter can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_firewall_filter.filterdemo filterDemo
```
