package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceServicesSecurityIntelligenceProfile_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_security_intelligence_profile.testacc_svcSecIntelProfile",
							"rule.#", "1"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_security_intelligence_profile.testacc_svcSecIntelProfile",
							"rule.#", "4"),
					),
				},
				{
					ResourceName:      "junos_services_security_intelligence_profile.testacc_svcSecIntelProfile",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
