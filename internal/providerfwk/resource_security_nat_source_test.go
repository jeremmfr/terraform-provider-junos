package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosSecurityNatSource_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityNatSourceConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"from.type", "zone"),
						resource.TestCheckTypeSetElemAttr("junos_security_nat_source.testacc_securitySNAT",
							"from.value.*", "testacc_securitySNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"to.type", "zone"),
						resource.TestCheckTypeSetElemAttr("junos_security_nat_source.testacc_securitySNAT",
							"to.value.*", "testacc_securitySNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.name", "testacc_securitySNATRule"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.source_address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.source_address.*", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.destination_address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.destination_address.*", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.protocol.*", "tcp"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.then.type", "pool"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.then.pool", "testacc_securitySNATPool"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"address.#", "2"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"address.0", "192.0.2.1/32"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"address.1", "192.0.2.64/27"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"routing_instance", "testacc_securitySNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"address_pooling", "paired"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"port_no_translation", "true"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"pool_utilization_alarm_raise_threshold", "80"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"pool_utilization_alarm_clear_threshold", "60"),
					),
				},
				{
					Config: testAccJunosSecurityNatSourceConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.#", "3"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.1.then.type", "off"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"address_pooling", "no-paired"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"port_overloading_factor", "3"),
					),
				},
				{
					ResourceName:      "junos_security_nat_source.testacc_securitySNAT",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityNatSourceConfigCreate() string {
	return `
resource "junos_security_nat_source" "testacc_securitySNAT" {
  name        = "testacc_securitySNAT"
  description = "testacc securitySNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT.name]
  }
  to {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT.name]
  }
  rule {
    name = "testacc_securitySNATRule"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["tcp"]
    }
    then {
      type = "pool"
      pool = junos_security_nat_source_pool.testacc_securitySNATPool.name
    }
  }
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool" {
  name                                   = "testacc_securitySNATPool"
  description                            = "testacc securitySNATPool"
  address                                = ["192.0.2.1/32", "192.0.2.64/27"]
  routing_instance                       = junos_routing_instance.testacc_securitySNAT.name
  address_pooling                        = "paired"
  port_no_translation                    = true
  pool_utilization_alarm_raise_threshold = 80
  pool_utilization_alarm_clear_threshold = 60
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool2" {
  name             = "testacc_securitySNATPool2"
  description      = "testacc securitySNATPool2"
  address          = ["192.0.2.2/32"]
  routing_instance = junos_routing_instance.testacc_securitySNAT.name
  port_range       = "1300"
}

resource "junos_security_zone" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
resource "junos_routing_instance" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
`
}

func testAccJunosSecurityNatSourceConfigUpdate() string {
	return `
resource "junos_security_nat_source" "testacc_securitySNAT" {
  depends_on = [
    junos_security_address_book.testacc_securitySNAT
  ]
  name = "testacc_securitySNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT.name]
  }
  to {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT.name]
  }
  rule {
    name = "testacc_securitySNATRule"
    match {
      source_address           = ["192.0.2.0/25"]
      source_address_name      = ["testacc_securitySNAT2"]
      source_port              = ["1024", "1021 to 1022"]
      destination_address      = ["192.0.2.128/25"]
      destination_address_name = ["testacc_securitySNAT"]
      destination_port         = ["80", "82 to 83"]
      protocol                 = ["tcp"]
    }
    then {
      type = "pool"
      pool = junos_security_nat_source_pool.testacc_securitySNATPool.name
    }
  }
  rule {
    name = "testacc_securitySNATRule2"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["udp"]
    }
    then {
      type = "off"
    }
  }
  rule {
    name = "testacc_securitySNATRule3"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      application         = ["junos-ssh", "junos-http"]
    }
    then {
      type = "off"
    }
  }
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool" {
  name                    = "testacc_securitySNATPool"
  address                 = ["192.0.2.1/32"]
  routing_instance        = junos_routing_instance.testacc_securitySNAT.name
  address_pooling         = "no-paired"
  port_overloading_factor = 3
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool2" {
  name             = "testacc_securitySNATPool2"
  description      = "testacc securitySNATPool2"
  address          = ["192.0.2.2/32"]
  routing_instance = junos_routing_instance.testacc_securitySNAT.name
  port_range       = "1300-62000"
}

resource "junos_security_zone" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
resource "junos_routing_instance" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
resource "junos_security_address_book" "testacc_securitySNAT" {
  network_address {
    name  = "testacc_securitySNAT"
    value = "192.0.2.128/27"
  }
  network_address {
    name  = "testacc_securitySNAT2"
    value = "192.0.2.160/27"
  }
}
`
}
