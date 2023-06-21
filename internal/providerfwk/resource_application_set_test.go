package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosApplicationSet_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosApplicationSetConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.#", "1"),
						resource.TestCheckResourceAttr("junos_application_set.testacc_app_set", "applications.0", "junos-ssh"),
					),
				},
				{
					Config: testAccJunosApplicationSetConfigUpdate(),
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

func testAccJunosApplicationSetConfigCreate() string {
	return `
resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh"]
}
`
}

func testAccJunosApplicationSetConfigUpdate() string {
	return `
resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh", "junos-telnet"]
  application_set = [
    junos_application_set.testacc_app_set2.name
  ]
}
resource "junos_application_set" "testacc_app_set2" {
  name         = "testacc_app_set2"
  applications = ["junos-ftp"]
  description  = "testacc appset2"
}
`
}
