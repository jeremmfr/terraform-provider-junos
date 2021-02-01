---
layout: "junos"
page_title: "Junos: junos_security_screen"
sidebar_current: "docs-junos-resource-security-screen"
description: |-
  Create a security screen (when Junos device supports it)
---

# junos_security_screen

Provides a security screen resource.

## Example Usage

```hcl
# Add a security screen
resource junos_security_screen "demo_screen" {
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

* `name` - (Required, Forces new resource)(`String`) The name of screen.
* `alarm_without_drop` - (Optional)(`Bool`) Do not drop packet, only generate alarm.
* `description` - (Optional)(`String`) Text description of screen.
* `icmp` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'icmp' configuration. See the [`icmp` arguments] (#icmp-arguments) block.
* `ip` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ip' configuration. See the [`ip` arguments] (#ip-arguments) block.
* `limit_session` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'limit-session' configuration.
  * `destination_ip_based` - (Optional)(`ListOfString`) Limit sessions to the same destination IP (1..2000000).
  * `source_ip_based` - (Optional)(`ListOfString`) Limit sessions from the same source IP (1..2000000).
* `tcp` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'tcp' configuration. See the [`tcp` arguments] (#tcp-arguments) block.
* `udp` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'udp' configuration. See the [`udp` arguments] (#udp-arguments) block.

---
#### icmp arguments
* `flood` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable icmp flood ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1..1000000 ICMP packets per second).
* `fragment` - (Optional)(`Bool`) Enable ICMP fragment ids option.
* `icmpv6_malformed` - (Optional)(`Bool`) Enable icmpv6 malformed ids option
* `large` - (Optional)(`Bool`) Enable large ICMP packet (size > 1024) ids option
* `ping_death` - (Optional)(`Bool`) Enable ping of death ids option
* `sweep` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable ip sweep ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1000..1000000 microseconds in which 10 ICMP packets are detected).

---
#### ip arguments
* `bad_option` - (Optional)(`Bool`) Enable ip with bad option ids option.
* `block_frag` - (Optional)(`Bool`) Enable ip fragment blocking ids option.
* `ipv6_extension_header` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ip ipv6-extension-header' configuration. See the [`ipv6_extension_header` arguments for ip] (#ipv6_extension_header-arguments-for-ip) block.
* `ipv6_extension_header_limit` - (Optional)(`Int`) Enable ipv6 extension header limit ids option (0..32).
* `ipv6_malformed_header` - (Optional)(`Bool`) Enable ipv6 malformed header ids option.
* `loose_source_route_option` - (Optional)(`Bool`) Enable ip with loose source route ids option.
* `record_route_option` - (Optional)(`Bool`) Enable ip with record route option ids option.
* `security_option` - (Optional)(`Bool`) Enable ip with security option ids option.
* `source_route_option` - (Optional)(`Bool`) Enable ip source route ids option.
* `spoofing` - (Optional)(`Bool`) Enable ip address spoofing ids option.
* `stream_option` - (Optional)(`Bool`) Enable ip with stream option ids option.
* `strict_source_route_option` - (Optional)(`Bool`) Enable ip with strict source route ids option.
* `tear_drop` - (Optional)(`Bool`) Enable tear drop ids option.
* `timestamp_option` - (Optional)(`Bool`) Enable ip with timestamp option ids option.
* `tunnel` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ip tunnel' configuration. See the [`tunnel` arguments for ip] (#tunnel-arguments-for-ip) block.
* `unknown_protocol` - (Optional)(`Bool`) Enable ip unknown protocol ids option.

---
#### tcp arguments
* `fin_no_ack` - (Optional)(`Bool`) Enable Fin bit with no ACK bit ids option.
* `land` - (Optional)(`Bool`) Enable land attack ids option.
* `no_flag` - (Optional)(`Bool`) Enable TCP packet without flag ids option.
* `port_scan` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable TCP port scan ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1000..1000000 microseconds in which 10 attack packets are detected).
* `syn_ack_ack_proxy` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable syn-ack-ack proxy ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1..250000 un-authenticated connections).
* `syn_fin` - (Optional)(`Bool`) Enable SYN and FIN bits set attack ids option.
* `syn_flood` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable SYN flood ids option. See the optional [`syn_flood` arguments for ip] (#syn_flood-arguments-for-tcp) block.
* `syn_frag` - (Optional)(`Bool`) Enable SYN fragment ids option.
* `sweep` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable TCP sweep ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1000..1000000 microseconds in which 10 TCP packets are detected).
* `winnuke` - (Optional)(`Bool`) Enable winnuke attack ids option.

---
#### udp arguments
* `flood` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to UDP flood ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1..1000000 UDP packets per second).
  * `whitelist` - (Optional)(`ListOfString`) List of UDP flood white list group name.
* `port_scan` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to UDP port scan ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1000..1000000 microseconds in which 10 attack packets are detected).
* `sweep` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to UDP sweep ids option.
  * `threshold` - (Optional)(`Int`) Threshold (1000..1000000 microseconds in which 10 UDP packets are detected).

---
#### ipv6_extension_header arguments for ip
* `ah_header` - (Optional)(`Bool`) Enable ipv6 Authentication Header ids option.
* `esp_header` - (Optional)(`Bool`) Enable ipv6 Encapsulating Security Payload header ids option.
* `hip_header` - (Optional)(`Bool`) Enable ipv6 Host Identify Protocol header ids option.
* `destination_header` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable ipv6 destination option header ids option.
  * `ilnp_nonce_option` - (Optional)(`Bool`) Enable Identifier-Locator Network Protocol Nonce option ids option.
  * `home_address_option` - (Optional)(`Bool`) Enable home address option ids option.
  * `line_identification_option` - (Optional)(`Bool`) Enable line identification option ids option.
  * `tunnel_encapsulation_limit_option` - (Optional)(`Bool`) Enable tunnel encapsulation limit option ids option.
  * `user_defined_option_type` - (Optional)(`ListOfString`) User-defined option type range. Need to be '(1..255)' or '(1..255) to (1..255)'. 
* `fragment_header` - (Optional)(`Bool`) Enable ipv6 fragment header ids option.
* `hop_by_hop_header` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable ipv6 hop by hop option header ids option.
  * `calipso_option` - (Optional)(`Bool`) Enable Common Architecture Label ipv6 Security Option ids option.
  * `rpl_option` - (Optional)(`Bool`) Enable Routing Protocol for Low-power and Lossy networks option ids option.
  * `smf_dpd_option` - (Optional)(`Bool`) Enable Simplified Multicast Forwarding ipv6 Duplicate Packet Detection option ids option.
  * `jumbo_payload_option` - (Optional)(`Bool`) Enable jumbo payload option ids option.
  * `quick_start_option` - (Optional)(`Bool`) Enable quick start option ids option.
  * `router_alert_option` - (Optional)(`Bool`) Enable router alert option ids option.
  * `user_defined_option_type` - (Optional)(`ListOfString`) User-defined option type range. Need to be '(1..255)' or '(1..255) to (1..255)'. 
* `mobility_header` - (Optional)(`Bool`) Enable ipv6 mobility header ids option.
* `no_next_header` - (Optional)(`Bool`) Enable ipv6 no next header ids option.
* `routing_header` - (Optional)(`Bool`) Enable ipv6 routing header ids option.
* `shim6_header` - (Optional)(`Bool`) Enable ipv6 shim header ids option.
* `user_defined_header_type` - (Optional)(`ListOfString`)  User-defined header type range. Need to be '(0..255)' or '(0..255) to (0..255)'.

---
#### tunnel arguments for ip
* `bad_inner_header` - (Optional)(`Bool`) Enable IP tunnel bad inner header ids option.
* `gre` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ip tunnel gre' configuration.
  * `gre_4in4` - (Optional)(`Bool`) Enable IP tunnel GRE 4in4 ids option.
  * `gre_4in6` - (Optional)(`Bool`) Enable IP tunnel GRE 4in6 ids option.
  * `gre_6in4` - (Optional)(`Bool`) Enable IP tunnel GRE 6in4 ids option.
  * `gre_6in6` - (Optional)(`Bool`) Enable IP tunnel GRE 6in6 ids option.
* `ip_in_udp_teredo` - (Optional)(`Bool`) Enable IP tunnel IPinUDP Teredo ids option.
* `ipip` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ip tunnel ipip' configuration.
  * `ipip_4in4` - (Optional)(`Bool`) Enable IP tunnel IPIP 4in4 ids option.
  * `ipip_4in6` - (Optional)(`Bool`) Enable IP tunnel IPIP 4in6 ids option.
  * `ipip_6in4` - (Optional)(`Bool`) Enable IP tunnel IPIP 6in4 ids option.
  * `ipip_6in6` - (Optional)(`Bool`) Enable IP tunnel IPIP 6in6 ids option.
  * `ipip_6over4` - (Optional)(`Bool`) Enable IP tunnel IPIP 6over4 ids option.
  * `ipip_6to4relay` - (Optional)(`Bool`) Enable IP tunnel IPIP 6to4 Relay ids option.
  * `dslite` - (Optional)(`Bool`) Enable IP tunnel IPIP DS-Lite ids option.
  * `isatap` - (Optional)(`Bool`) Enable IP tunnel IPIP ISATAP ids option.

---
#### syn_flood arguments for tcp
* `alarm_threshold` - (Optional)(`Int`) Alarm threshold (1..500000 requests per second).
* `attack_threshold` - (Optional)(`Int`) Attack threshold (1..500000 proxied requests per second).
* `destination_threshold` - (Optional)(`Int`) Destination threshold (4..500000 SYN pps).
* `source_threshold` - (Optional)(`Int`) Source threshold (4..500000 SYN pps).
* `timeout` - (Optional)(`Int`) SYN flood ager timeout (1..50 seconds).
* `whitelist` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each white-list to declare.
  * `name` - (Required)(`String`) White-list name.
  * `destination_address` - (Optional)(`ListOfString`) Destination address. Need to be a valid CIDR network.
  * `source_address` - (Optional)(`ListOfString`) Source address. Need to be a valid CIDR network.

## Import

Junos security screen can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_screen.demo_screen demo_screen
```
