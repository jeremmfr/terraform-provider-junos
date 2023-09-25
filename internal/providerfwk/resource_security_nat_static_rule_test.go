package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityNatStaticRule_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"name", "testacc_securityNATSttRule"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"destination_address", "192.0.2.0/28"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.type", "prefix"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.routing_instance", "testacc_securityNATSttRule"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.prefix", "192.0.2.128/28"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"destination_address", "192.0.2.0/27"),
						resource.TestCheckResourceAttr("junos_security_nat_static_rule.testacc_securityNATSttRule",
							"then.prefix", "192.0.2.64/27"),
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
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
