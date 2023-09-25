package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_source_address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_source_address.*", "testacc_address1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_destination_address.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_destination_address.*", "any"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_application.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_application.*", "junos-ssh"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.log_init", "true"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.log_close", "true"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.count", "true"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.#", "2"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.1.then", "reject"),
						resource.TestCheckTypeSetElemAttr("junos_security_policy.testacc_securityPolicy",
							"policy.1.match_source_address.*", "testacc_address1"),
					),
				},
				{
					ResourceName:      "junos_security_policy.testacc_securityPolicy",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
