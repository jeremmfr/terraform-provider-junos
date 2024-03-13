---
page_title: "Junos: junos_interface_physical"
---

# junos_interface_physical

Get configuration from a physical interface  
(as with a junos_interface_physical resource import).

## Example Usage

```hcl
# Search interface with name
data "junos_interface_physical" "interface_physical_demo" {
  config_interface = "ge-0/0/3"
}
```

## Argument Reference

The following arguments are supported:

- **config_interface** (Optional, String)  
  Specifies the interface part for search.  
  Command is `show configuration interfaces <config_interface>`
- **match** (Optional, String)  
  Regex string to filter lines and find only one interface.

~> **NOTE:** If more or less than a single match is returned by the search, Terraform will fail.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **name** (String)  
  Name of physical interface (without dot).
- **description** (String)  
  Description for interface.
- **disable** (Boolean)  
  Interface disabled.
- **encapsulation** (String)  
  Physical link-layer encapsulation.
- **esi** (Block)  
  ESI Config parameters.  
  See [below for nested schema](#esi-attributes).
- **ether_opts** (Block)  
  The `ether-options` configuration.
  - **ae_8023ad** (String)  
    Name of an aggregated Ethernet interface to join.
  - **auto_negotiation** (Boolean)  
    Enable auto-negotiation.
  - **no_auto_negotiation** (Boolean)  
    Don't enable auto-negotiation.
  - **flow_control** (Boolean)  
    Enable flow control.
  - **no_flow_control** (Boolean)  
    Don't enable flow control.
  - **loopback** (Boolean)  
    Enable loopback.
  - **no_loopback** (Boolean)  
    Don't enable loopback.
  - **redundant_parent** (String)  
    Name of a redundant ethernet interface to join.
- **flexible_vlan_tagging** (Boolean)  
  Support for no tagging, or single and double 802.1q VLAN tagging.
- **gigether_opts** (Block)  
  The `gigether-options` configuration.
  - **ae_8023ad** (String)  
    Name of an aggregated Ethernet interface to join.
  - **auto_negotiation** (Boolean)  
    Enable auto-negotiation.
  - **no_auto_negotiation** (Boolean)  
    Don't enable auto-negotiation.
  - **flow_control** (Boolean)  
    Enable flow control.
  - **no_flow_control** (Boolean)  
    Don't enable flow control.
  - **loopback** (Boolean)  
    Enable loopback.
  - **no_loopback** (Boolean)  
    Don't enable loopback.
  - **redundant_parent** (String)  
    Name of a redundant ethernet interface to join.
- **gratuitous_arp_reply** (Boolean)  
  Enable gratuitous ARP reply.
- **hold_time_down** (Number)  
  Link down hold time (milliseconds).
- **hold_time_up** (Number)  
  Link up hold time (milliseconds).
- **link_mode** (String)  
  Link operational mode.
- **mtu** (Number)  
  Maximum transmission unit.
- **no_gratuitous_arp_reply** (Boolean)  
  Don't enable gratuitous ARP reply.
- **no_gratuitous_arp_request** (Boolean)  
  Ignore gratuitous ARP request.
- **parent_ether_opts** (Block)  
  The `aggregated-ether-options` or `redundant-ether-options` configuration
  (it depends on the interface `name`).  
  See [below for nested schema](#parent_ether_opts-attributes).
- **speed** (String)  
  Link speed.
- **storm_control** (String)  
  Storm control profile name to bind.
- **trunk** (Boolean)  
  Interface mode is trunk.
- **trunk_non_els** (Boolean)  
  Port mode is trunk.  
  When use `port-mode` instead of `interface-mode` on non-ELS devices.
- **vlan_members** (List of String)  
  List of vlan membership for this interface.
- **vlan_native** (Number)  
  Vlan for untagged frames.
- **vlan_native_non_els** (String)  
  Vlan for untagged frames.  
  When use `native-vlan-id` in `unit 0 family ethernet-switching`
  instead of interface root level on non-ELS devices.
- **vlan_tagging** (Boolean)  
  802.1q VLAN tagging support.

---

### esi attributes

- **mode** (String)  
  ESI Mode.
- **auto_derive_lacp** (Boolean)  
  Auto-derive ESI value for the interface.
- **df_election_type** (String)  
  DF Election Type.
- **identifier** (String)  
  The ESI value for the interface.
- **source_bmac** (String)  
  Unicast Source B-MAC address per ESI for PBB-EVPN.

---

### parent_ether_opts attributes

- **bfd_liveness_detection** (Block)  
  The `bfd-liveness-detection` in `aggregated-ether-options` configuration.  
  See [below for nested schema](#bfd_liveness_detection-attributes-in-parent_ether_opts).
- **flow_control** (Boolean)  
  Enable flow control.
- **no_flow_control** (Boolean)  
  Don't enable flow control.
- **lacp** (Block)  
  The `lacp` configuration.
  - **mode** (String)  
    Active or passive.
  - **admin_key** (Number)  
    Node's administrative key.
  - **periodic** (String)  
    Timer interval for periodic transmission of LACP packets.
  - **sync_reset** (String)  
    On minimum-link failure notify out of sync to peer.
  - **system_id** (String)  
    Node's System ID, encoded as a MAC address
  - **system_priority** (Number)  
    Priority of the system.
- **loopback** (Boolean)  
  Enable loopback.
- **no_loopback** (Boolean)  
  Don't enable loopback.
- **link_speed** (String)  
  Link speed of individual interface that joins the AE.
- **mc_ae** (Block)  
  Multi-chassis aggregation (MC-AE) network device configuration.  
  See [below for nested schema](#mc_ae-arguments-in-parent_ether_opts).
- **minimum_bandwidth** (String)  
  Minimum bandwidth configured for aggregated bundle.
- **minimum_links** (Number)  
  Minimum number of aggregated/active links.
- **redundancy_group** (Number)  
  Redundancy group of this interface for reth interface.
- **source_address_filter** (List of String)  
  Source address filters.
- **source_filtering** (Boolean)  
  Enable source address filtering.

---

### bfd_liveness_detection attributes in parent_ether_opts

- **local_address** (String)  
  BFD local address.
- **authentication_algorithm** (String)  
  Authentication algorithm name.
- **authentication_key_chain** (String)  
  Authentication Key chain name.
- **authentication_loose_check** (Boolean)  
  Verify authentication only if authentication is negotiated.
- **detection_time_threshold** (Number)  
  High detection-time triggering a trap.
- **holddown_interval** (Number)  
  Time to hold the session-UP notification to the client.
- **minimum_interval** (Number)  
  Minimum transmit and receive interval.
- **minimum_receive_interval** (Number)  
  Minimum receive interval.
- **multiplier** (Number)  
  Detection time multiplier.
- **neighbor** (String)  
  BFD neighbor address.
- **no_adaptation** (Boolean)  
  Disable adaptation.
- **transmit_interval_minimum_interval** (Number)  
  Minimum transmit interval
- **transmit_interval_threshold** (Number)  
  High transmit interval triggering a trap.
- **version** (String)  
  BFD protocol version number.

---

### mc_ae arguments in parent_ether_opts

- **chassis_id** (Number)  
  Chassis id of MC-AE network device.
- **mc_ae_id** (Number)  
  MC-AE group id.
- **mode** (String)  
  Mode of the MC-AE.
- **status_control** (String)  
  Status of the MC-AE chassis.
- **enhanced_convergence** (Boolean)  
  Optimized convergence time for MC-AE.
- **events_iccp_peer_down** (Block)  
  Define behavior in the event of ICCP peer down.  
  - **force_icl_down** (Boolean)  
    Bring down ICL logical interface.
  - **prefer_status_control_active** (Boolean)  
    Keep this node up (recommended only on status-control active).
- **init_delay_time** (Number)  
  Init delay timer for mcae sm for min traffic loss (seconds).
- **redundancy_group** (Number)  
  Redundancy group id.
- **revert_time** (Number)  
  Wait interval before performing switchover (minute).
- **switchover_mode** (String)  
  Switchover mode.
