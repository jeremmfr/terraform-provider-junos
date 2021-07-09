package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityNatSource_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" && os.Getenv("TESTACC_ROUTER") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityNatSourceConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"from.0.type", "zone"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"from.0.value.0", "testacc_securitySNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"to.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"to.0.type", "zone"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"to.0.value.0", "testacc_securitySNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.name", "testacc_securitySNATRule"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.0.source_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.0.source_address.0", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.0.destination_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.0.destination_address.0", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.0.protocol.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.match.0.protocol.0", "tcp"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.then.0.type", "pool"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.0.then.0.pool", "testacc_securitySNATPool"),
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
							"rule.#", "2"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.1.match.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.1.then.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_source.testacc_securitySNAT",
							"rule.1.then.0.type", "off"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"address_pooling", "no-paired"),
						resource.TestCheckResourceAttr("junos_security_nat_source_pool.testacc_securitySNATPool",
							"port_no_translation", "false"),
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
resource junos_security_nat_source testacc_securitySNAT {
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
resource junos_security_nat_source_pool testacc_securitySNATPool {
  name                                   = "testacc_securitySNATPool"
  address                                = ["192.0.2.1/32", "192.0.2.64/27"]
  routing_instance                       = junos_routing_instance.testacc_securitySNAT.name
  address_pooling                        = "paired"
  port_no_translation                    = true
  pool_utilization_alarm_raise_threshold = 80
  pool_utilization_alarm_clear_threshold = 60
}

resource junos_security_zone testacc_securitySNAT {
  name = "testacc_securitySNAT"
}
resource junos_routing_instance testacc_securitySNAT {
  name = "testacc_securitySNAT"
}
`
}

func testAccJunosSecurityNatSourceConfigUpdate() string {
	return `
resource junos_security_nat_source testacc_securitySNAT {
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
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["tcp"]
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
}
resource junos_security_nat_source_pool testacc_securitySNATPool {
  name                    = "testacc_securitySNATPool"
  address                 = ["192.0.2.1/32"]
  routing_instance        = junos_routing_instance.testacc_securitySNAT.name
  address_pooling         = "no-paired"
  port_overloading_factor = 3
}

resource junos_security_zone testacc_securitySNAT {
  name = "testacc_securitySNAT"
}
resource junos_routing_instance testacc_securitySNAT {
  name = "testacc_securitySNAT"
}
`
}
