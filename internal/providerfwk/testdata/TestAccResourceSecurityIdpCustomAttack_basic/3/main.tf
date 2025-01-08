resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name     = "testacc/#1_"
  severity = "minor"
  attack_type_signature {
    context          = "http-url"
    direction        = "any"
    negate           = true
    pattern          = "test"
    pattern_pcre     = "test"
    protocol_binding = "ip protocol-number 58"
    regexp           = "test"
    shellcode        = "all"
    protocol_ipv4 {
      checksum_validate_match = "equal"
      checksum_validate_value = 1
      destination_match       = "not-equal"
      destination_value       = "192.0.2.3"
      identification_match    = "equal"
      identification_value    = 2
      ihl_match               = "less-than"
      ihl_value               = 3
      ip_flags                = ["df", "mf"]
      protocol_match          = "equal"
      protocol_value          = 4
      source_match            = "greater-than"
      source_value            = "192.0.2.4"
      tos_match               = "equal"
      tos_value               = 6
      total_length_match      = "equal"
      total_length_value      = 7
      ttl_match               = "equal"
      ttl_value               = 8
    }
    protocol_tcp {
      ack_number_match        = "equal"
      ack_number_value        = 10
      checksum_validate_match = "greater-than"
      checksum_validate_value = 11
      data_length_match       = "equal"
      data_length_value       = 12
      destination_port_match  = "not-equal"
      destination_port_value  = 13
      header_length_match     = "equal"
      header_length_value     = 14
      mss_match               = "equal"
      mss_value               = 15
      option_match            = "equal"
      option_value            = 16
      reserved_match          = "less-than"
      reserved_value          = 5
      sequence_number_match   = "equal"
      sequence_number_value   = 18
      source_port_match       = "equal"
      source_port_value       = 19
      tcp_flags               = ["fin", "no-ack"]
      urgent_pointer_match    = "equal"
      urgent_pointer_value    = 20
      window_scale_match      = "equal"
      window_scale_value      = 21
      window_size_match       = "equal"
      window_size_value       = 22
    }
  }
}
