package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityIdpCustomAttack_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityIdpCustomAttackConfigCreate(),
				},
				{
					Config: testAccResourceSecurityIdpCustomAttackConfigUpdate(),
				},
				{
					Config: testAccResourceSecurityIdpCustomAttackConfigUpdate2(),
				},
				{
					ResourceName:      "junos_security_idp_custom_attack.testacc_idpCustomAttack",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSecurityIdpCustomAttackConfigUpdate3(),
				},
				{
					Config: testAccResourceSecurityIdpCustomAttackConfigUpdate4(),
				},
				{
					Config: testAccResourceSecurityIdpCustomAttackConfigUpdate5(),
				},
			},
		})
	}
}

func testAccResourceSecurityIdpCustomAttackConfigCreate() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "ignore"
  severity           = "info"
  attack_type_anomaly {
    direction = "any"
    service   = "TELNET"
    test      = "SUBOPTION_OVERFLOW"
    shellcode = "all"
  }
}
`
}

func testAccResourceSecurityIdpCustomAttackConfigUpdate() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "none"
  severity           = "minor"
  attack_type_chain {
    protocol_binding = "application HTTP"
    member {
      name = "testacc/#1_chain_member1"
      attack_type_anomaly {
        direction = "any"
        test      = "MISSING_HOST"
        shellcode = "all"
      }
    }
    member {
      name = "testacc/#1_chain_member2"
      attack_type_anomaly {
        direction = "any"
        test      = "ACCEPT_LANGUAGE_OVERFLOW"
        shellcode = "all"
      }
    }
  }
}
`
}

func testAccResourceSecurityIdpCustomAttackConfigUpdate2() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "none"
  severity           = "minor"
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
`
}

func testAccResourceSecurityIdpCustomAttackConfigUpdate3() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "none"
  severity           = "minor"
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
`
}

func testAccResourceSecurityIdpCustomAttackConfigUpdate4() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "none"
  severity           = "minor"
  attack_type_signature {
    context   = "http-url"
    direction = "client-to-server"
    protocol_icmp {
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
`
}

func testAccResourceSecurityIdpCustomAttackConfigUpdate5() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttack" {
  name               = "testacc/#1_"
  recommended_action = "none"
  severity           = "minor"
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
`
}
