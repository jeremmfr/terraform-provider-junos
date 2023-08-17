package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosSecurityNatStaticRule_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityNatStaticRuleConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"name", "testacc_securityNATSttRule"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"destination_address", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.type", "prefix"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.routing_instance", "testacc_securityNATSttRule"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.prefix", "192.0.2.128/25"),
					),
				},
				{
					Config: testAccJunosSecurityNatStaticRuleConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"destination_address", "192.0.2.0/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.prefix", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule2",
							"destination_address_name", "testacc_securityNATSttRule2"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule2",
							"source_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule2",
							"source_port.#", "2"),
					),
				},
				{
					ResourceName:      "junos_security_nat_static_rule.testacc_securityNATSttRule",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_nat_static_rule.testacc_securityNATSttRule2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_nat_static_rule.testacc_securityNATSttRule3",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityNatStaticRuleConfigCreate2(),
				},
			},
		})
	}
}

func testAccJunosSecurityNatStaticRuleConfigCreate() string {
	return `
resource "junos_security_nat_static" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRule.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule" {
  name                = "testacc_securityNATSttRule"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.2.0/25"
  then {
    type             = "prefix"
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    prefix           = "192.0.2.128/25"
  }
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRuleInet" {
  name                = "testacc_securityNATSttRuleInet"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "64:ff9b::/96"
  then {
    type = "inet"
  }
}

resource "junos_security_zone" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
resource "junos_routing_instance" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
`
}

func testAccJunosSecurityNatStaticRuleConfigUpdate() string {
	return `
resource "junos_security_nat_static" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRule.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule" {
  name                = "testacc_securityNATSttRule"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.2.0/26"
  then {
    type             = "prefix"
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    prefix           = "192.0.2.64/26"
  }
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule2" {
  depends_on = [
    junos_security_address_book.testacc_securityNATSttRule
  ]
  name                     = "testacc_securityNATSttRule2"
  rule_set                 = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address_name = "testacc_securityNATSttRule2"
  source_address = [
    "192.0.2.128/26"
  ]
  source_port = [
    "1024",
    "1025 to 1026",
  ]
  then {
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    type             = "prefix-name"
    prefix           = "testacc_securityNATSttRule-prefix"
  }
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRule3" {
  depends_on = [
    junos_security_address_book.testacc_securityNATSttRule
  ]
  name                = "testacc_securityNATSttRule3"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRule.name
  destination_address = "192.0.3.1/32"
  source_address_name = [
    "testacc_securityNATSttRule-src"
  ]
  destination_port    = 81
  destination_port_to = 82
  then {
    routing_instance = junos_routing_instance.testacc_securityNATSttRule.name
    type             = "prefix"
    prefix           = "192.0.3.2/32"
    mapped_port      = 8081
    mapped_port_to   = 8082
  }
}

resource "junos_security_zone" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
resource "junos_routing_instance" "testacc_securityNATSttRule" {
  name = "testacc_securityNATSttRule"
}
resource "junos_security_address_book" "testacc_securityNATSttRule" {
  network_address {
    name  = "testacc_securityNATSttRule2"
    value = "192.0.2.128/27"
  }
  network_address {
    name  = "testacc_securityNATSttRule-prefix"
    value = "192.0.2.160/27"
  }
  network_address {
    name  = "testacc_securityNATSttRule-src"
    value = "192.0.2.224/27"
  }
}
`
}

func testAccJunosSecurityNatStaticRuleConfigCreate2() string {
	return `
resource "junos_security_nat_static" "testacc_securityNATSttRuleInet" {
  name = "testacc_securityNATSttRuleInet"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityNATSttRuleInet.name]
  }
  configure_rules_singly = true
}
resource "junos_security_nat_static_rule" "testacc_securityNATSttRuleInet" {
  name                = "testacc_securityNATSttRuleInet"
  rule_set            = junos_security_nat_static.testacc_securityNATSttRuleInet.name
  destination_address = "64:ff9b::/96"
  then {
    type             = "inet"
    routing_instance = junos_routing_instance.testacc_securityNATSttRuleInet.name
  }
}
resource "junos_security_zone" "testacc_securityNATSttRuleInet" {
  name = "testacc_securityNATSttRuleInet"
}
resource "junos_routing_instance" "testacc_securityNATSttRuleInet" {
  name = "testacc_securityNATSttRuleInet"
}
`
}
