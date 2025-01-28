package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityDynamicAddressName_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_name.testacc_dyn_add_name",
							"profile_feed_name", "feedtfacc"),
						resource.TestCheckResourceAttr("junos_security_dynamic_address_name.testacc_dyn_add_name2",
							"profile_category.property.#", "2"),
					),
				},
				{
					ResourceName:      "junos_security_dynamic_address_name.testacc_dyn_add_name",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_dynamic_address_name.testacc_dyn_add_name2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("junos_security_dynamic_address_name.testacc_dyn_add_name",
							"profile_category.name"),
					),
				},
			},
		})
	}
}
