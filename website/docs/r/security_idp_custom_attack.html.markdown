---
layout: "junos"
page_title: "Junos: junos_security_idp_custom_attack"
sidebar_current: "docs-junos-resource-security-idp-custom-attack"
description: |-
  Create a security idp custom-attack (when Junos device supports it)
---

# junos_security_idp_custom_attack

Provides a security idp custom-attack resource.

## Example Usage

```hcl
# Add a idp custom-attack
resource junos_security_idp_custom_attack "demo_idp_custom_attack" {
  name               = "SSH:BRUTE-FORCE-CUSTOM"
  recommended_action = "drop"
  severity           = "minor"
  time_binding_count = 120
  time_binding_scope = "peer"
  attack_type_signature {
    protocol_binding = "application SSH"
    context          = "first-data-packet"
    pattern          = "\\[SSH\\].*"
    direction        = "client-to-server"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of idp custom-attack.
* `recommended_action` - (Required)(`String`) Recommended Action. Need to be 'close', 'close-client', 'close-server', 'drop', 'drop-packet', 'ignore' or 'none'.
* `severity` - (Required)(`String`) Select the severity that matches the lethality of this attack on your network. Need to be 'critical', 'info' 'major', 'minor' or 'warning'.
* `attack_type_anomaly` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'attack-type anomaly' configuration. Need to set one of three: `attack_type_anomaly`, `attack_type_chain` or `attack_type_signature`. See the [`attack_type_anomaly` arguments] (#attack_type_anomaly-arguments) block.
* `attack_type_chain` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'attack-type chain' configuration. Need to set one of three: `attack_type_anomaly`, `attack_type_chain` or `attack_type_signature`. See the [`attack_type_chain` arguments] (#attack_type_chain-arguments) block.
* `attack_type_signature` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'attack-type signature' configuration. Need to set one of three: `attack_type_anomaly`, `attack_type_chain` or `attack_type_signature`. See the [`attack_type_signature` arguments] (#attack_type_signature-arguments) block.
* `time_binding_count` - (Optional)(`Int`) Number of times this attack is to be triggered.
* `time_binding_scope` - (Optional)(`String`) Scope within which the count occurs. Need to be 'destination', 'peer' or 'source'.

---
#### attack_type_anomaly arguments
* `direction` - (Required)(`String`) Connection direction of the attack. Need to be 'any', 'client-to-server' or 'server-to-client'.
* `service` - (Required)(`String`) Service name.
* `test` - (Required)(`String`) Protocol anomaly condition to be checked.
* `shellcode` - (Optional)(`String`) Specify shellcode flag for this attack. Need to be 'all', 'intel', 'no-shellcode' or 'sparc'.

---
#### attack_type_chain arguments
* `member` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each member attack to declare.
  * `name` - (Required)(`String`) Custom attack name.
  * `attack_type_anomaly` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'attack-type anomaly' configuration. Need to set one of two: `attack_type_anomaly` or `attack_type_signature`. See the [`attack_type_anomaly` arguments] (#attack_type_anomaly-arguments) block but without `service` argument.
  * `attack_type_signature` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'attack-type signature' configuration. Need to set one of two: `attack_type_anomaly` or `attack_type_signature`. See the [`attack_type_signature` arguments] (#attack_type_signature-arguments) block but without `protocol_binding` argument.
* `expression` - (Optional)(`String`) Boolean Expression.
* `order` - (Optional)(`Bool`) Attacks should match in the order in which they are defined.
* `protocol_binding` - (Optional)(`String`) Protocol binding over which attack will be detected. Need to start with 'application', 'icmp', 'ip', 'rpc', 'tcp' or 'udp' string.
* `reset` - (Optional)(`Bool`) Repeat match should generate a new alert.
* `scope` - (Optional)(`String`) Scope of the attack. Need to be 'session' or 'transaction'. 

---
#### attack_type_signature arguments
* `context` - (Required)(`String`) Context.
* `direction` - (Required)(`String`) Connection direction of the attack. Need to be 'any', 'client-to-server' or 'server-to-client'.
* `negate` - (Optional)(`Bool`) Trigger the attack if condition is not met.
* `pattern` - (Optional)(`String`) Pattern is the signature of the attack you want to detect.
* `pattern_pcre` - (Optional)(`String`) Attack signature pattern in PCRE format.
* `protocol_icmp` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'protocol icmp' configuration.
  * `checksum_validate_match` - (Optional)(`String`) Condition for validate checksum field against calculated checksum. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `checksum_validate_value` - (Optional)(`Int`) Value for validate checksum field against calculated checksum.
  * `code_match` - (Optional)(`String`) Condition for code field. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `code_value` - (Optional)(`Int`) Value for code field.
  * `data_length_match` - (Optional)(`String`) Condition for size of IP datagram subtracted by ICMP header length. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `data_length_value` - (Optional)(`Int`) Value for size of IP datagram subtracted by ICMP header length.
  * `identification_match` - (Optional)(`String`) Condition for identifier in echo request/reply. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `identification_value` - (Optional)(`Int`) Value for identifier in echo request/reply.
  * `sequence_number_match` - (Optional)(`String`) Condition for sequence number. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `sequence_number_value` - (Optional)(`Int`) Value for sequence number.
  * `type_match` - (Optional)(`String`) Condition for type. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `type_value` - (Optional)(`Int`) Value for type.
* `protocol_icmpv6` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'protocol icmpv6' configuration. Same arguments as for 'protocol_icmp'.
* `protocol_ipv4` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'protocol ipv4' configuration.
  * `checksum_validate_match` - (Optional)(`String`) Condition for validate checksum field against calculated checksum. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `checksum_validate_value` - (Optional)(`Int`) Value for validate checksum field against calculated checksum.
  * `destination_match` - (Optional)(`String`) Condition for destination IP-address. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `destination_value` - (Optional)(`String`) Value for destination IP-address.
  * `identification_match` - (Optional)(`String`) Condition for fragment identification. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `identification_value` - (Optional)(`Int`) Value for fragment identification.
  * `ihl_match` - (Optional)(`String`) Condition for header length in words. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `ihl_value` - (Optional)(`Int`) Value for header length in words.
  * `ip_flags` - (Optional)(`ListOfString`) IP Flag bits.
  * `protocol_match` - (Optional)(`String`) Condition for transport layer protocol. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `protocol_value` - (Optional)(`Int`) Value for transport layer protocol.
  * `source_match` - (Optional)(`String`) Condition for source IP-address. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `source_value` - (Optional)(`String`) Value for source IP-address.
  * `tos_match` - (Optional)(`String`) Condition for type of service. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `tos_value` - (Optional)(`String`) Value for type of service.
  * `total_length_match` - (Optional)(`String`) Condition for total length of IP datagram. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `total_length_value` - (Optional)(`String`) Value for total length of IP datagram.
  * `ttl_match` - (Optional)(`String`) Condition for time to live. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `ttl_value` - (Optional)(`String`) Value for time to live.
* `protocol_ipv6` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'protocol ipv6' configuration.
  * `destination_match` - (Optional)(`String`) Condition for destination IP-address. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `destination_value` - (Optional)(`String`) Value for destination IP-address.
  * `extension_header_destination_option_home_address_match` - (Optional)(`String`) Condition for home address of the mobile node in destination option extension header. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `extension_header_destination_option_home_address_value` - (Optional)(`Int`) Value for home address of the mobile node in destination option extension header.
  * `extension_header_destination_option_type_match` - (Optional)(`String`) Condition for header type in destination option extension header. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `extension_header_destination_option_type_value` - (Optional)(`Int`) Value for header type in  destination option extension header.
  * `extension_header_routing_header_type_match` - (Optional)(`String`) Condition for header type in routing extension header. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `extension_header_routing_header_type_value` - (Optional)(`Int`) Value for header type in routing extension header.
  * `flow_label_match` - (Optional)(`String`) Condition for flow label identification. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `flow_label_value` - (Optional)(`String`) Value for flow label identification.
  * `hop_limit_match` - (Optional)(`String`) Condition for hop limit. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `hop_limit_value` - (Optional)(`String`) Value for hop limit.
  * `next_header_match` - (Optional)(`String`) Condition for the header following the basic IPv6 header. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `next_header_value` - (Optional)(`String`) Value for the header following the basic IPv6 header.
  * `payload_length_match` - (Optional)(`String`) Condition for length of the payload in the IPv6 datagram. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `payload_length_value` - (Optional)(`String`) Value for length of the payload in the IPv6 datagram.
  * `source_match` - (Optional)(`String`) Condition for source IP-address. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `source_value` - (Optional)(`String`) Value for source IP-address.
  * `traffic_class_match` - (Optional)(`String`) Condition for traffic class. Similar to TOS in IPv4. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `traffic_class_value` - (Optional)(`String`) Value for traffic class. Similar to TOS in IPv4.
* `protocol_tcp` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'protocol tcp' configuration.
  * `ack_number_match` - (Optional)(`String`) Condition for acknowledgement number. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `ack_number_value` - (Optional)(`String`) Value for acknowledgement number.
  * `checksum_validate_match` - (Optional)(`String`) Condition for validate checksum field against calculated checksum. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `checksum_validate_value` - (Optional)(`Int`) Value for validate checksum field against calculated checksum.
  * `data_length_match` - (Optional)(`String`) Condition for size of IP datagram subtracted by TCP header length. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `data_length_value` - (Optional)(`Int`) Value for size of IP datagram subtracted by TCP header length.
  * `destination_port_match` - (Optional)(`String`) Condition for destination port. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `destination_port_value` - (Optional)(`Int`) Value for destination port.
  * `header_length_match` - (Optional)(`String`) Condition for header length in words. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `header_length_value` - (Optional)(`Int`) Value for header length in words.
  * `mss_match` - (Optional)(`String`) Condition for maximum segment size. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `mss_value` - (Optional)(`Int`) Value for maximum segment size.
  * `option_match` - (Optional)(`String`) Condition for kind. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `option_value` - (Optional)(`Int`) Value for kind.
  * `reserved_match` - (Optional)(`String`) Condition for three reserved bits. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `reserved_value` - (Optional)(`Int`) Value for three reserved bits.
  * `sequence_number_match` - (Optional)(`String`) Condition for sequence number. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `sequence_number_value` - (Optional)(`Int`) Value for sequence number.
  * `source_port_match` - (Optional)(`String`) Condition for source port. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `source_port_value` - (Optional)(`Int`) Value for source port.
  * `tcp_flags` - (Optional)(`ListOfString`) TCP header flags.
  * `urgent_pointer_match` - (Optional)(`String`) Condition for urgent pointer. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `urgent_pointer_value` - (Optional)(`Int`) Value for urgent pointer.
  * `window_scale_match` - (Optional)(`String`) Condition for window scale. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `window_scale_value` - (Optional)(`Int`) Value for sindow scale.
  * `window_size_match` - (Optional)(`String`) Condition for window size. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `window_size_value` - (Optional)(`Int`) Value for sindow size.
* `protocol_udp` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'protocol udp' configuration.
  * `checksum_validate_match` - (Optional)(`String`) Condition for validate checksum field against calculated checksum. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `checksum_validate_value` - (Optional)(`Int`) Value for validate checksum field against calculated checksum.
  * `data_length_match` - (Optional)(`String`) Condition for size of IP datagram subtracted by UDP header length. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `data_length_value` - (Optional)(`Int`) Value for size of IP datagram subtracted by UDP header length.
  * `destination_port_match` - (Optional)(`String`) Condition for destination port. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `destination_port_value` - (Optional)(`Int`) Value for destination port.
  * `source_port_match` - (Optional)(`String`) Condition for source port. Need to be 'equal', 'greater-than', 'less-than' or 'not-equal'.
  * `source_port_value` - (Optional)(`Int`) Value for source port.
* `protocol_binding` - (Optional)(`String`) Protocol binding over which attack will be detected. Need to start with 'application', 'icmp', 'ip', 'rpc', 'tcp' or 'udp' string.
* `regexp` - (Optional)(`String`) Regular expression used for matching repetition of patterns.
* `shellcode` - (Optional)(`String`) Specify shellcode flag for this attack. Need to be 'all', 'intel', 'no-shellcode' or 'sparc'.

## Import

Junos security idp custom-attack can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_security_idp_custom_attack.demo_idp_custom_attack 'SSH:BRUTE-FORCE-CUSTOM'
```
