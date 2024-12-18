package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityIdpCustomAttack_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					// 1
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 2
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 3
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 4
					ResourceName:      "junos_security_idp_custom_attack.testacc_idpCustomAttack",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 5
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 6
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 7
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
