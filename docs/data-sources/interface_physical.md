---
page_title: "Junos: junos_interface_physical"
---

# junos_interface_physical

Get information on a physical interface  
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

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.
- **name** (String)  
  Name of physical interface (without dot).
- **ae_lacp** (String, **Deprecated**)  
  LACP option in aggregated-ether-options.
- **ae_link_speed** (String, **Deprecated**)  
  Link speed of individual interface that joins the AE.
- **ae_minimum_links** (Number, **Deprecated**)  
  Minimum number of aggregated links.
- **description** (String)  
  Description for interface.
- **disable** (Boolean)  
  Interface disabled.
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
- **ether802_3ad** (String, **Deprecated**)  
  Link of 802.3ad interface.
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
- **parent_ether_opts** (Block)  
  The `aggregated-ether-options` or `redundant-ether-options` configuration
  (it depends on the interface `name`).  
  See [below for nested schema](#parent_ether_opts-attributes).
- **trunk** (Boolean)  
  Interface mode is trunk.
- **vlan_members** (List of String)  
  List of vlan membership for this interface.
- **vlan_native** (Number)  
  Vlan for untagged frames.
- **vlan_tagging** (Boolean)  
  802.1q VLAN tagging support.

---

### esi attributes

- **mode** (String)  
  ESI Mode
- **identifier** (String)  
  The ESI value for the interface
- **auto_derive_lacp** (Boolean)  
  Auto-derive ESI value for the interface
- **df_election_type** (String)  
  DF Election Type
- **source_bmac** (String)  
  Unicast Source B-MAC address per ESI for PBB-EVPN

---

### parent_ether_opts attributes

- **bfd_liveness_detection** (Block)  
  The `bfd-liveness-detection` in `aggregated-ether-options` configuration.  
  See [below for nested schema](#bfd_liveness_detection-attributes-in-parent_ether_opts) block.
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
