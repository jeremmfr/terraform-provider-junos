---
page_title: "Junos: junos_security_screen"
---

# junos_security_screen

Provides a security screen resource.

## Example Usage

```hcl
# Add a security screen
resource "junos_security_screen" "demo_screen" {
  name               = "demo_screen"
  alarm_without_drop = true
  description        = "desc screen"
  icmp {
    flood {}
    ping_death = true
  }
  ip {
    spoofing = true
  }
  limit_session {
    destination_ip_based = 2000
    source_ip_based      = 3000
  }
  tcp {
    syn_flood {}
  }
  udp {
    flood {}
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of screen.
- **alarm_without_drop** (Optional, Boolean)  
  Do not drop packet, only generate alarm.
- **description** (Optional, String)  
  Text description of screen.
- **icmp** (Optional, Block)  
  Declare `icmp` configuration.  
  See [below for nested schema](#icmp-arguments).
- **ip** (Optional, Block)  
  Declare `ip` configuration.  
  See [below for nested schema](#ip-arguments).
- **limit_session** (Optional, Block)  
  Declare `limit-session` configuration.
  - **destination_ip_based** (Optional, Number)  
    Limit sessions to the same destination IP (1..2000000).
  - **source_ip_based** (Optional, Number)  
    Limit sessions from the same source IP (1..2000000).
- **tcp** (Optional, Block)  
  Declare `tcp` configuration.  
  See [below for nested schema](#tcp-arguments).
- **udp** (Optional, Block)  
  Declare `udp` configuration.  
  See [below for nested schema](#udp-arguments).

---

### icmp arguments

- **flood** (Optional, Block)  
  Enable icmp flood ids option.
  - **threshold** (Optional, Number)  
    Threshold (1..1000000 ICMP packets per second).
- **fragment** (Optional, Boolean)  
  Enable ICMP fragment ids option.
- **icmpv6_malformed** (Optional, Boolean)  
  Enable icmpv6 malformed ids option
- **large** (Optional, Boolean)  
  Enable large ICMP packet (size > 1024) ids option
- **ping_death** (Optional, Boolean)  
  Enable ping of death ids option
- **sweep** (Optional, Block)  
  Enable ip sweep ids option.
  - **threshold** (Optional, Number)  
    Threshold (1000..1000000 microseconds in which 10 ICMP packets are detected).

---

### ip arguments

- **bad_option** (Optional, Boolean)  
  Enable ip with bad option ids option.
- **block_frag** (Optional, Boolean)  
  Enable ip fragment blocking ids option.
- **ipv6_extension_header** (Optional, Block)  
  Declare `ip ipv6-extension-header` configuration.  
  See [below for nested schema](#ipv6_extension_header-arguments-for-ip).
- **ipv6_extension_header_limit** (Optional, Number)  
  Enable ipv6 extension header limit ids option (0..32).
- **ipv6_malformed_header** (Optional, Boolean)  
  Enable ipv6 malformed header ids option.
- **loose_source_route_option** (Optional, Boolean)  
  Enable ip with loose source route ids option.
- **record_route_option** (Optional, Boolean)  
  Enable ip with record route option ids option.
- **security_option** (Optional, Boolean)  
  Enable ip with security option ids option.
- **source_route_option** (Optional, Boolean)  
  Enable ip source route ids option.
- **spoofing** (Optional, Boolean)  
  Enable ip address spoofing ids option.
- **stream_option** (Optional, Boolean)  
  Enable ip with stream option ids option.
- **strict_source_route_option** (Optional, Boolean)  
  Enable ip with strict source route ids option.
- **tear_drop** (Optional, Boolean)  
  Enable tear drop ids option.
- **timestamp_option** (Optional, Boolean)  
  Enable ip with timestamp option ids option.
- **tunnel** (Optional, Block)  
  Declare `ip tunnel` configuration.  
  See [below for nested schema](#tunnel-arguments-for-ip).
- **unknown_protocol** (Optional, Boolean)  
  Enable ip unknown protocol ids option.

---

### tcp arguments

- **fin_no_ack** (Optional, Boolean)  
  Enable Fin bit with no ACK bit ids option.
- **land** (Optional, Boolean)  
  Enable land attack ids option.
- **no_flag** (Optional, Boolean)  
  Enable TCP packet without flag ids option.
- **port_scan** (Optional, Block)  
  Enable TCP port scan ids option.
  - **threshold** (Optional, Number)  
    Threshold (1000..1000000 microseconds in which 10 attack packets are detected).
- **syn_ack_ack_proxy** (Optional, Block)  
  Enable syn-ack-ack proxy ids option.
  - **threshold** (Optional, Number)  
    Threshold (1..250000 un-authenticated connections).
- **syn_fin** (Optional, Boolean)  
  Enable SYN and FIN bits set attack ids option.
- **syn_flood** (Optional, Block)  
  Enable SYN flood ids option.  
  See [below for nested schema](#syn_flood-arguments-for-tcp).
- **syn_frag** (Optional, Boolean)  
  Enable SYN fragment ids option.
- **sweep** (Optional, Block)  
  Enable TCP sweep ids option.
  - **threshold** (Optional, Number)  
    Threshold (1000..1000000 microseconds in which 10 TCP packets are detected).
- **winnuke** (Optional, Boolean)  
  Enable winnuke attack ids option.

---

### udp arguments

- **flood** (Optional, Block)  
  UDP flood ids option.
  - **threshold** (Optional, Number)  
    Threshold (1..1000000 UDP packets per second).
  - **whitelist** (Optional, Set of String)  
    List of UDP flood white list group name.
- **port_scan** (Optional, Block)  
  UDP port scan ids option.
  - **threshold** (Optional, Number)  
    Threshold (1000..1000000 microseconds in which 10 attack packets are detected).
- **sweep** (Optional, Block)  
  UDP sweep ids option.
  - **threshold** (Optional, Number)  
    Threshold (1000..1000000 microseconds in which 10 UDP packets are detected).

---

### ipv6_extension_header arguments for ip

- **ah_header** (Optional, Boolean)  
  Enable ipv6 Authentication Header ids option.
- **esp_header** (Optional, Boolean)  
  Enable ipv6 Encapsulating Security Payload header ids option.
- **hip_header** (Optional, Boolean)  
  Enable ipv6 Host Identify Protocol header ids option.
- **destination_header** (Optional, Block)  
  Enable ipv6 destination option header ids option.
  - **ilnp_nonce_option** (Optional, Boolean)  
    Enable Identifier-Locator Network Protocol Nonce option ids option.
  - **home_address_option** (Optional, Boolean)  
    Enable home address option ids option.
  - **line_identification_option** (Optional, Boolean)  
    Enable line identification option ids option.
  - **tunnel_encapsulation_limit_option** (Optional, Boolean)  
    Enable tunnel encapsulation limit option ids option.
  - **user_defined_option_type** (Optional, List of String)  
    User-defined option type range.  
    Need to be `(1..255)` or `(1..255) to (1..255)`.
- **fragment_header** (Optional, Boolean)  
  Enable ipv6 fragment header ids option.
- **hop_by_hop_header** (Optional, Block)  
  Enable ipv6 hop by hop option header ids option.
  - **calipso_option** (Optional, Boolean)  
    Enable Common Architecture Label ipv6 Security Option ids option.
  - **rpl_option** (Optional, Boolean)  
    Enable Routing Protocol for Low-power and Lossy networks option ids option.
  - **smf_dpd_option** (Optional, Boolean)  
    Enable Simplified Multicast Forwarding ipv6 Duplicate Packet Detection option ids option.
  - **jumbo_payload_option** (Optional, Boolean)  
    Enable jumbo payload option ids option.
  - **quick_start_option** (Optional, Boolean)  
    Enable quick start option ids option.
  - **router_alert_option** (Optional, Boolean)  
    Enable router alert option ids option.
  - **user_defined_option_type** (Optional, List of String)  
    User-defined option type range.  
    Need to be `(1..255)` or `(1..255) to (1..255)`.
- **mobility_header** (Optional, Boolean)  
  Enable ipv6 mobility header ids option.
- **no_next_header** (Optional, Boolean)  
  Enable ipv6 no next header ids option.
- **routing_header** (Optional, Boolean)  
  Enable ipv6 routing header ids option.
- **shim6_header** (Optional, Boolean)  
  Enable ipv6 shim header ids option.
- **user_defined_header_type** (Optional, List of String)  
  User-defined header type range.  
  Need to be `(0..255)` or `(0..255) to (0..255)`.

---

### tunnel arguments for ip

- **bad_inner_header** (Optional, Boolean)  
  Enable IP tunnel bad inner header ids option.
- **gre** (Optional, Block)  
  Declare `ip tunnel gre` configuration.
  - **gre_4in4** (Optional, Boolean)  
    Enable IP tunnel GRE 4in4 ids option.
  - **gre_4in6** (Optional, Boolean)  
    Enable IP tunnel GRE 4in6 ids option.
  - **gre_6in4** (Optional, Boolean)  
    Enable IP tunnel GRE 6in4 ids option.
  - **gre_6in6** (Optional, Boolean)  
    Enable IP tunnel GRE 6in6 ids option.
- **ip_in_udp_teredo** (Optional, Boolean)  
  Enable IP tunnel IPinUDP Teredo ids option.
- **ipip** (Optional, Block)  
  Declare `ip tunnel ipip` configuration.
  - **ipip_4in4** (Optional, Boolean)  
    Enable IP tunnel IPIP 4in4 ids option.
  - **ipip_4in6** (Optional, Boolean)  
    Enable IP tunnel IPIP 4in6 ids option.
  - **ipip_6in4** (Optional, Boolean)  
    Enable IP tunnel IPIP 6in4 ids option.
  - **ipip_6in6** (Optional, Boolean)  
    Enable IP tunnel IPIP 6in6 ids option.
  - **ipip_6over4** (Optional, Boolean)  
    Enable IP tunnel IPIP 6over4 ids option.
  - **ipip_6to4relay** (Optional, Boolean)  
    Enable IP tunnel IPIP 6to4 Relay ids option.
  - **dslite** (Optional, Boolean)  
    Enable IP tunnel IPIP DS-Lite ids option.
  - **isatap** (Optional, Boolean)  
    Enable IP tunnel IPIP ISATAP ids option.

---

### syn_flood arguments for tcp

- **alarm_threshold** (Optional, Number)  
  Alarm threshold (1..500000 requests per second).
- **attack_threshold** (Optional, Number)  
  Attack threshold (1..500000 proxied requests per second).
- **destination_threshold** (Optional, Number)  
  Destination threshold (4..500000 SYN pps).
- **source_threshold** (Optional, Number)  
  Source threshold (4..500000 SYN pps).
- **timeout** (Optional, Number)  
  SYN flood ager timeout (1..50 seconds).
- **whitelist** (Optional, Block Set)  
  For each name of white-list to declare.
  - **name** (Required, String)  
    White-list name.
  - **destination_address** (Optional, Set of String)  
    Destination address.  
    Need to be a valid CIDR network.
  - **source_address** (Optional, Set of String)  
    Source address.  
    Need to be a valid CIDR network.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security screen can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_screen.demo_screen demo_screen
```
