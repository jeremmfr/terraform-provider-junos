package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityNatDestination_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" && os.Getenv("TESTACC_ROUTER") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityNatDestinationConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"from.0.type", "zone"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"from.0.value.0", "testacc_securityDNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.name", "testacc_securityDNATRule"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.destination_address", "192.0.2.1/32"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.then.0.type", "pool"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.then.0.pool", "testacc_securityDNATPool"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.then.0.pool", "testacc_securityDNATPool"),
						resource.TestCheckResourceAttr("junos_security_nat_destination_pool.testacc_securityDNATPool",
							"address", "192.0.2.1/32"),
						resource.TestCheckResourceAttr("junos_security_nat_destination_pool.testacc_securityDNATPool",
							"address_to", "192.0.2.2/32"),
						resource.TestCheckResourceAttr("junos_security_nat_destination_pool.testacc_securityDNATPool",
							"routing_instance", "testacc_securityDNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_destination_pool.testacc_securityDNATPool2",
							"address_port", "80"),
					),
				},
				{
					Config: testAccJunosSecurityNatDestinationConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.#", "3"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.1.destination_address_name", "testacc_securityDNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.1.then.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.1.then.0.type", "pool"),
					),
				},
				{
					ResourceName:      "junos_security_nat_destination.testacc_securityDNAT",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityNatDestinationConfigCreate() string {
	return `
resource junos_security_nat_destination testacc_securityDNAT {
  name = "testacc_securityDNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityDNAT.name]
  }
  rule {
    name                = "testacc_securityDNATRule"
    destination_address = "192.0.2.1/32"
    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool.name
    }
  }
}
resource junos_security_nat_destination_pool testacc_securityDNATPool {
  name             = "testacc_securityDNATPool"
  address          = "192.0.2.1/32"
  address_to       = "192.0.2.2/32"
  routing_instance = junos_routing_instance.testacc_securityDNAT.name
}
resource junos_security_nat_destination_pool testacc_securityDNATPool2 {
  name             = "testacc_securityDNATPool2"
  address          = "192.0.2.1/32"
  address_port     = 80
  routing_instance = junos_routing_instance.testacc_securityDNAT.name
}

resource junos_security_zone testacc_securityDNAT {
  name = "testacc_securityDNAT"
}
resource junos_routing_instance testacc_securityDNAT {
  name = "testacc_securityDNAT"
}
`
}

func testAccJunosSecurityNatDestinationConfigUpdate() string {
	return `
resource junos_security_nat_destination testacc_securityDNAT {
  depends_on = [
    junos_security_address_book.testacc_securityDNAT
  ]
  name = "testacc_securityDNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityDNAT.name]
  }
  rule {
    name                = "testacc_securityDNATRule"
    destination_address = "192.0.2.1/32"
    source_address_name = [
      "testacc_securityDNAT-src",
    ]
    protocol = ["tcp", "50"]
    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool.name
    }
  }
  rule {
    name                     = "testacc_securityDNATRule2"
    destination_address_name = "testacc_securityDNAT"
    destination_port = [
      "81",
      "82 to 83",
    ]
    source_address = [
      "192.0.2.128/26",
    ]
    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool2.name
    }
  }
  rule {
    name                = "testacc_securityDNATRule3"
    destination_address = "192.0.2.1/32"
    application = [
      "junos-ssh", "junos-http",
    ]

    then {
      type = "pool"
      pool = junos_security_nat_destination_pool.testacc_securityDNATPool2.name
    }
  }
}
resource junos_security_nat_destination_pool testacc_securityDNATPool {
  name             = "testacc_securityDNATPool"
  address          = "192.0.2.1/32"
  address_to       = "192.0.2.2/32"
  routing_instance = junos_routing_instance.testacc_securityDNAT.name
}
resource junos_security_nat_destination_pool testacc_securityDNATPool2 {
  name             = "testacc_securityDNATPool2"
  address          = "192.0.2.1/32"
  address_port     = 80
  routing_instance = junos_routing_instance.testacc_securityDNAT.name
}

resource junos_security_zone testacc_securityDNAT {
  name = "testacc_securityDNAT"
}
resource junos_routing_instance testacc_securityDNAT {
  name = "testacc_securityDNAT"
}
resource "junos_security_address_book" "testacc_securityDNAT" {
  network_address {
    name  = "testacc_securityDNAT"
    value = "192.0.2.128/27"
  }
  network_address {
    name  = "testacc_securityDNAT-src"
    value = "192.0.2.160/27"
  }
}
`
}
