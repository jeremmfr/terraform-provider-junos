package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityNatSource_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
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
					ConfigDirectory: config.TestStepDirectory(),
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
