package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityNatStatic_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.type", "zone"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.value.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"from.value.0", "testacc_securityNATStt"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.#", "2"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.name", "testacc_securityNATSttRule"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.destination_address", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.type", "prefix"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.routing_instance", "testacc_securityNATStt"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.prefix", "192.0.2.128/25"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.#", "3"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.destination_address", "192.0.2.0/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.0.then.prefix", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.1.destination_address_name", "testacc_securityNATSttRule2"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.1.source_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_nat_static.testacc_securityNATStt",
							"rule.1.source_port.#", "2"),
					),
				},
				{
					ResourceName:      "junos_security_nat_static.testacc_securityNATStt",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:  "junos_security_nat_static.testacc_securityNATStt_singly",
					ImportState:   true,
					ImportStateId: "testacc_securityNATStt_singly" + junos.IDSeparator + "no_rules",
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
