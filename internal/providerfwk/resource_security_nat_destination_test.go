package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityNatDestination_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"from.type", "zone"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"from.value.0", "testacc_securityDNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.name", "testacc_securityDNATRule"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.destination_address", "192.0.2.1/32"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.then.type", "pool"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.0.then.pool", "testacc_securityDNATPool"),
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
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.#", "3"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.1.destination_address_name", "testacc_securityDNAT"),
						resource.TestCheckResourceAttr("junos_security_nat_destination.testacc_securityDNAT",
							"rule.1.then.type", "pool"),
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
