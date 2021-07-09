---
layout: "junos"
page_title: "Junos: junos_interface_physical"
sidebar_current: "docs-junos-data-source-interface-physical"
description: |-
  Get information on a physical interface (as with an junos_interface_physcial resource import)
---

# junos_interface_physical

Get information on a physical interface.

## Example Usage

```hcl
# Search interface with name
data junos_interface_physical "interface_physical_demo" {
  config_interface = "ge-0/0/3"
}
```

## Argument Reference

The following arguments are supported:

* `config_interface` - (Optional)(`String`) Specifies the interface part for search. Command is 'show configuration interfaces <config_interface>'
* `match` - (Optional)(`String`) Regex string to filter lines and find only one interface.

~> **NOTE:** If more or less than a single match is returned by the search, Terraform will fail.

## Attributes Reference

* `id` - Like resource it's the `name` of interface
* `name` - Name of physical interface (without dot).
* `ae_lacp` (**DEPRECATED**) - LACP option in aggregated-ether-options.
* `ae_link_speed` (**DEPRECATED**)  - Link speed of individual interface that joins the AE.
* `ae_minimum_links` (**DEPRECATED**) - Minimum number of aggregated links.
* `description` - Description for interface.
* `esi` - Define ESI Config parameters. See the [`esi` attributes](#esi-attributes) block.
* `ether_opts` - Declare 'ether-options' configuration.
  * `ae_8023ad` - Name of an aggregated Ethernet interface to join.
  * `auto_negotiation` - Enable auto-negotiation.
  * `no_auto_negotiation` - Don't enable auto-negotiation.
  * `flow_control` - Enable flow control.
  * `no_flow_control` - Don't enable flow control.
  * `loopback` - Enable loopback.
  * `no_loopback` - Don't enable loopback.
  * `redundant_parent` - Name of a redundant ethernet interface to join.
* `ether802_3ad` (**DEPRECATED**) - Link of 802.3ad interface.
* `gigether_opts` - Declare 'gigether-options' configuration.
  * `ae_8023ad` - Name of an aggregated Ethernet interface to join.
  * `auto_negotiation` - Enable auto-negotiation.
  * `no_auto_negotiation` - Don't enable auto-negotiation.
  * `flow_control` - Enable flow control.
  * `no_flow_control` - Don't enable flow control.
  * `loopback` - Enable loopback.
  * `no_loopback` - Don't enable loopback.
  * `redundant_parent` - Name of a redundant ethernet interface to join.
* `parent_ether_opts` - Declare 'aggregated-ether-options' or 'redundant-ether-options' configuration (it depends on the interface 'name'). See the [`parent_ether_opts` attributes](#parent_ether_opts-attributes) block.
* `trunk` - Interface mode is trunk.
* `vlan_members` - List of vlan membership for this interface.
* `vlan_native` - Vlan for untagged frames.
* `vlan_tagging` - 802.1q VLAN tagging support.

---

### esi attributes

* `mode` - ESI Mode
* `identifier` - The ESI value for the interface
* `auto_derive_lacp` - Auto-derive ESI value for the interface
* `df_election_type` - DF Election Type
* `source_bmac` - Unicast Source B-MAC address per ESI for PBB-EVPN

---

### parent_ether_opts attributes

* `bfd_liveness_detection` - Declare 'bfd-liveness-detection' in 'aggregated-ether-options' configuration. See the [`bfd_liveness_detection` attributes in parent_ether_opts](#bfd_liveness_detection-attributes-in-parent_ether_opts) block.
* `flow_control` - Enable flow control.
* `no_flow_control` - Don't enable flow control.
* `lacp` - Declare 'lacp' configuration.
  * `mode` - Active or passive.
  * `admin_key` - Node's administrative key.
  * `periodic` - Timer interval for periodic transmission of LACP packets.
  * `sync_reset` - On minimum-link failure notify out of sync to peer.
  * `system_id` - Node's System ID, encoded as a MAC address
  * `system_priority` - Priority of the system.
* `loopback` - Enable loopback.
* `no_loopback` - Don't enable loopback.
* `link_speed` - Link speed of individual interface that joins the AE.
* `minimum_bandwidth` - Minimum bandwidth configured for aggregated bundle.
* `minimum_links` - Minimum number of aggregated/active links.
* `redundancy_group` - Redundancy group of this interface for reth interface.
* `source_address_filter` - Source address filters.
* `source_filtering` - Enable source address filtering.

---

### bfd_liveness_detection attributes in parent_ether_opts

* `local_address` - BFD local address.
* `authentication_algorithm` - Authentication algorithm name.
* `authentication_key_chain` - Authentication Key chain name.
* `authentication_loose_check` - Verify authentication only if authentication is negotiated.
* `detection_time_threshold` - High detection-time triggering a trap.
* `holddown_interval` - Time to hold the session-UP notification to the client.
* `minimum_interval` - Minimum transmit and receive interval.
* `minimum_receive_interval` - Minimum receive interval.
* `multiplier` - Detection time multiplier.
* `neighbor` - BFD neighbor address.
* `no_adaptation` - Disable adaptation.
* `transmit_interval_minimum_interval` - Minimum transmit interval
* `transmit_interval_threshold` - High transmit interval triggering a trap.
* `version` - BFD protocol version number.
