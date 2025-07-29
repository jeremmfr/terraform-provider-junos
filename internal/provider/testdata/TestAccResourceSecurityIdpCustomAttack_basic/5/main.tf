resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name     = "testacc/#1_"
  severity = "minor"
  attack_type_signature {
    context   = "http-url"
    direction = "any"
    protocol_ipv6 {
      destination_match                                      = "not-equal"
      destination_value                                      = "2001:db8:85a4::1"
      extension_header_destination_option_home_address_match = "equal"
      extension_header_destination_option_home_address_value = "2001:db8:85a4::3"
      extension_header_destination_option_type_match         = "equal"
      extension_header_destination_option_type_value         = 0
      extension_header_routing_header_type_match             = "equal"
      extension_header_routing_header_type_value             = 2
      flow_label_match                                       = "equal"
      flow_label_value                                       = 3
      hop_limit_match                                        = "greater-than"
      hop_limit_value                                        = 4
      next_header_match                                      = "less-than"
      next_header_value                                      = 5
      payload_length_match                                   = "equal"
      payload_length_value                                   = 6
      source_match                                           = "equal"
      source_value                                           = "2001:db8:85a4::2"
      traffic_class_match                                    = "equal"
      traffic_class_value                                    = 7
    }
    protocol_udp {
      checksum_validate_match = "equal"
      checksum_validate_value = 8
      data_length_match       = "equal"
      data_length_value       = 9
      destination_port_match  = "greater-than"
      destination_port_value  = 10
      source_port_match       = "not-equal"
      source_port_value       = 11
    }
  }
}
