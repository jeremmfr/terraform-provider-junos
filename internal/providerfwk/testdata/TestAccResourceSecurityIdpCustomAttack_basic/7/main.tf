resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name     = "testacc/#1_"
  severity = "minor"
  attack_type_signature {
    context   = "http-url"
    direction = "any"
    protocol_icmpv6 {
      checksum_validate_match = "equal"
      checksum_validate_value = 0
      code_match              = "not-equal"
      code_value              = 1
      data_length_match       = "greater-than"
      data_length_value       = 2
      identification_match    = "less-than"
      identification_value    = 3
      sequence_number_match   = "equal"
      sequence_number_value   = 4
      type_match              = "not-equal"
      type_value              = 5
    }
  }
}
