package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceApplicationSet_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.#", "1"),
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.0", "junos-ssh"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.#", "2"),
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.0", "junos-ssh"),
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.1", "junos-telnet"),
					),
				},
				{
					ResourceName:      "junos_application_set.testacc_app_set",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
