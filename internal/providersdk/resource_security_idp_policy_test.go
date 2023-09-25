package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityIdpPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityIdpPolicyConfigCreate(),
				},
				{
					ResourceName:      "junos_security_idp_policy.testacc_idp_pol",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSecurityIdpPolicyConfigUpdate(),
				},
				{
					Config: testAccResourceSecurityIdpPolicyConfigPostCheck(),
				},
			},
		})
	}
}

func testAccResourceSecurityIdpPolicyConfigCreate() string {
	return `
resource "junos_security" "testacc_secIdpPolicy" {
  lifecycle {
    ignore_changes = [
      log
    ]
  }
  idp_sensor_configuration {
    packet_log {
      source_address = "192.0.2.4"
      host_address   = "192.0.2.5"
      host_port      = 514
    }
  }
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
}
resource "junos_security_zone" "testacc_idp_pol_from" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_idp_pol_from"
}
resource "junos_security_zone" "testacc_idp_pol_to" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_idp_pol_to"
}
resource "junos_security_idp_policy" "testacc_idp_pol" {
  name = "testacc_idp/#1"
  ips_rule {
    name        = "rules_#B"
    description = "rules _ test #B"
    terminal    = true
    match {
      application             = "junos:telnet"
      custom_attack_group     = ["test1"]
      custom_attack           = ["test2"]
      destination_address     = ["192.0.2.0/24"]
      dynamic_attack_group    = ["test3"]
      from_zone               = junos_security_zone.testacc_idp_pol_from.name
      predefined_attack_group = ["test4"]
      predefined_attack       = ["test5"]
      source_address          = ["192.0.2.0/24"]
      to_zone                 = junos_security_zone.testacc_idp_pol_to.name
    }
    then {
      action                                      = "drop-connection"
      ip_action                                   = "ip-close"
      ip_action_log                               = true
      ip_action_log_create                        = true
      ip_action_refresh_timeout                   = true
      ip_action_target                            = "service"
      ip_action_timeout                           = 60
      notification                                = true
      notification_log_attacks                    = true
      notification_log_attacks_alert              = true
      notification_packet_log                     = true
      notification_packet_log_post_attack         = 10
      notification_packet_log_post_attack_timeout = 90
      notification_packet_log_pre_attack          = 20
      severity                                    = "info"
    }
  }
  ips_rule {
    name        = "rules_#A"
    description = "rules _ test #A"
    match {
      destination_address_except = ["192.0.2.1/32"]
      from_zone                  = junos_security_zone.testacc_idp_pol_from.name
      source_address_except      = ["192.0.2.254/32"]
      to_zone                    = junos_security_zone.testacc_idp_pol_to.name
    }
    then {
      action = "drop-packet"
    }
  }
}
`
}

func testAccResourceSecurityIdpPolicyConfigUpdate() string {
	return `
resource "junos_security" "testacc_secIdpPolicy" {
  idp_sensor_configuration {
    packet_log {
      source_address = "192.0.2.4"
      host_address   = "192.0.2.5"
      host_port      = 514
    }
  }
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
}
resource "junos_security_idp_policy" "testacc_idp_pol" {
  name = "testacc_idp/#1"
  ips_rule {
    name        = "rules_#B"
    description = "rules _ test #B"
    match {
      application = "junos:telnet"
    }
    then {
      action = "drop-connection"
    }
  }
  exempt_rule {
    name        = "rules_#A"
    description = "rules _ test #A"
    match {
      destination_address_except = ["192.0.2.1/32"]
      source_address_except      = ["192.0.2.254/32"]
    }
  }
}
`
}

func testAccResourceSecurityIdpPolicyConfigPostCheck() string {
	return `
resource "junos_security" "testacc_secIdpPolicy" {
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
}
`
}
