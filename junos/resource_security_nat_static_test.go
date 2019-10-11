package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccJunosSecurityNatStatic_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityNatStaticConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.0.type", "zone"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.0.value.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.0.value.0", "testacc_securityNATStt"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.name", "testacc_securityNATSttRule"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.destination_address", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.0.type", "prefix"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.0.routing_instance", "testacc_securityNATStt"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.0.prefix", "192.0.2.128/25"),
					),
				},
				{
					Config: testAccJunosSecurityNatStaticConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.#", "2"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.destination_address", "192.0.2.0/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.0.prefix", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.1.destination_address", "192.0.2.128/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.1.then.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.1.then.0.type", "prefix"),
					),
				},
				{
					ResourceName:      "junos_security_nat_static.testacc_securityNATStt",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityNatStaticConfigCreate() string {
	return fmt.Sprintf(`
resource junos_security_nat_static testacc_securityNATStt {
  name = "testacc_securityNATStt"
  from {
    type = "zone"
    value = [ junos_security_zone.testacc_securityNATStt.name ]
  }
  rule {
    name = "testacc_securityNATSttRule"
    destination_address = "192.0.2.0/25"
    then {
      type = "prefix"
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      prefix = "192.0.2.128/25"
    }
  }
}

resource junos_security_zone testacc_securityNATStt {
  name = "testacc_securityNATStt"
}
resource junos_routing_instance testacc_securityNATStt {
  name = "testacc_securityNATStt"
}
`)
}
func testAccJunosSecurityNatStaticConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_security_nat_static testacc_securityNATStt {
  name = "testacc_securityNATStt"
  from {
    type = "zone"
    value = [ junos_security_zone.testacc_securityNATStt.name ]
  }
  rule {
    name = "testacc_securityNATSttRule"
    destination_address = "192.0.2.0/26"
    then {
      type = "prefix"
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      prefix = "192.0.2.64/26"
    }
  }
  rule {
    name = "testacc_securityNATSttRule2"
    destination_address = "192.0.2.128/26"
    then {
      routing_instance = junos_routing_instance.testacc_securityNATStt.name
      type = "prefix"
       prefix = "192.0.2.192/26"
    }
  }
}

resource junos_security_zone testacc_securityNATStt {
  name = "testacc_securityNATStt"
}
resource junos_routing_instance testacc_securityNATStt {
  name = "testacc_securityNATStt"
}
`)
}
